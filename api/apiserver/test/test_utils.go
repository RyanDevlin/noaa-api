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

package test

import (
	"bytes"
	"database/sql"
	"net/http"
	"testing"
	"time"

	// The blank import here is used to import the pq PostgreSQL drivers
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

type mockCo2Row struct {
	Year                int
	Month               int
	Day                 int
	Date_decimal        float32
	Average             float32
	Ndays               int
	One_year_ago        float32
	Ten_years_ago       float32
	Increase_since_1800 float32
	YYYYMMDD            time.Time
}

// NewMockCo2Db returns an sqlmock database to be used for unit tests.
func NewMockCo2Db() (*sql.DB, sqlmock.Sqlmock, *sqlmock.Rows, []mockCo2Row, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	columns := []string{"year", "month", "day", "date_decimal", "average", "ndays", "one_year_ago", "ten_years_ago", "Increase_since_1800", "YYYYMMDD"}

	rows := sqlmock.NewRows(columns)

	data := GetMockCo2Rows()

	return db, mock, rows, data, nil
}

func PrintServerResponse(t *testing.T, r *http.Response, body []byte) {
	t.Logf("Server Response: \n\tStatus Code: %v\n\tHeader Content-Type: %v\n\tBody: \n%v\n", r.StatusCode, r.Header.Get("Content-Type"), string(indent([]byte{'\t'}, body)))
}

// indent was borred directly from https://pkg.go.dev/github.com/openconfig/goyang/pkg/indent
func indent(indent, b []byte) []byte {
	if len(indent) == 0 || len(b) == 0 {
		return b
	}
	lines := bytes.SplitAfter(b, []byte{'\n'})
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return bytes.Join(append([][]byte{{}}, lines...), indent)
}

// GetMockCo2Rows returns a list hardcoded co2 measurement data used to mock the database.
func GetMockCo2Rows() []mockCo2Row {
	return []mockCo2Row{
		{
			Year:                1974,
			Month:               5,
			Day:                 19,
			Date_decimal:        1974.3795,
			Average:             333.37,
			Ndays:               5,
			One_year_ago:        -999.99,
			Ten_years_ago:       -999.99,
			Increase_since_1800: 50.4,
			YYYYMMDD:            time.Date(1974, time.Month(5), 19, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:                1974,
			Month:               5,
			Day:                 26,
			Date_decimal:        1974.3986,
			Average:             332.95,
			Ndays:               6,
			One_year_ago:        -999.99,
			Ten_years_ago:       -999.99,
			Increase_since_1800: 50.06,
			YYYYMMDD:            time.Date(1974, time.Month(5), 26, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:                1984,
			Month:               1,
			Day:                 1,
			Date_decimal:        1984.0014,
			Average:             344.19,
			Ndays:               5,
			One_year_ago:        341.51,
			Ten_years_ago:       -999.99,
			Increase_since_1800: 64.53,
			YYYYMMDD:            time.Date(1984, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:                1984,
			Month:               1,
			Day:                 8,
			Date_decimal:        1984.0205,
			Average:             343.89,
			Ndays:               6,
			One_year_ago:        341.86,
			Ten_years_ago:       -999.99,
			Increase_since_1800: 64.09,
			YYYYMMDD:            time.Date(1984, time.Month(1), 8, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:                2000,
			Month:               1,
			Day:                 2,
			Date_decimal:        2000.0041,
			Average:             368.89,
			Ndays:               7,
			One_year_ago:        367.99,
			Ten_years_ago:       353.64,
			Increase_since_1800: 88.9,
			YYYYMMDD:            time.Date(2000, time.Month(1), 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:                2000,
			Month:               1,
			Day:                 9,
			Date_decimal:        2000.0232,
			Average:             369.03,
			Ndays:               7,
			One_year_ago:        368.23,
			Ten_years_ago:       353.63,
			Increase_since_1800: 88.88,
			YYYYMMDD:            time.Date(2000, time.Month(1), 9, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:                2018,
			Month:               9,
			Day:                 2,
			Date_decimal:        2018.6699,
			Average:             405.68,
			Ndays:               7,
			One_year_ago:        404.11,
			Ten_years_ago:       383.72,
			Increase_since_1800: 128.89,
			YYYYMMDD:            time.Date(2018, time.Month(9), 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:                2018,
			Month:               10,
			Day:                 7,
			Date_decimal:        2018.7658,
			Average:             405.77,
			Ndays:               7,
			One_year_ago:        403.58,
			Ten_years_ago:       383.01,
			Increase_since_1800: 129.42,
			YYYYMMDD:            time.Date(2018, time.Month(10), 7, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:                2020,
			Month:               2,
			Day:                 2,
			Date_decimal:        2020.0888,
			Average:             414.53,
			Ndays:               7,
			One_year_ago:        411.31,
			Ten_years_ago:       390.87,
			Increase_since_1800: 133.87,
			YYYYMMDD:            time.Date(2020, time.Month(2), 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:                2020,
			Month:               5,
			Day:                 24,
			Date_decimal:        2020.3948,
			Average:             417.67,
			Ndays:               7,
			One_year_ago:        414.62,
			Ten_years_ago:       392.85,
			Increase_since_1800: 134.36,
			YYYYMMDD:            time.Date(2020, time.Month(5), 24, 0, 0, 0, 0, time.UTC),
		},
	}
}
