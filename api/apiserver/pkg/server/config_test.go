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
package server

import (
	"apiserver/test"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

var configData = []byte(
	"HttpPort: 8080\nLogLevel: 5\nDBConnTimeout: 0",
)

func TestConfig(t *testing.T) {
	log.SetLevel(1)

	apiserver := &ApiServer{}
	serverError := apiserver.ServerInit()
	if serverError == nil {
		t.Error("Expected config file error from ServerInit(), got nil error response.")
		return
	}
	test.ErrorLog(t, serverError)

	err := ioutil.WriteFile("config.yaml", configData, 0755)
	if err != nil {
		t.Errorf("Unable to write config file: %v", err)
		return
	}

	defer func() {
		e := os.Remove("config.yaml")
		if e != nil {
			t.Error(e)
			return
		}
	}()

	db_user := os.Getenv("PLANET_DB_USER")
	db_pass := os.Getenv("PLANET_DB_PASS")
	db_host := os.Getenv("PLANET_DB_HOST")
	if verbose {
		log.SetLevel(4)

		t.Logf("Environment variable PLANET_DB_USER set to: '%s'", db_user)
		t.Logf("Environment variable PLANET_DB_PASS set to: '**************'")
		t.Logf("Environment variable PLANET_DB_HOST set to: '%s'", db_host)
	}
	t.Log("Unsetting environment variables....")

	os.Unsetenv("PLANET_DB_USER")
	os.Unsetenv("PLANET_DB_PASS")
	os.Unsetenv("PLANET_DB_HOST")

	serverError = apiserver.ServerInit()
	if serverError == nil {
		t.Error("Expected environment variable error from ServerInit(), got nil error response.")
		return
	}
	test.ErrorLog(t, serverError)

	t.Log("Resetting environment variables....")
	os.Setenv("PLANET_DB_USER", db_user)
	os.Setenv("PLANET_DB_PASS", db_pass)
	os.Setenv("PLANET_DB_HOST", db_host)

	serverError = apiserver.ServerInit()
	if serverError != nil {
		test.ErrorLog(t, serverError)
		t.Error("Configuration failed during testing.")
		return
	}
}
