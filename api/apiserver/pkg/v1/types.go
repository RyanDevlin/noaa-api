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

package v1

// A QueryParam is any parameter that can be passed as a query.
// The parameter is represented as a string key mapped to an interface.
// This is used to model all params such as 'test=true' or 'nums=[2,4]'.
type QueryParam struct {
	Key   string
	Value []string
}

// Implementing a QueryFilter allows for reducing a map of data from the DB
// down to a smaller output. Filter functions ingest a QueryParam and data
// from the DB represented as map[string]interface{}. The logic in the Filter
// method is then used to discard entries from the given data and return data
// that is less than or equal to the original input.
type QueryFilter interface {
	Filter(map[string]interface{}) (map[string]interface{}, error)
}

/*
// QueryParams acts as a mapping from
type QueryParams struct {
	Params map[string]QueryFilter
}*/
