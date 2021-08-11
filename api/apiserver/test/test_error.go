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

package test

import (
	"apiserver/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

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

// ErrorLog in the test package is used to mock the behavior of ErrorLog
// in the utils package.
func ErrorLog(t *testing.T, serverError *utils.ServerError) {
	errString := fmt.Sprintf(
		"MOCK ERROR LOG = %s:%d (%s) - %s (%s).",
		serverError.File,
		serverError.Line,
		http.StatusText(serverError.HttpCode),
		serverError.Error.Error(),
		serverError.Message,
	)
	if serverError.Fatal {
		t.Log(errString)
		t.Log("MOCK ERROR LOG = This error would have been fatal.")
		return
	}
	t.Log(errString)
}

// HttpJsonError in the test package is used to mock the behavior of HttpJsonError
// in the utils package.
func HttpJsonError(w http.ResponseWriter, err *utils.ServerError) {
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
