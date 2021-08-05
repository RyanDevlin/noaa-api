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
	utils "apiserver/pkg/utils"
	"apiserver/pkg/v1/co2Weekly"
	"encoding/json"
	"net/http"
)

func co2WeeklyHandlerFactory(apiserver *ApiServer) utils.ApiHandler {
	return utils.ApiHandler(func(w http.ResponseWriter, r *http.Request) *utils.ServerError {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		co2table, err := apiserver.DBGetCo2Table()
		if err != nil {
			return utils.NewError(err, "failed to connect to database", 500, false)
		}

		// Filter data based on query params
		filters := co2FilterFactory(r)
		for _, filter := range filters {
			co2table, err = filter(co2table)
			if err != nil {
				message := err.Error() + ": " + utils.ParseQuery(r).Encode()
				return utils.NewError(err, message, 400, false)
			}
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		if err := enc.Encode(co2table); err != nil {
			return utils.NewError(err, "error encoding data as json", 500, false)
		}
		return nil
	})
}

func co2FilterFactory(r *http.Request) []co2Weekly.FilterFunc {
	params := utils.ParseQuery(r)

	var filters []co2Weekly.FilterFunc
	for key, val := range params {
		switch key {
		case "year":
			co2Year := &co2Weekly.Co2Year{
				Params: val,
			}
			filters = append(filters, co2Year.Filter)

		case "month":
			co2Month := &co2Weekly.Co2Month{
				Params: val,
			}
			filters = append(filters, co2Month.Filter)
		case "gt":
			co2GreaterThan := &co2Weekly.Co2GreaterThan{
				Params: val,
			}
			filters = append(filters, co2GreaterThan.Filter)
		case "lt":
			co2LessThan := &co2Weekly.Co2LessThan{
				Params: val,
			}
			filters = append(filters, co2LessThan.Filter)
		case "gte":
			co2Gte := &co2Weekly.Co2Gte{
				Params: val,
			}
			filters = append(filters, co2Gte.Filter)
		case "lte":
			co2Lte := &co2Weekly.Co2Lte{
				Params: val,
			}
			filters = append(filters, co2Lte.Filter)
		case "simple":
			co2Simple := &co2Weekly.Co2Simple{
				Params: val,
			}
			filters = append(filters, co2Simple.Filter)
		}
	}
	return filters
}
