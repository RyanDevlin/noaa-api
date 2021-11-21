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

package ch4

import (
	"apiserver/pkg/database"
	"apiserver/pkg/database/models"
	"apiserver/pkg/server/handlers"
	"apiserver/pkg/utils"
	"apiserver/test"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var verbose bool

func TestMain(m *testing.M) {
	flag.BoolVar(&verbose, "verbose", false, "turns on verbose logging for test cases")
	flag.Parse()
	os.Exit(m.Run())
}

/* HELPER FUNCTIONS */

func RunTest(t *testing.T, testName string, testVal interface{}, sqlString string, query string, validValues []string, config *handlers.ApiHandlerConfig) {
	db, mock, rows, data, err := newMockDb()
	if err != nil {
		t.Errorf("error generating mock database: %s", err.Error())
		return
	}
	defer db.Close()

	mock.ExpectQuery(sqlString).WillReturnRows(rows)

	err = configureDbRows(t, testName, testVal, rows, data)
	if err != nil {
		t.Error(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest("GET", query, nil)
	req = test.SetReqIdTest(req) // Attach a mock UUID to this test request
	w := httptest.NewRecorder()

	t.Logf("Attempting request: %s %s",
		req.Method,
		req.RequestURI,
	)

	config.Database = &database.Database{
		DB: db,
	}

	// Execute method under test
	if err := Get(ctx, config, w, req); err != nil {
		// Because the Get function returned an error, we must use HttpJsonError to write to the ResponseRecorder before checking the result.
		// Normally Get() writes to the io buffer, but when encountering an error it won't have a chance to.
		utils.HttpJsonError(w, req, err)
		resp := w.Result()

		if verbose {
			test.ErrorLog(t, err)

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			test.PrintServerResponse(t, resp, body)
		}

		if strings.Contains(err.Error.Error(), "could not match actual sql") {
			// This can occur with a badly written test case. Usually if the SQL query
			// regex in the test case does not match what is actually used by the server.
			test.ErrorLog(t, err)
			t.Error("Test failed. Is the 'sqlString' regex correct?")
			return
		}

		validateErrorResponse(t, resp, validValues)
		return
	}

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	if verbose {
		t.Logf("SQL query: '%s'", sqlString)
		test.PrintServerResponse(t, resp, body)
	}

	validateResponse(t, body, validValues)

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func configureDbRows(t *testing.T, testName string, testVal interface{}, rows *sqlmock.Rows, data []mockCh4Row) error {
	for i, v := range data {
		switch testName {
		case "TestCh4GetAll", "TestCh4TrendGetAll":
			// Add all entries to mock database response
			rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
		case "TestCh4GetYear", "TestCh4TrendGetYear":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if v.Year == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
			}
		case "TestCh4GetMonth", "TestCh4TrendGetMonth":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if v.Month == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
			}
		case "TestCh4GetGt", "TestCh4TrendGetGt":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Trend") {
				if v.Trend > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			} else {
				if v.Average > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			}
		case "TestCh4GetGte", "TestCh4TrendGetGte":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Trend") {
				if v.Trend >= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			} else {
				if v.Average >= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			}
		case "TestCh4GetLt", "TestCh4TrendGetLt":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Trend") {
				if v.Trend < testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			} else {
				if v.Average < testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			}
		case "TestCh4GetLte", "TestCh4TrendGetLte":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Trend") {
				if v.Trend <= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			} else {
				if v.Average <= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			}
		case "TestCh4GetLimit", "TestCh4TrendGetLimit":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if i < testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
			}
		case "TestCh4GetOffset", "TestCh4TrendGetOffset":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if i+1 > testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
			}
		case "TestCh4GetPage", "TestCh4TrendGetPage":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if i+1 > testVal.(int) && i+1 < 5 {
				rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
			}
		case "TestCh4GetCombo", "TestCh4TrendGetCombo":
			if _, ok := testVal.([]float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type []float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Trend") {
				if fmt.Sprintf("%.2f", v.Trend) == fmt.Sprintf("%.2f", testVal.([]float32)[0]) || fmt.Sprintf("%.2f", v.Trend) == fmt.Sprintf("%.2f", testVal.([]float32)[1]) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			} else {
				if fmt.Sprintf("%.2f", v.Average) == fmt.Sprintf("%.2f", testVal.([]float32)[0]) || fmt.Sprintf("%.2f", v.Average) == fmt.Sprintf("%.2f", testVal.([]float32)[1]) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			}
		case "TestCh4GetNull", "TestCh4TrendGetNull":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Trend") {
				if v.Trend < testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			} else {
				if v.Average < testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.DateDecimal, v.Average, v.AverageUncertainty, v.Trend, v.TrendUncertainty, v.Timestamp)
				}
			}
		case "TestCh4Errors", "TestCh4TrendErrors":
			if testVal != nil {
				return fmt.Errorf("Test value '%v' for test '%v' is not nil.", testVal, testName)
			}
		default:
			return fmt.Errorf("A test named '%s' has not been implemented in the 'RunTest' function.", testName)
		}
	}
	return nil
}

