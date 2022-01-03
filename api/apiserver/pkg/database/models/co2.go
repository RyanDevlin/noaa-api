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

package models

import (
	"database/sql"
	"time"
)

const (
	// Co2PpmMax is the maximum ppm value that may be used in a query for Co2 data
	Co2PpmMax = 1000

	// Co2PpmMin is the minimum ppm value that may be used in a query for Co2 data
	Co2PpmMin = 0
)

// Co2Table represents a list of Co2Entry objects
type Co2Table []interface{}

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
	Year                  int
	Month                 int
	Day                   int
	Average               float32
	IncSincePreIndustrial float32
}

// Load imports the results of a database query into a Co2Table slice
func (co2Table *Co2Table) Load(rows *sql.Rows, simple bool) error {
	if !simple {
		var co2entry Co2Entry
		if err := rows.Scan(&co2entry.Year, &co2entry.Month, &co2entry.Day, &co2entry.DateDecimal, &co2entry.Average, &co2entry.NumDays, &co2entry.OneYearAgo, &co2entry.TenYearsAgo, &co2entry.IncSincePreIndustrial, &co2entry.Timestamp); err != nil {
			return err
		}
		*co2Table = append(*co2Table, co2entry)
	} else {
		var co2entry Co2EntrySimple
		if err := rows.Scan(&co2entry.Year, &co2entry.Month, &co2entry.Day, &co2entry.Average, &co2entry.IncSincePreIndustrial); err != nil {
			return err
		}
		*co2Table = append(*co2Table, co2entry)
	}
	return nil
}
