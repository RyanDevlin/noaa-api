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
	"apiserver/test"
	"context"
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

	"github.com/DATA-DOG/go-sqlmock"
)

var verbose bool

func TestMain(m *testing.M) {
	flag.BoolVar(&verbose, "verbose", false, "turns on verbose logging for test cases")
	flag.Parse()
	os.Exit(m.Run())
}

/* HELPER FUNCTIONS */

func RunTest(t *testing.T, testName string, testVal interface{}, sqlString string, query string, validValues []string) {
	db, mock, rows, data, err := test.NewMockCo2Db()
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
	w := httptest.NewRecorder()

	config := &handlers.ApiHandlerConfig{
		Database: &database.Database{
			DB: db,
		},
	}

	t.Logf("Attempting request: %s %s",
		req.Method,
		req.RequestURI,
	)

	// Execute method under test
	if err := Get(ctx, config, w, req); err != nil {
		// Because the Get function returned an error, we must use HttpJsonError to write to the ResponseRecorder before checking the result.
		// Normally Get() writes to the buffer, but when encountering an error it won't have a chance to.
		test.HttpJsonError(w, err)
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

func configureDbRows(t *testing.T, testName string, testVal interface{}, rows *sqlmock.Rows, data []test.MockCo2Row) error {
	for i, v := range data {
		switch testName {
		case "TestGetAll", "TestIncreaseGetAll":
			// Add all entries to mock database response
			rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
		case "TestGetYear", "TestIncreaseGetYear":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if v.Year == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetMonth", "TestIncreaseGetMonth":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if v.Month == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetGt", "TestIncreaseGetGt":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.Increase_since_1800 > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			} else {
				if v.Average > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			}
		case "TestGetGte", "TestIncreaseGetGte":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.Increase_since_1800 >= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			} else {
				if v.Average >= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			}
		case "TestGetLt", "TestIncreaseGetLt":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.Increase_since_1800 < testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			} else {
				if v.Average < testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			}
		case "TestGetLte", "TestIncreaseGetLte":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.Increase_since_1800 <= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			} else {
				if v.Average <= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			}
		case "TestGetLimit", "TestIncreaseGetLimit":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if i < testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetOffset", "TestIncreaseGetOffset":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if i+1 > testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetPage", "TestIncreaseGetPage":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if i+1 > testVal.(int) && i+1 < 5 {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetCombo", "TestIncreaseGetCombo":
			if _, ok := testVal.([]float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type []float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.Increase_since_1800 == testVal.([]float32)[0] || v.Increase_since_1800 == testVal.([]float32)[1] {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			} else {
				if v.Average == testVal.([]float32)[0] || v.Average == testVal.([]float32)[1] {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			}
		case "TestGetNull", "TestIncreaseGetNull":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.Increase_since_1800 > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			} else {
				if v.Average > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
				}
			}
		case "TestErrors", "TestIncreaseErrors":
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
		return
	}

	for i, val := range results.Results {
		y, m, d, err := parseDate(validDates[i])
		if err != nil {
			t.Error(err)
		}

		result := val.(map[string]interface{})
		for k, v := range result {
			switch entry := v.(type) {
			case float64:
				if k == "Year" && entry != float64(y) {
					t.Errorf("Entry with year: '%v' was not present in JSON response or appeared in an incorrect order.", y)
					t.Errorf("Wanted: '%v', Got: '%v'.", y, entry)
				}
				if k == "Month" && entry != float64(m) {
					t.Errorf("Entry with month: '%v' was not present in JSON response or appeared in an incorrect order.", m)
					t.Errorf("Wanted: '%v', Got: '%v'.", m, entry)
				}
				if k == "Day" && entry != float64(d) {
					t.Errorf("Entry with day: '%v' was not present in JSON response or appeared in an incorrect order.", d)
					t.Errorf("Wanted: '%v', Got: '%v'.", d, entry)
				}
			}
		}
	}
}

func parseDate(date string) (year int, month int, day int, err error) {
	ymd := strings.Split(date, "-")
	if len(ymd) != 3 {
		return -1, -1, -1, fmt.Errorf("Improper format for date '%v'. Should be yyyy-mm-dd.", ymd)
	}

	year, err = strconv.Atoi(ymd[0])
	if err != nil {
		return -1, -1, -1, err
	}
	month, err = strconv.Atoi(ymd[1])
	if err != nil {
		return -1, -1, -1, err
	}
	day, err = strconv.Atoi(ymd[2])
	if err != nil {
		return -1, -1, -1, err
	}
	return year, month, day, nil
}

func validateErrorResponse(t *testing.T, r *http.Response, expectedCode []string) {
	if len(expectedCode) != 1 {
		t.Errorf("Exactly one http status code required for error validation. Got: %v", expectedCode)
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
