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

package models

import (
	"time"
)

const (
	// The maximum ppm value that may be used in a query for Co2 data
	Co2PpmMax = 1000
	
	// The minimum ppm value that may be used in a query for Co2 data
	Co2PpmMin = 0
)

// Co2Table represents a mapping of dates to Co2Entry or Co2EntrySimple structs.
// The index of the Co2Table map must be '<year>-<month>-<day>'
type Co2Table map[string]interface{}

// Co2Entry represents the JSON data to be returned from an individual Co2 measurement in the database.
type Co2Entry struct {
	Year                  int
	Month                 int
	Day                   int
	DateDecimal           float32
	Average               float32
	NumDays               int
	OneYearAgo            float32
	TenYearsAgo           float32
	IncSincePreIndustrial float32
	Timestamp             time.Time
}

// Co2EntrySimple represents the simplified JSON data to be returned from an individual Co2 measurement in the database.
type Co2EntrySimple struct {
	Average               float32
	IncSincePreIndustrial float32
}

/*

	Here lies a ton of code I wrote to filter data from the database server-side.
	This has now been simplified by properly querying the database, thus this
	code is useless. I'm leaving it here on the off chance I need it again during this
	initial architecture push. Also I put too much work into it :(


	Draw things out on paper before you write the code 🤦‍♂️.

        .
       -|-
        |
    .-'~~~`-.
  .'         `.
  |  R  I  P  |
  |           |
  |           |
\\|           |//
^^^^^^^^^^^^^^^^^
*/

/*
func (co2Table *Co2Table) Filter(r *http.Request) *utils.ServerError {
	params := utils.ParseQuery(r)
	for key, val := range params {
		query := Query{key, val}
		err := query.execute(co2Table)
		if err != nil {
			message := err.Error() + ": " + utils.ParseQuery(r).Encode()
			return utils.NewError(fmt.Errorf("error when parsing query parameters"), message, 400, false)
		}
	}
	return nil
}

func (query Query) execute(data *Co2Table) error {
	var err error
	switch query.filterType {
	case "year":
		err = dateParse(data, query.params, 0)
	case "month":
		err = dateParse(data, query.params, 1)
	case "gt":
		err = filterPpmCompare(data, query.params, ">")
	case "lt":
		err = filterPpmCompare(data, query.params, "<")
	case "gte":
		err = filterPpmCompare(data, query.params, ">=")
	case "lte":
		err = filterPpmCompare(data, query.params, "<=")
	case "simple":
		err = simple(data, query.params)
	}
	return err
}

// ====	HELPER FUNCTIONS ====

func dateParse(table *Co2Table, params []string, index int) error {
	result := make(Co2Table)

	for key, val := range *table {
		date := strings.Split(key, "-")
		if index < 0 || index > len(date) {
			return fmt.Errorf("internal error") //TODO: is this correct?
		}

		match, err := paramMatch(params, date[index])
		if err != nil {
			return err
		}
		if match {
			result[key] = val
		}
	}

	*table = result
	return nil
}

func paramMatch(params []string, input string) (bool, error) {
	for _, v := range params {
		if _, err := strconv.Atoi(v); err != nil {
			return false, fmt.Errorf("malformed query parameters, invalid date value")
		}
		if v == input {
			return true, nil
		}
	}
	return false, nil
}

func filterPpmCompare(table *Co2Table, params []string, comparison string) error {
	result := make(Co2Table)

	ppm, err := validateAndDigestPpm(params)
	if err != nil {
		return err
	}
	for key, val := range *table {
		average := float32(reflect.ValueOf(val).FieldByName("Average").Float())

		switch comparison {
		case ">":
			if average > ppm {
				result[key] = val
			}
		case "<":
			if average < ppm {
				result[key] = val
			}
		case ">=":
			if average >= ppm {
				result[key] = val
			}
		case "<=":
			if average <= ppm {
				result[key] = val
			}
		default:
			return fmt.Errorf("(internal) malformed ppm comparison string '%s'", comparison) //TODO: is this correct?
		}
	}

	*table = result
	return nil
}

func validateAndDigestPpm(param []string) (float32, error) {
	if len(param) != 1 {
		return 0, fmt.Errorf("malformed query parameters, too many ppm constraints")
	}

	ppm, err := strconv.ParseFloat(param[0], 32)
	if err != nil {
		return 0, fmt.Errorf("malformed query parameters, ppm value should be a decimal number")
	}
	if !(ppm <= Co2PpmMax && ppm >= Co2PpmMin) {
		return 0, fmt.Errorf("malformed query parameters, ppm query range is %v to %v", Co2PpmMin, Co2PpmMax)
	}
	return float32(ppm), nil
}

func validateAndDigestBool(param []string) (bool, error) {
	if len(param) != 1 {
		return false, fmt.Errorf("malformed query parameters, only one boolean value per argument")
	}

	result, err := strconv.ParseBool(param[0])
	if err != nil {
		return false, fmt.Errorf("malformed query parameters, invalid boolean value")
	}
	return result, nil
}

func simple(table *Co2Table, params []string) error {
	result := make(Co2Table)

	param, err := validateAndDigestBool(params)
	if err != nil {
		return err
	}
	if !param {
		table = &result
		return nil
	}

	for key, val := range *table {
		result[key] = convToSimple(val)
	}

	*table = result
	return nil
}

func convToSimple(data interface{}) interface{} {
	simple := Co2EntrySimple{}

	dataVal := reflect.ValueOf(data)
	simpleVal := reflect.ValueOf(simple)
	simplePtr := reflect.ValueOf(&simple)

	for i := 0; i < simpleVal.NumField(); i++ {
		field := dataVal.FieldByName(simpleVal.Type().Field(i).Name)
		simplePtr.Elem().Field(i).Set(field)
	}
	return simple
}
*/
