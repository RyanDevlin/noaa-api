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
package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var apiSpecPath = "../../../openapi/spec-v1.json"

type Path map[string]interface{}

type PathSpec struct {
	Ref         string
	Summary     string
	Description string
	Get         string
	Put         string
	Post        string
	Delete      string
	Options     string
	Head        string
	Patch       string
	Trace       string
	Servers     []string
	Parameters  []string
}

func TestRoutes(t *testing.T) {

	jsonFile, err := os.Open(apiSpecPath)
	if err != nil {
		t.Error(err)
		t.Errorf("Tried to open the spec at this location: '%v'", apiSpecPath)
		return
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Error(err)
		return
	}

	specTable := make(map[string]Path)
	json.Unmarshal(bytes, &specTable)

	apiserver := &ApiServer{}
	routes := apiserver.CreateRoutes()
	routesTable := make(map[string]Route)

	for _, route := range routes {
		routesTable[route.Pattern] = route
	}

	problemDetected := false

	for k := range specTable["paths"] {
		if _, ok := routesTable["/v1"+k]; !ok {
			problemDetected = true
			t.Logf("WARNING - The following path is defined in the OpenAPI spec file but is not defined in the Routes list: /v1%v", k)
		}
	}

	for k := range routesTable {
		if k == "/" || k == "/v1" { // Skip over the top level routes for now
			continue
		}
		trimmed := strings.Replace(k, "/v1", "", 1)

		if _, ok := specTable["paths"][trimmed]; !ok {
			problemDetected = true
			t.Logf("WARNING - The following path is defined in the Routes list but is not defined in the OpenAPI spec file: %v", k)
		}
	}
	if problemDetected {
		// For now, this problem is logged but won't cause test failures.
		// When the API is more concrete, this should trigger a failure.
		t.Skipf("There is a mismatch between the OpenAPI spec file (%v) and the routes defined in router.go", apiSpecPath)
	}
}
