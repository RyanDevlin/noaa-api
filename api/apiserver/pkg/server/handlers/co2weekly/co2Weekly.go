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

package co2weekly

import (
	"apiserver/pkg/database"
	"apiserver/pkg/database/models"
	"apiserver/pkg/server/handlers"
	utils "apiserver/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func Get(ctx context.Context, handlerConfig *handlers.ApiHandlerConfig, w http.ResponseWriter, r *http.Request) *utils.ServerError {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	query := database.DBQuery{
		Table:   "public.co2_weekly_mlo",
		Cols:    []string{"*"},
		OrderBy: "year,month,day",
	}

	filters, internalArgs, err := parseParams(r)
	if err != nil {
		return err
	}

	if len(internalArgs) != 0 {
		parseInternalArgs(internalArgs, &query)
	}

	query.Where = filters

	co2Table, dberr := handlerConfig.Database.Query(query)
	if err != nil {
		return utils.NewError(dberr, "failed to connect to database", 500, false)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(co2Table); err != nil {
		return utils.NewError(err, "error encoding data as json", 500, false)
	}
	return nil
}

// parseParams returns a list of SQL WHERE directives and a map of internal arguments
// to the server, derived from http.Request parameters.
func parseParams(r *http.Request) ([]string, map[string]interface{}, *utils.ServerError) {
	params := utils.ParseQuery(r)
	var sqlFilters []string
	internalArgs := make(map[string]interface{})

	for key, val := range params {
		err := parseParam(key, val, &sqlFilters, internalArgs)
		if err != nil {
			message := err.Error() + ": " + key + "=[" + strings.Join(val, ",") + "]"
			return nil, nil, utils.NewError(fmt.Errorf("error when parsing query parameters"), message, 400, false)
		}
	}
	return sqlFilters, internalArgs, nil
}

// parseParam appends a single boolean expression to the sqlFilters list. This list of expressions is later passed
// directly to the WHERE clause of an SQL query. parseParam also will add specific arguments to the internalArgs map
// to be later used by the server.
func parseParam(filterType string, params []string, sqlFilters *[]string, internalArgs map[string]interface{}) error {

	switch filterType {
	case "year", "month":
		result, err := dateParse(params, filterType)
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "gt":
		result, err := ppmParse(params, ">")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "lt":
		result, err := ppmParse(params, "<")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "gte":
		result, err := ppmParse(params, ">=")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "lte":
		result, err := ppmParse(params, "<=")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "simple":
		result, err := validateBool(params)
		if err != nil {
			return err
		}
		internalArgs[filterType] = result
	}
	return nil
}

// dateParse
func dateParse(params []string, section string) (string, error) {
	result := section + " in ('"
	numParams := len(params) - 1
	for i, v := range params {
		err := validateDate(v, section)
		if err != nil {
			return "", err
		}

		if i == numParams {
			result += v + "')"
			break
		}
		result += v + "', '"
	}
	return result, nil
}

func ppmParse(params []string, comparison string) (string, error) {

	err := validatePpm(params)
	if err != nil {
		return "", err
	}

	// Note: validatePpm checks that params only has one element so it is okay to hardcode this here.
	return "average " + comparison + " " + params[0], nil
}

// parseInternalArgs iterates through arguments originally provided as query params and changes the default query accordingly.
func parseInternalArgs(internalArgs map[string]interface{}, query *database.DBQuery) error {
	for key, val := range internalArgs {
		switch key {
		case "simple":
			if result := val.(bool); result {
				query.Cols = []string{"year", "month", "day", "average", "Increase_since_1800"}
				query.Simple = result
			}

		}
	}
	return nil
}

/* VALIDATION */

// validateDate validates a date parameter against the current API spec.
func validateDate(val string, section string) error {
	date, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("malformed query parameters, invalid date value")
	}

	switch section {
	case "year":
		if date < 0 || date > 3000 {
			return fmt.Errorf("invalid year value. Years must be between 0 and 3000")
		}
	case "month":
		if date < 1 || date > 12 {
			return fmt.Errorf("invalid month value. Months must be between 1 and 12")
		}
	}
	return nil
}

// validatePpm validates a ppm parameter against the current API spec.
func validatePpm(param []string) error {
	if len(param) != 1 {
		return fmt.Errorf("malformed query parameters, too many ppm constraints")
	}

	ppm, err := strconv.ParseFloat(param[0], 32)
	if err != nil {
		return fmt.Errorf("malformed query parameters, ppm value should be a decimal number")
	}
	if !(ppm <= models.Co2PpmMax && ppm >= models.Co2PpmMin) {
		return fmt.Errorf("malformed query parameters, ppm query range is %v to %v", models.Co2PpmMin, models.Co2PpmMax)
	}
	return nil
}

// validateBool validates a boolean parameter.
func validateBool(param []string) (bool, error) {
	if len(param) != 1 {
		return false, fmt.Errorf("malformed query parameters, only one boolean value per argument")
	}

	result, err := strconv.ParseBool(param[0])
	if err != nil {
		return false, fmt.Errorf("malformed query parameters, invalid boolean value")
	}
	return result, nil
}
