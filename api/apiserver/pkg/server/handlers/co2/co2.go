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
	"apiserver/pkg/server/handlers"
	"apiserver/pkg/utils"
	"context"
	"encoding/json"
	"net/http"
)

// Get is an ApiHandlerFunc type. It queries the database for requested co2weekly data and returns a JSON representation of the data
// to the client.
func Get(ctx context.Context, handlerConfig *handlers.ApiHandlerConfig, w http.ResponseWriter, r *http.Request) *utils.ServerError {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	query := database.NewQuery("public.co2_weekly_mlo", []string{"*"}, "year,month,day")

	filters, internalArgs, err := ParseParams(r, handlerConfig.PathParam, handlerConfig.SortBy)
	if err != nil {
		return err
	}

	if len(internalArgs) != 0 {
		ParseInternalArgs(internalArgs, &query)
	}

	query.Where = filters

	co2Table := models.Co2Table{}
	dberr := handlerConfig.Database.Query(query, &co2Table)
	if dberr != nil {
		return utils.NewError(dberr, "internal database error", 500, false)
	}

	// This prevents the 'Results' part of the response from being omitted if
	// there are no results.
	if len(co2Table) == 0 {
		co2Table = models.Co2Table{
			nil,
		}
	}

	// Parse RequestID param
	id, idError := utils.GetReqId(r)
	if idError != nil {
		return utils.NewError(idError, "cannot extract request ID", 500, false)
	}

	resp := models.ServerResp{
		Results:   co2Table,
		Status:    "OK",
		RequestId: id,
		Error:     nil,
	}

	enc := json.NewEncoder(w)
	if query.Pretty {
		enc.SetIndent("", "    ")
	}
	if err := enc.Encode(resp); err != nil {
		return utils.NewError(err, "error encoding data as json", 500, false)
	}
	return nil
}
