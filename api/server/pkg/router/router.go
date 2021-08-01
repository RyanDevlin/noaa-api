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

package router

import (
	"net/http"
	"strings"

	ep "github.com/RyanDevlin/planetpulse/api/server/pkg/endpoints"
	utils "github.com/RyanDevlin/planetpulse/api/server/pkg/utils"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = utils.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		ep.Index,
	},

	Route{
		"Index",
		"GET",
		"/v1/",
		ep.Index,
	},

	Route{
		"Co2WeeklyGet",
		strings.ToUpper("Get"),
		"/v1/co2/weekly",
		ep.Co2WeeklyGet,
	},

	Route{
		"Co2WeeklyIdGet",
		strings.ToUpper("Get"),
		"/v1/co2/weekly/{id}",
		ep.Co2WeeklyIdGet,
	},

	Route{
		"Co2WeeklyIncreaseGet",
		strings.ToUpper("Get"),
		"/v1/co2/weekly/increase",
		ep.Co2WeeklyIncreaseGet,
	},

	Route{
		"Co2WeeklyPpmGet",
		strings.ToUpper("Get"),
		"/v1/co2/weekly/{ppm}",
		ep.Co2WeeklyPpmGet,
	},

	Route{
		"HealthGet",
		strings.ToUpper("Get"),
		"/v1/health",
		ep.HealthGet,
	},
}
