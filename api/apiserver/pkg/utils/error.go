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
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

// ServerError represents an error on the server. It is used not only to provide context to the internal
// logger, but also to provide context to the response sent to the client.
type ServerError struct {
	Error    error
	Message  string
	HttpCode int
	Fatal    bool

	File string
	Line int
}

// ErrorResp represents the format of the JSON response returned to the client when an error occurs.
type ErrorResp struct {
	Description string
	Content     ErrorType
}

// ErrorType represents the context of an error returned to the client.
type ErrorType struct {
	Code    int
	Message string
}

// NewError returns a new ServerError object used to encode contextual information about a runtime error
func NewError(err error, message string, code int, fatal bool) *ServerError {
	_, file, line, _ := runtime.Caller(1)
	return &ServerError{Error: err, Message: message, HttpCode: code, Fatal: fatal, File: filepath.Base(file), Line: line}
}

// ErrorLog uses the configured logger to report context from a server error.
// Based on the server error configuration, this function may trigger a program exit and return 1.
func ErrorLog(serverError *ServerError) {
	errString := fmt.Sprintf(
		"%s:%d (%s) - %s (%s).",
		serverError.File,
		serverError.Line,
		http.StatusText(serverError.HttpCode),
		serverError.Error.Error(),
		serverError.Message,
	)
	if serverError.Fatal {
		// A stacktrace is only logged if loglevel >= 6
		log.WithField("stack", string(debug.Stack())).Trace("STACK TRACE:")
		log.Fatal(errString)
		return
	}
	log.Errorf(errString)
	log.WithField("stack", string(debug.Stack())).Trace("STACK TRACE:")
}

// HttpJsonError extracts metadata from a ServerError and returns this
// information to the client.
func HttpJsonError(w http.ResponseWriter, err *ServerError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.HttpCode)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.SetEscapeHTML(false)
	enc.Encode(
		ErrorResp{
			Description: http.StatusText(err.HttpCode),
			Content: ErrorType{
				Code:    err.HttpCode,
				Message: err.Message,
			},
		},
	)
}
