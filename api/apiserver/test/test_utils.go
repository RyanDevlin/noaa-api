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

package test

import (
	"apiserver/pkg/utils"
	"bytes"
	"context"
	"net/http"
	"testing"

	// The blank import here is used to import the pq PostgreSQL drivers

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func PrintServerResponse(t *testing.T, r *http.Response, body []byte) {
	t.Logf("Server Response: \n\tStatus Code: %v\n\tHeader Content-Type: %v\n\tBody: \n%v\n", r.StatusCode, r.Header.Get("Content-Type"), string(indent([]byte{'\t'}, body)))
}

// indent was copied directly from https://pkg.go.dev/github.com/openconfig/goyang/pkg/indent
func indent(indent, b []byte) []byte {
	if len(indent) == 0 || len(b) == 0 {
		return b
	}
	lines := bytes.SplitAfter(b, []byte{'\n'})
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return bytes.Join(append([][]byte{{}}, lines...), indent)
}

// SetReqIdTest is used to attach an ID to a request context during testing
func SetReqIdTest(r *http.Request) *http.Request {
	ctx := r.Context()
	id := uuid.New()
	ctx = context.WithValue(ctx, utils.RequestIdDefaultKey, id.String())
	return r.WithContext(ctx)
}
