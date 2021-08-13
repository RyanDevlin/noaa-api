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

package handlers

import (
	"apiserver/pkg/database"
	utils "apiserver/pkg/utils"
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// ApiHandler exposes HTTP handler methods.
type ApiHandler struct {
	Handler ApiHandlerFunc
	Config  *ApiHandlerConfig
}

// ApiHandlerConfig represents configuration parameters to be passed to an ApiHandlerFunc.
// It is a struct to allow for future extensions.
type ApiHandlerConfig struct {
	Database *database.Database
}

// ApiHandlerFunc represents an http handler used to serve data at a specific URL path.
type ApiHandlerFunc func(context.Context, *ApiHandlerConfig, http.ResponseWriter, *http.Request) *utils.ServerError

// NewHandler wraps the ServeHTTP method. It returns an http.Handler used to call ServeHTTP and time its execution for logging purposes.
func NewHandler(ctx context.Context, handler ApiHandler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler.ServeHTTP(ctx, w, r)

		log.Infof(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

// ServeHTTP executes an http handler, logs any errors during execution,
// and rerturns error messages to the client.
func (apiHandler ApiHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if e := apiHandler.Handler(ctx, apiHandler.Config, w, r); e != nil {
		utils.ErrorLog(e)
		utils.HttpJsonError(w, e)
	}
}
