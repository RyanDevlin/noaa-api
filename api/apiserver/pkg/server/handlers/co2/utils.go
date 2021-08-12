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

package co2

import (
	"apiserver/pkg/database"
	"apiserver/pkg/database/models"
	"apiserver/pkg/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// parseParams returns a list of SQL WHERE directives and a map of internal arguments
// to the server, derived from http.Request parameters.
func ParseParams(r *http.Request) ([]string, map[string]interface{}, *utils.ServerError) {
	params := utils.ParseQuery(r)
	var sqlFilters []string
	internalArgs := make(map[string]interface{})

	for key, val := range params {
		err := parseParam(key, val, r.URL.Path, &sqlFilters, internalArgs)
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
func parseParam(filterType string, params []string, urlPath string, sqlFilters *[]string, internalArgs map[string]interface{}) error {

	switch filterType {
	case "year", "month":
		result, err := dateParse(params, filterType)
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "gt":
		result, err := ppmParse(params, urlPath, ">")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "lt":
		result, err := ppmParse(params, urlPath, "<")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "gte":
		result, err := ppmParse(params, urlPath, ">=")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "lte":
		result, err := ppmParse(params, urlPath, "<=")
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
	case "limit":
		result, err := validateInt(params, 0)
		if err != nil {
			return err
		}
		internalArgs[filterType] = result
	case "offset":
		result, err := validateInt(params, 0)
		if err != nil {
			return err
		}
		internalArgs[filterType] = result
	case "page":
		result, err := validateInt(params, 1)
		result-- // validateInt ensures result > 0. This is done so page # '1' is indexed as '0'.
		if err != nil {
			return err
		}
		internalArgs[filterType] = result
	case "pretty":
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

func ppmParse(params []string, urlPath string, comparison string) (string, error) {

	err := validatePpm(params)
	if err != nil {
		return "", err
	}

	// Note: validatePpm checks that params only has one element so it is okay to hardcode params[0] here.
	switch urlPath {
	case "/v1/co2/weekly":
		return "average " + comparison + " " + params[0], nil
	case "/v1/co2/weekly/increase":
		return "increase_since_1800 " + comparison + " " + params[0], nil
	default:
		return "", fmt.Errorf("the path '%v' is not known", urlPath)
	}
}

// parseInternalArgs iterates through arguments originally provided as query params and changes the default query accordingly.
func ParseInternalArgs(internalArgs map[string]interface{}, query *database.DBQuery) error {
	for key, val := range internalArgs {
		switch key {
		case "simple":
			if result, ok := val.(bool); ok {
				query.Cols = []string{"year", "month", "day", "average", "increase_since_1800"}
				query.Simple = result
			}
		case "limit":
			if result, ok := val.(int); ok {
				query.Limit = result
			}
		case "offset":
			if result, ok := val.(int); ok {
				query.Offset = result
			}
		case "page":
			if result, ok := val.(int); ok {
				query.Page = result
			}
		case "pretty":
			if result, ok := val.(bool); ok {
				query.Pretty = result
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
		return false, fmt.Errorf("malformed query parameters, only one boolean value allowed for this argument")
	}

	result, err := strconv.ParseBool(param[0])
	if err != nil {
		return false, fmt.Errorf("malformed query parameters, invalid boolean value")
	}
	return result, nil
}

// validateBool validates an integer parameter.
func validateInt(param []string, min int) (int, error) {
	if len(param) != 1 {
		return -1, fmt.Errorf("malformed query parameters, only one integer value allowed for this argument")
	}

	result, err := strconv.ParseInt(param[0], 10, 32)
	if err != nil {
		return -1, fmt.Errorf("malformed query parameters, invalid integer value")
	}

	if int(result) < min {
		return 0, fmt.Errorf("malformed query parameters, integer value cannot be less than %v", min)
	}
	return int(result), nil
}
