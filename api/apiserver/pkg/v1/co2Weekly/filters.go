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

package co2Weekly

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type QueryFilter interface {
	Filter(Co2Table) (Co2Table, error)
}

type FilterFunc func(Co2Table) (Co2Table, error)

func (co2Year *Co2Year) Filter(table Co2Table) (Co2Table, error) {
	result := make(map[string]interface{})
	for key, val := range table {
		year := strings.Split(key, "-")[0] // Year is the first portion of the key: 'year-month-day'
		if paramMatch(co2Year.Params, year) {
			result[key] = val
		}
	}
	return result, nil
}

func (co2Month *Co2Month) Filter(table Co2Table) (Co2Table, error) {
	result := make(map[string]interface{})
	for key, val := range table {
		month := strings.Split(key, "-")[1] // Month is the second portion of the key: 'year-month-day'
		if paramMatch(co2Month.Params, month) {
			result[key] = val
		}
	}
	return result, nil
}

func (co2GreaterThan *Co2GreaterThan) Filter(table Co2Table) (Co2Table, error) {
	return filterPpmCompare(table, co2GreaterThan.Params, ">")
}

func (co2LessThan *Co2LessThan) Filter(table Co2Table) (Co2Table, error) {
	return filterPpmCompare(table, co2LessThan.Params, "<")
}

func (gte *Co2Gte) Filter(table Co2Table) (Co2Table, error) {
	return filterPpmCompare(table, gte.Params, ">=")
}

func (lte *Co2Lte) Filter(table Co2Table) (Co2Table, error) {
	return filterPpmCompare(table, lte.Params, "<=")
}

func (co2Simple *Co2Simple) Filter(table Co2Table) (Co2Table, error) {
	result := make(map[string]interface{})

	param, err := validateAndDigestBool(co2Simple.Params)
	if err != nil {
		return nil, err
	}
	if !param {
		return table, nil
	}

	for key, val := range table {
		result[key] = convToSimple(val)
	}
	return result, nil
}

/* ====	HELPER FUNCTIONS ==== */

func paramMatch(params []string, input string) bool {
	for _, v := range params {
		if v == input {
			return true
		}
	}
	return false
}

func filterPpmCompare(table Co2Table, params []string, comparison string) (Co2Table, error) {
	result := make(map[string]interface{})

	ppm, err := validateAndDigestPpm(params)
	if err != nil {
		return nil, err
	}
	fmt.Print(ppm)
	for key, val := range table {
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
			return nil, fmt.Errorf("filterPpmCompare: malformed comparison string")
		}
	}
	return result, nil
}

func validateAndDigestPpm(param []string) (float32, error) {
	if len(param) != 1 {
		return 0, fmt.Errorf("malformed ppm value in query parameters. No more than one ppm value per argument")
	}

	ppm, err := strconv.ParseFloat(param[0], 32)
	if err != nil {
		return 0, fmt.Errorf("malformed ppm value in query parameters. Cannot convert to float32")
	}
	if !(ppm <= co2PpmMax && ppm >= co2PpmMin) {
		return 0, fmt.Errorf("malformed ppm value in query parameters. Ppm range is %v to %v", co2PpmMax, co2PpmMin)
	}
	return float32(ppm), nil
}

func validateAndDigestBool(param []string) (bool, error) {
	if len(param) != 1 {
		return false, fmt.Errorf("malformed boolean value in query parameters. No more than one boolean value per argument")
	}

	result, err := strconv.ParseBool(param[0])
	if err != nil {
		return false, fmt.Errorf("malformed boolean value in query parameters. Cannot convert to boolean")
	}
	return result, nil
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
