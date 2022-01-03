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

package co2

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

func configureDbRows(t *testing.T, testName string, testVal interface{}, rows *sqlmock.Rows, data []mockCo2Row) error {
	for i, v := range data {
		switch testName {
		case "TestCo2GetAll", "TestCo2IncreaseGetAll":
			// Add all entries to mock database response
			rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
		case "TestCo2GetYear", "TestCo2IncreaseGetYear":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if v.Year == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
			}
		case "TestCo2GetMonth", "TestCo2IncreaseGetMonth":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if v.Month == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
			}
		case "TestCo2GetGt", "TestCo2IncreaseGetGt":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.IncreaseSince1800 > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			} else {
				if v.Average > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			}
		case "TestCo2GetGte", "TestCo2IncreaseGetGte":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.IncreaseSince1800 >= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			} else {
				if v.Average >= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			}
		case "TestCo2GetLt", "TestCo2IncreaseGetLt":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.IncreaseSince1800 < testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			} else {
				if v.Average < testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			}
		case "TestCo2GetLte", "TestCo2IncreaseGetLte":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.IncreaseSince1800 <= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			} else {
				if v.Average <= testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			}
		case "TestCo2GetLimit", "TestCo2IncreaseGetLimit":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if i < testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
			}
		case "TestCo2GetOffset", "TestCo2IncreaseGetOffset":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if i+1 > testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
			}
		case "TestCo2GetPage", "TestCo2IncreaseGetPage":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Add desired entries to mock database response
			if i+1 > testVal.(int) && i+1 < 5 {
				rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
			}
		case "TestCo2GetCombo", "TestCo2IncreaseGetCombo":
			if _, ok := testVal.([]float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type []float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.IncreaseSince1800 == testVal.([]float32)[0] || v.IncreaseSince1800 == testVal.([]float32)[1] {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			} else {
				if v.Average == testVal.([]float32)[0] || v.Average == testVal.([]float32)[1] {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			}
		case "TestCo2GetNull", "TestCo2IncreaseGetNull":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Add desired entries to mock database response
			if strings.Contains(testName, "Increase") {
				if v.IncreaseSince1800 > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			} else {
				if v.Average > testVal.(float32) {
					rows.AddRow(v.Year, v.Month, v.Day, v.DateDecimal, v.Average, v.Ndays, v.OneYearAgo, v.TenYearsAgo, v.IncreaseSince1800, v.YYYYMMDD)
				}
			}
		case "TestCo2Errors", "TestCo2IncreaseErrors":
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
					t.Errorf("Entry with year '%v' was not present in JSON response or appeared in an incorrect order.", y)
					t.Errorf("Wanted: '%v', Got: '%v'.", y, entry)
				}
				if k == "Month" && entry != float64(m) {
					t.Errorf("Entry with month '%v' was not present in JSON response or appeared in an incorrect order.", m)
					t.Errorf("Wanted: '%v', Got: '%v'.", m, entry)
				}
				if k == "Day" && entry != float64(d) {
					t.Errorf("Entry with day '%v' was not present in JSON response or appeared in an incorrect order.", d)
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

type mockCo2Row struct {
	Year              int
	Month             int
	Day               int
	DateDecimal       float32
	Average           float32
	Ndays             int
	OneYearAgo        float32
	TenYearsAgo       float32
	IncreaseSince1800 float32
	YYYYMMDD          time.Time
}

// newMockDb returns an sqlmock database to be used for unit tests.
func newMockDb() (*sql.DB, sqlmock.Sqlmock, *sqlmock.Rows, []mockCo2Row, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	columns := []string{"year", "month", "day", "date_decimal", "average", "ndays", "one_year_ago", "ten_years_ago", "Increase_since_1800", "YYYYMMDD"}

	rows := sqlmock.NewRows(columns)

	data := GetMockCo2Rows()

	return db, mock, rows, data, nil
}

// GetMockCo2Rows returns a list hardcoded co2 measurement data used to mock the database.
func GetMockCo2Rows() []mockCo2Row {
	return []mockCo2Row{
		{
			Year:              1974,
			Month:             5,
			Day:               19,
			DateDecimal:       1974.3795,
			Average:           333.37,
			Ndays:             5,
			OneYearAgo:        -999.99,
			TenYearsAgo:       -999.99,
			IncreaseSince1800: 50.4,
			YYYYMMDD:          time.Date(1974, time.Month(5), 19, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:              1974,
			Month:             5,
			Day:               26,
			DateDecimal:       1974.3986,
			Average:           332.95,
			Ndays:             6,
			OneYearAgo:        -999.99,
			TenYearsAgo:       -999.99,
			IncreaseSince1800: 50.06,
			YYYYMMDD:          time.Date(1974, time.Month(5), 26, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:              1984,
			Month:             1,
			Day:               1,
			DateDecimal:       1984.0014,
			Average:           344.19,
			Ndays:             5,
			OneYearAgo:        341.51,
			TenYearsAgo:       -999.99,
			IncreaseSince1800: 64.53,
			YYYYMMDD:          time.Date(1984, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:              1984,
			Month:             1,
			Day:               8,
			DateDecimal:       1984.0205,
			Average:           343.89,
			Ndays:             6,
			OneYearAgo:        341.86,
			TenYearsAgo:       -999.99,
			IncreaseSince1800: 64.09,
			YYYYMMDD:          time.Date(1984, time.Month(1), 8, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:              2000,
			Month:             1,
			Day:               2,
			DateDecimal:       2000.0041,
			Average:           368.89,
			Ndays:             7,
			OneYearAgo:        367.99,
			TenYearsAgo:       353.64,
			IncreaseSince1800: 88.9,
			YYYYMMDD:          time.Date(2000, time.Month(1), 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:              2000,
			Month:             1,
			Day:               9,
			DateDecimal:       2000.0232,
			Average:           369.03,
			Ndays:             7,
			OneYearAgo:        368.23,
			TenYearsAgo:       353.63,
			IncreaseSince1800: 88.88,
			YYYYMMDD:          time.Date(2000, time.Month(1), 9, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:              2018,
			Month:             9,
			Day:               2,
			DateDecimal:       2018.6699,
			Average:           405.68,
			Ndays:             7,
			OneYearAgo:        404.11,
			TenYearsAgo:       383.72,
			IncreaseSince1800: 128.89,
			YYYYMMDD:          time.Date(2018, time.Month(9), 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:              2018,
			Month:             10,
			Day:               7,
			DateDecimal:       2018.7658,
			Average:           405.77,
			Ndays:             7,
			OneYearAgo:        403.58,
			TenYearsAgo:       383.01,
			IncreaseSince1800: 129.42,
			YYYYMMDD:          time.Date(2018, time.Month(10), 7, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:              2020,
			Month:             2,
			Day:               2,
			DateDecimal:       2020.0888,
			Average:           414.53,
			Ndays:             7,
			OneYearAgo:        411.31,
			TenYearsAgo:       390.87,
			IncreaseSince1800: 133.87,
			YYYYMMDD:          time.Date(2020, time.Month(2), 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Year:              2020,
			Month:             5,
			Day:               24,
			DateDecimal:       2020.3948,
			Average:           417.67,
			Ndays:             7,
			OneYearAgo:        414.62,
			TenYearsAgo:       392.85,
			IncreaseSince1800: 134.36,
			YYYYMMDD:          time.Date(2020, time.Month(5), 24, 0, 0, 0, 0, time.UTC),
		},
	}
}
