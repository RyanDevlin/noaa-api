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
	"apiserver/pkg/utils"
	"net/http"
	"time"
)

const (
	co2PpmMax = 1000
	co2PpmMin = 0
)

type NoaaData interface {
	Filter(r *http.Request) *utils.ServerError
}

// The index of the Co2Table map must be '<year>-<month>-<day>'
type Co2Table map[string]interface{}

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

type Co2EntrySimple struct {
	Average               float32
	IncSincePreIndustrial float32
}
