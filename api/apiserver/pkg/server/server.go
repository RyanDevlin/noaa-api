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
	utils "apiserver/pkg/utils"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (apiserver *ApiServer) Start() {
	if err := apiserver.ServerInit(); err != nil {
		utils.HttpErrorLog(err)
	}

	defer apiserver.Db.Close()
	log.Info("Server started.")

	//log.Fatal(http.ListenAndServeTLS(":"+apiserver.Config.HttpsPort, "apiserver.crt", "apiserver.key", apiserver.Router))
	log.Fatal(http.ListenAndServe(":"+apiserver.Config.HttpPort, apiserver.Router))
}

func (apiserver *ApiServer) ServerInit() *utils.ServerError {
	// Configure server parameters. If this fails, a fatal log.Fatal will be called
	// and the server process will be terminated
	err := apiserver.configure()
	if err != nil {
		return utils.NewError(err, "apiserver configuration failed", 500, true)
	}

	// Setup Logging Level
	log.SetLevel(log.Level(apiserver.Config.LogLevel))
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Generate routes
	router := NewRouter(apiserver.CreateRoutes(), apiserver)
	apiserver.Router = router

	// Establish database connection. If this fails the server will recover and
	// begin serving, but will only return error messages to the client until a
	// db connection is established.
	err = apiserver.DBConnect()
	if err != nil {
		return utils.NewError(err, "error establishing database connection", 500, false)
	}

	return nil
}
