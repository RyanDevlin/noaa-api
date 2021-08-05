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
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configure loads in configuration parameters from ENV vars and the api yaml config file.
// This function returns an ApiConfig object representing the server configuration.
func (apiserver *ApiServer) configure() error {
	serviceconfig, err := serviceConfig()
	if err != nil {
		return err
	}

	dbconfig, err := dbConfig()
	if err != nil {
		return err
	}

	apiserver.Config = &ApiConfig{
		ServiceConfig: serviceconfig,
		DBConfig:      dbconfig,
	}
	apiserver.configured = true
	return nil
}

func serviceConfig() (*ServiceConfig, error) {
	// Values for the service config are read from a config.yaml file in the same directory as the executable
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return &ServiceConfig{}, err
	}

	// Defaults
	viper.SetDefault("HttpPort", "8080")
	viper.SetDefault("HttpsPort", "8443")
	viper.SetDefault("LogLevel", "4")
	viper.SetDefault("DBConnTimeout", "5")

	var serviceconfig ServiceConfig
	err = viper.Unmarshal(&serviceconfig)
	if err != nil {
		return &ServiceConfig{}, err
	}

	err = validateConfig(serviceconfig)
	return &serviceconfig, err
}

func dbConfig() (*DBConfig, error) {
	// All environment vars for the API server should be prefixed with 'PLANET_'
	// eg. 'export PLANET_DB_PASSWORD="hunter2"'
	viper.SetEnvPrefix("planet")

	// Defaults
	viper.SetDefault("db_port", "5432")

	viper.AutomaticEnv()

	dbconfig := DBConfig{
		DBHost: viper.GetString("db_host"),
		DBUser: viper.GetString("db_user"),
		DBPass: viper.GetString("db_pass"),
		DBPort: viper.GetString("db_port"),
	}
	err := validateConfig(dbconfig)
	return &dbconfig, err
}

func validateConfig(config interface{}) error {
	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		return validateErrorHandler(reflect.TypeOf(config), err)
	}
	return nil
}

func validateErrorHandler(obj reflect.Type, err error) error {
	for _, err := range err.(validator.ValidationErrors) {

		if err.Tag() == "required" {
			if field, ok := obj.FieldByName(err.Field()); ok {
				if env, ok := field.Tag.Lookup("env"); ok {
					log.Error("'" + env + "' is a required environment variable.")
				}
			}
		}
	}
	return fmt.Errorf("failed to load database parameters from environment variables")
}
