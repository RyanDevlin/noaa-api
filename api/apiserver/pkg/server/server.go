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
	utils "apiserver/pkg/utils"
	"context"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Start initializes the API server and begins listening on the configured port.
func (apiserver *ApiServer) Start() {
	if err := apiserver.ServerInit(); err != nil {
		utils.ErrorLog(err)
	}

	defer apiserver.Database.DB.Close()
	log.Info("Server started.")

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(apiserver.Config.HttpPort), apiserver.Router))
}

// ServerInit initializes the API server. The initialization process loads configuration data
// from config.yaml and environment variables, configures the logger, creates a top level context, establishes
// a database connection, and generates a router to forward requests to handler functions.
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

	// Create a toplevel conext to be passed to all handlers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Generate routes
	apiserver.Router = apiserver.NewRouter(ctx, apiserver.CreateRoutes())

	// Establish database connection. If this fails the server will recover and
	// begin serving, but will only return error messages to the client until a
	// db connection is established.
	err = apiserver.Database.Connect()
	if err != nil {
		utils.ErrorLog(utils.NewError(err, "error establishing database connection", 500, false))
	}

	return nil
}
