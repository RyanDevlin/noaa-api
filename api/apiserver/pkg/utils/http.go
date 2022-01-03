/*
Copyright 2021 The PlanetPulse Authors.

Planet Pulse is an API designed to serve climate data pulled from NOAA's
Global Monitoring Laboratory FTP server. This API is based on the
OpenAPI v3 specification.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

A copy of the GNU General Public License can be found here:
https://www.gnu.org/licenses/

Contact: planetpulse.api@gmail.com
*/

package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ParseQuery expands the query parameters passed to an endpoint
// to account for array-like parameters.
// This allows one to search for, say, the following:
//
//		example.com/v1/stuff?day=1,2,3&day=4
//
// The expansion will allow the day slice to become:
//
//		day := ["1", "2", "3". "4"]
func ParseQuery(r *http.Request) url.Values {
	params := r.URL.Query()
	for key, val := range params {
		var expanded []string
		for _, elem := range val {
			array := strings.Split(elem, ",")
			expanded = append(expanded, array...)
		}
		params[key] = expanded
	}
	return params
}

// SetCSPHeaders is called by each route to set Content Security Policy headers in http responses.
// Currently this is used mainly to allow the favicon to be requested by the client.
func SetCSPHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src '*.planetpulse.io'; object-src: 'none';")
		next.ServeHTTP(w, r)
	})
}

// SetCORSHeaders is called by each route to set Cross-Origin Resource Sharing headers in http responses.
// This is used to allow requests from any origin.
func SetCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// gzipResponseWriter TODO
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write writes to the *gzip.Writer io.Writer which compresses all data
// and writes to the original http.ResponseWriter
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Gzip is a middleware function that enables gzip compression on all
// http responses so long as the header "Accept-Encoding": gzip is
// present in the request
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)
	})
}
