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

package co2

import (
	"apiserver/pkg/database"
	"apiserver/pkg/database/models"
	"apiserver/pkg/utils"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// parseParams returns a list of SQL WHERE directives and a map of internal arguments
// to the server, derived from http.Request parameters. The urlParams tells the function if
// the parameters should be derived mainly from the query params (eg. '/v1/co2/weekly?year=2020&gte=417')
// or the url path (eg. '/v1/co2/weekly/317.22?simple=true'). This is needed because when specifying
// a specific resource in the url path, filters like gt,gte,lt,lte, etc. are not needed as only one
// resource is returned.
func ParseParams(r *http.Request, pathParam bool, sortBy string) ([]string, map[string]interface{}, *utils.ServerError) {
	params := utils.ParseQuery(r)
	var sqlFilters []string
	internalArgs := make(map[string]interface{})
	var err error

	if pathParam {
		err = parsePathParams(r.URL.Path, sortBy, &sqlFilters)
		if err != nil {
			message := err.Error() + ": " + path.Dir(r.URL.Path) + "=[" + path.Base(r.URL.Path) + "]"
			return nil, nil, utils.NewError(fmt.Errorf("error when parsing path parameter"), message, 400, false)
		}
	}

	for key, val := range params {
		if pathParam {
			err = parseSingleResource(key, val, sortBy, &sqlFilters, internalArgs)
		} else {
			err = parseParam(key, val, sortBy, &sqlFilters, internalArgs)
		}
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
func parseParam(filterType string, params []string, sortBy string, sqlFilters *[]string, internalArgs map[string]interface{}) error {

	switch filterType {
	case "year", "month":
		result, err := dateParse(params, filterType)
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "gt":
		ppm, err := getPPM(params, true)
		if err != nil {
			return err
		}
		result, err := ppmParse(ppm, sortBy, ">")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "lt":
		ppm, err := getPPM(params, false)
		if err != nil {
			return err
		}
		result, err := ppmParse(ppm, sortBy, "<")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "gte":
		ppm, err := getPPM(params, true)
		if err != nil {
			return err
		}
		result, err := ppmParse(ppm, sortBy, ">=")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	case "lte":
		ppm, err := getPPM(params, false)
		if err != nil {
			return err
		}
		result, err := ppmParse(ppm, sortBy, "<=")
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
		result, err := validateInt(params, 0, 10000)
		if err != nil {
			return err
		}
		internalArgs[filterType] = result
	case "offset":
		result, err := validateInt(params, 0, 10000)
		if err != nil {
			return err
		}
		internalArgs[filterType] = result
	case "page":
		result, err := validateInt(params, 1, 10000)
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

// parseSingleResource appends a single boolean expression to the sqlFilters list. This list of expressions is later passed
// directly to the WHERE clause of an SQL query. parseParam also will add specific arguments to the internalArgs map
// to be later used by the server.
func parseSingleResource(filterType string, params []string, urlPath string, sqlFilters *[]string, internalArgs map[string]interface{}) error {

	switch filterType {
	case "simple":
		result, err := validateBool(params)
		if err != nil {
			return err
		}
		internalArgs[filterType] = result
	case "limit":
		result, err := validateInt(params, 0, 10000)
		if err != nil {
			return err
		}
		internalArgs[filterType] = result
	case "offset":
		result, err := validateInt(params, 0, 10000)
		if err != nil {
			return err
		}
		internalArgs[filterType] = result
	case "page":
		result, err := validateInt(params, 1, 10000)
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

func parsePathParams(urlPath string, sortBy string, sqlFilters *[]string) error {
	val := path.Base(urlPath)

	switch sortBy {
	case "average", "increase":
		_, err := validatePpm(val)
		if err != nil {
			return err
		}
		result, err := ppmParse(val, sortBy, "=")
		if err != nil {
			return err
		}
		*sqlFilters = append(*sqlFilters, result)
	default:
		return fmt.Errorf("cannot sort database results by '%v'. Unknown column", sortBy)
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

func ppmParse(ppm string, sortBy string, comparison string) (string, error) {
	switch sortBy {
	case "average":
		return "average " + comparison + " " + ppm, nil
	case "increase":
		return "increase_since_1800 " + comparison + " " + ppm, nil
	default:
		return "", fmt.Errorf("cannot sort results by '%v'. '%v' is not a column in the database", sortBy, sortBy)
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

func getPPM(array []string, max bool) (string, error) {
	target, err := validatePpm(array[0])
	if err != nil {
		return "", err
	}
	for _, value := range array {
		curr, err := validatePpm(value)
		if err != nil {
			return "", err
		}
		if max && curr > target {
			target = curr
		} else if !max && curr < target {
			target = curr
		}
	}
	return strconv.FormatFloat(target, 'f', 2, 32), nil
}

// validatePpm validates a ppm parameter against the current API spec.
func validatePpm(ppmStr string) (float64, error) {

	ppm, err := strconv.ParseFloat(ppmStr, 32)
	if err != nil {
		return 0, fmt.Errorf("malformed query parameters, ppm value should be a decimal number")
	}
	if !(ppm <= models.Co2PpmMax && ppm >= models.Co2PpmMin) {
		return 0, fmt.Errorf("malformed query parameters, ppm query range is %v to %v", models.Co2PpmMin, models.Co2PpmMax)
	}
	return ppm, nil
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
func validateInt(param []string, min int, max int) (int, error) {
	if len(param) != 1 {
		return -1, fmt.Errorf("malformed query parameters, only one integer value allowed for this argument")
	}

	result, err := strconv.ParseInt(param[0], 10, 32)
	if err != nil {
		return -1, fmt.Errorf("malformed query parameters, invalid integer value")
	}

	if int(result) < min {
		return 0, fmt.Errorf("malformed query parameters, integer value cannot be less than %v", min)
	} else if int(result) > max {
		return 0, fmt.Errorf("malformed query parameters, integer value cannot be greater than %v", max)
	}
	return int(result), nil
}
