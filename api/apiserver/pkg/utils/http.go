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

API version: 0.1.0
Contact: planetpulse.api@gmail.com
*/

package utils

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
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

// RequestID represents a request's unique ID value. This is used with a golang
// context to trace the request through log messages. It should be used as the key
// when calling context.WithValue() and should always be set to "RequestID". This
// allows the server logic to look up "RequestID" in the request context and obtain
// the ID value.
type RequestID string

// RequestIdDefaultKey is the default value used as a "tag" to extract the RequestIdDefaultKey
// from a request's context.
const RequestIdDefaultKey = RequestID("RequestID")

// SetReqId initializes a new reuest with a unique ID value in a safe manner.
// Setting the requestId key/value pair without this function will usually result
// in server errors.
func SetReqId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		id := uuid.New()
		log.Tracef("Initializing new request with ID: %v\n", id)

		ctx = context.WithValue(ctx, RequestIdDefaultKey, id.String())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetReqId extracts a request's unique ID value and returns it as a string.
// If extraction fails, an error is returned.
func GetReqId(r *http.Request) (string, error) {
	ctx := r.Context()
	id, ok := ctx.Value(RequestIdDefaultKey).(string)
	if !ok {
		return "", fmt.Errorf("request ID key was not set with SetReqId function")
	}
	return id, nil
}
