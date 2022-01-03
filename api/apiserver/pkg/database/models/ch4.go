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
	// Ch4PpbMax is the maximum ppb value that may be used in a query for Ch4 data
	Ch4PpbMax = 3000

	// Ch4PpbMin is the minimum ppb value that may be used in a query for Ch4 data
	Ch4PpbMin = 0
)

// Ch4Table represents a list of Ch4Entry objects
type Ch4Table []interface{}

// Ch4Entry represents the JSON data to be returned from an individual Ch4 measurement in the database.
type Ch4Entry struct {
	Year               int
	Month              int
	DateDecimal        float32
	Average            float32
	AverageUncertainty float32
	Trend              float32
	TrendUncertainty   float32
	Timestamp          time.Time
}

// Ch4EntrySimple represents the simplified JSON data to be returned from an individual Ch4 measurement in the database.
type Ch4EntrySimple struct {
	Year    int
	Month   int
	Average float32
	Trend   float32
}

// Load imports the results of a database query into a Ch4Table slice
func (ch4Table *Ch4Table) Load(rows *sql.Rows, simple bool) error {
	if !simple {
		var ch4entry Ch4Entry
		if err := rows.Scan(&ch4entry.Year, &ch4entry.Month, &ch4entry.DateDecimal, &ch4entry.Average, &ch4entry.AverageUncertainty, &ch4entry.Trend, &ch4entry.TrendUncertainty, &ch4entry.Timestamp); err != nil {
			return err
		}
		*ch4Table = append(*ch4Table, ch4entry)
	} else {
		var ch4entry Ch4EntrySimple
		if err := rows.Scan(&ch4entry.Year, &ch4entry.Month, &ch4entry.Average, &ch4entry.Trend); err != nil {
			return err
		}
		*ch4Table = append(*ch4Table, ch4entry)
	}
	return nil
}
