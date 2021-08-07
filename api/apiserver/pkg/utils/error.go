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
	v1 "apiserver/pkg/v1"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

type ServerError struct {
	Error    error
	Message  string
	HttpCode int
	Fatal    bool

	File string
	Line int
}

type ApiHandler func(http.ResponseWriter, *http.Request) *ServerError
type ServerHandler func() *ServerError

func NewError(err error, message string, code int, fatal bool) *ServerError {
	_, file, line, _ := runtime.Caller(1)
	return &ServerError{Error: err, Message: message, HttpCode: code, Fatal: fatal, File: filepath.Base(file), Line: line}
}

// ServerHTTP executes an http handler, logs any errors during execution,
// and rerturns error messages to the client.
func (fn ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		HttpErrorLog(e)
		//http.Error(w, strings.Title(e.Message), e.HttpCode)
		JsonError(w, e)
	}
}

// ExecuteInternalTask executes a server handler (some sort of internal server
// process) and logs any errors.
func (fn ServerHandler) ExecuteInternalTask() {
	if e := fn(); e != nil {
		HttpErrorLog(e)
		//http.Error(w, strings.Title(e.Error.Error()), e.HttpCode)
	}
}

func JsonError(w http.ResponseWriter, err *ServerError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.HttpCode)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.SetEscapeHTML(false)
	enc.Encode(
		v1.ErrorResp{
			Description: http.StatusText(err.HttpCode),
			Content: v1.ErrorType{
				Code:    err.HttpCode,
				Message: err.Message,
			},
		},
	)
}

func HttpErrorLog(serverError *ServerError) {
	errString := fmt.Sprintf(
		"%s:%d (%s) - %s (%s).",
		serverError.File,
		serverError.Line,
		http.StatusText(serverError.HttpCode),
		serverError.Error.Error(),
		serverError.Message,
	)
	if serverError.Fatal {
		log.WithField("stack", string(debug.Stack())).Trace("STACK TRACE:")
		log.Fatal(errString)
		return
	}
	log.Errorf(errString)
	log.WithField("stack", string(debug.Stack())).Trace("STACK TRACE:")
}

func InternalErrorLog(serverError *ServerError) {
	errString := fmt.Sprintf(
		"%s:%d - %s (%s).",
		serverError.File,
		serverError.Line,
		serverError.Error.Error(),
		serverError.Message,
	)
	if serverError.Fatal {
		log.WithField("stack", string(debug.Stack())).Trace("STACK TRACE:")
		log.Fatal(errString)
		return
	}
	log.Errorf(errString)
	log.WithField("stack", string(debug.Stack())).Trace("STACK TRACE:")
}
