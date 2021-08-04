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
	"database/sql"
	"log"
	"net/http"

	v1 "apiserver/pkg/v1"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	Config *v1.ApiConfig
	Db     *sql.DB
	Router *mux.Router

	// True if the server has been configured
	configured bool
}

func (apiserver *ApiServer) ServerInit() error {
	log.Printf("Server started.")

	// Configure server parameters
	err := apiserver.configure()
	if err != nil {
		return err
	}

	// Establish database connection
	err = apiserver.planetDBConnect()
	if err != nil {
		return err
	}

	// Generate routes
	router := NewRouter(apiserver.CreateRoutes(), apiserver)
	apiserver.Router = router
	return nil
}

func (apiserver *ApiServer) Start() {
	defer apiserver.Db.Close()
	log.Fatal(http.ListenAndServe(":"+apiserver.Config.ServicePort, apiserver.Router))
}