func validateResponse(t *testing.T, body []byte, validDates []string) {
	results := models.ServerResp{}
	json.Unmarshal(body, &results)

	if len(results.Results) != len(validDates) {
		// If we expect there to be no results, the test case should initialize validDates as an empty slice.
		if len(validDates) == 0 && results.Results[0] == nil {
			return
		}
		t.Error("Incorrect number of values returned from query.")
		t.Errorf("Expected '%v' values, instead got '%v' values.", len(validDates), len(results.Results))
		return
	}

	for i, val := range results.Results {
		dateDecimal, err := strconv.ParseFloat(validDates[i], 32)
		if err != nil {
			t.Error(err)
		}

		result := val.(map[string]interface{})
		for k, v := range result {
			switch entry := v.(type) {
			case float64:
				if k == "DateDecimal" && fmt.Sprintf("%.2f", entry) != fmt.Sprintf("%.2f", dateDecimal) {
					t.Errorf("Entry date with decimal representation '%v' was not present in JSON response or appeared in an incorrect order.", dateDecimal)
					t.Errorf("Wanted: '%v', Got: '%v'.", dateDecimal, entry)
				}
			}
		}
	}
}

func validateErrorResponse(t *testing.T, r *http.Response, expectedCode []string) {
	if len(expectedCode) != 1 {
		t.Errorf("Exactly one http status code required for error validation. Got: %v", expectedCode)
		t.Errorf("Should this test case return a database error?")
		return
	}

	val, err := strconv.Atoi(expectedCode[0])
	if err != nil {
		t.Errorf("Problem converting http status code string to int: %v", err)
		return
	}

	if val < 100 || val > 599 {
		t.Errorf("Invalid http status code value '%v'. Valid http status codes range from [100-599]", val)
		return
	}

	if r.StatusCode != val {
		t.Errorf("Response status code '%v' does not match expected code '%v'.", r.StatusCode, val)
		return
	}
	t.Logf("Response code '%v: %v'. This is correct.", val, http.StatusText(val))
}

type mockCh4Row struct {
	Year               int
	Month              int
	DateDecimal        float32
	Average            float32
	AverageUncertainty float32
	Trend              float32
	TrendUncertainty   float32
	Timestamp          time.Time
}

// newMockDb returns an sqlmock database to be used for unit tests.
func newMockDb() (*sql.DB, sqlmock.Sqlmock, *sqlmock.Rows, []mockCh4Row, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	columns := []string{"year", "month", "date_decimal", "average", "average_unc", "trend", "trend_unc", "YYYYMMDD"}

	rows := sqlmock.NewRows(columns)

	data := GetMockCh4Rows()

	return db, mock, rows, data, nil
}

// GetMockCh4Rows returns a list hardcoded Ch4 measurement data used to mock the database.
func GetMockCh4Rows() []mockCh4Row {
	return []mockCh4Row{
		{
			Year:               1983,
			Month:              7,
			DateDecimal:        1983.542,
			Average:            1625.4,
			AverageUncertainty: 2.4,
			Trend:              1634.5,
			TrendUncertainty:   1.5,
			Timestamp:          time.Date(1983, time.Month(7), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:               1983,
			Month:              8,
			DateDecimal:        1983.625,
			Average:            1627.5,
			AverageUncertainty: 2.9,
			Trend:              1635.1,
			TrendUncertainty:   1.4,
			Timestamp:          time.Date(1983, time.Month(8), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:               1990,
			Month:              1,
			DateDecimal:        1990.042,
			Average:            1712.1,
			AverageUncertainty: 1.2,
			Trend:              1710.4,
			TrendUncertainty:   0.6,
			Timestamp:          time.Date(1990, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:               1990,
			Month:              2,
			DateDecimal:        1990.125,
			Average:            1713.5,
			AverageUncertainty: 1.3,
			Trend:              1711.1,
			TrendUncertainty:   0.6,
			Timestamp:          time.Date(1990, time.Month(2), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:               2000,
			Month:              1,
			DateDecimal:        2000.042,
			Average:            1776.1,
			AverageUncertainty: 1.1,
			Trend:              1773.5,
			TrendUncertainty:   0.7,
			Timestamp:          time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:               2000,
			Month:              2,
			DateDecimal:        2000.125,
			Average:            1776,
			AverageUncertainty: 1.4,
			Trend:              1773.4,
			TrendUncertainty:   0.7,
			Timestamp:          time.Date(2000, time.Month(2), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:               2020,
			Month:              10,
			DateDecimal:        2020.792,
			Average:            1890.1,
			AverageUncertainty: -9.9,
			Trend:              1883.9,
			TrendUncertainty:   -9.9,
			Timestamp:          time.Date(1983, time.Month(10), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:               2020,
			Month:              11,
			DateDecimal:        2020.875,
			Average:            1891.7,
			AverageUncertainty: -9.9,
			Trend:              1885,
			TrendUncertainty:   -9.9,
			Timestamp:          time.Date(1983, time.Month(11), 1, 0, 0, 0, 0, time.UTC),
		},
	}
}
