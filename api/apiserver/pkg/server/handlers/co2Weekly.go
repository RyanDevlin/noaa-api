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

package handlers

import (
	utils "apiserver/pkg/utils"
	"context"
	"encoding/json"
	"net/http"
)

func GetCo2Weekly(ctx context.Context, handlerConfig *ApiHandlerConfig, w http.ResponseWriter, r *http.Request) *utils.ServerError {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	co2Table, err := handlerConfig.Database.Query("SELECT * FROM public.co2_weekly_mlo")
	if err != nil {
		return utils.NewError(err, "failed to connect to database", 500, false)
	}

	// Filter data based on query params
	if err := co2Table.Filter(r); err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(co2Table); err != nil {
		return utils.NewError(err, "error encoding data as json", 500, false)
	}
	return nil
}
