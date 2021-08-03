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
package utils

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	v1 "apiserver/pkg/v1"
)

// Configure loads in configuration parameters from ENV vars and the api yaml config file.
// This function returns an ApiConfig object representing the server configuration.
func Configure() (*v1.ApiConfig, error) {
	serviceconfig, err := serviceConfig()
	if err != nil {
		return &v1.ApiConfig{}, err
	}

	dbconfig, err := dbConfig()
	if err != nil {
		return &v1.ApiConfig{}, err
	}

	return &v1.ApiConfig{
		ServiceConfig: serviceconfig,
		DBConfig:      dbconfig,
	}, nil
}

func serviceConfig() (*v1.ServiceConfig, error) {
	// All environment vars for the API server should be prefixed with 'PLANET_'
	// eg. 'export PLANET_DB_PASSWORD="hunter2"'
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
	if err != nil {
		return &v1.ServiceConfig{}, err
	}

	// Defaults
	viper.SetDefault("ServicePort", "8080")

	var serviceconfig v1.ServiceConfig
	err = viper.Unmarshal(&serviceconfig)
	if err != nil {
		return &v1.ServiceConfig{}, err
	}

	err = validateConfig(serviceconfig)
	return &serviceconfig, err
}

func dbConfig() (*v1.DBConfig, error) {
	// All environment vars for the API server should be prefixed with 'PLANET_'
	// eg. 'export PLANET_DB_PASSWORD="hunter2"'
	viper.SetEnvPrefix("planet")

	// Defaults
	viper.SetDefault("db_port", "5432")

	viper.AutomaticEnv()

	dbconfig := v1.DBConfig{
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

		fmt.Println("Error: Validation failed for " + err.StructNamespace())
		if err.Tag() == "required" {
			if field, ok := obj.FieldByName(err.Field()); ok {
				if env, ok := field.Tag.Lookup("env"); ok {
					fmt.Println("'" + env + "' is a required environment variable.")
				}
			}
		}
		fmt.Println()
	}
	return fmt.Errorf("database environment variable validation Failed")
}
