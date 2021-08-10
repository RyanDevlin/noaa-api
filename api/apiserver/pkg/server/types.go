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
	"apiserver/pkg/database"
	"apiserver/pkg/server/handlers"

	"github.com/gorilla/mux"
)

// ApiServer provides a way to interact with server components and underlying server methods.
type ApiServer struct {
	Config   *ApiConfig
	Database *database.Database
	Router   *mux.Router
}

// ApiConfig represents configuration parameters for the API server
type ApiConfig struct {
	// (OPTIONAL) The port the server will listen for HTTP traffic on
	HttpPort int

	// (OPTIONAL) The port the server will listen for HTTPS traffic on
	HttpsPort int

	// (OPTIONAL) The global server log level
	LogLevel int
}

// Route represents an HTTP route (a mapping from a URL path to a handler function).
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler handlers.ApiHandler
}

// Routes represents a list of Routes
type Routes []Route

// YamlConfig represents all parameters to be loaded from the config.yaml file.
type YamlConfig struct {
	// (OPTIONAL) The port the server will listen for HTTP traffic on
	HttpPort int `env:"false" name:"HttpPort" validate:"gte=0,lte=65535"`

	// (OPTIONAL) The port the server will listen for HTTPS traffic on
	HttpsPort int `env:"false" name:"HttpsPort" validate:"gte=0,lte=65535"`

	// (OPTIONAL) The global server log level
	LogLevel int `env:"false" name:"LogLevel" validate:"gte=0,lte=6"`

	// (OPTIONAL) The connection timeout in seconds used when connecting to the database
	DBConnTimeout int `env:"false" name:"DBConnTimeout" validate:"gte=0,lte=120"`
}

// EnvConfig represents all parameters to be loaded from environment variables.
type EnvConfig struct {
	// The database endpoint
	DBHost string `env:"true" name:"PLANET_DB_HOST" validate:"required"`

	// The database username
	DBUser string `env:"true" name:"PLANET_DB_USER" validate:"required"`

	// The database password
	DBPass string `env:"true" name:"PLANET_DB_PASS" validate:"required"`

	// (OPTIONAL) The port the database listens on
	DBPort int `env:"true" name:"PLANET_DB_PORT" validate:"gte=0,lte=65535"`
}
