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

package co2weekly

import (
	"apiserver/pkg/database"
	"apiserver/pkg/server/handlers"
	"apiserver/test"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAll(t *testing.T) {
	sqlString := regexp.QuoteMeta(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day`)
	query := "/v1/co2/weekly"
	validKeys := []string{"1974-05-19", "1974-05-26", "1984-01-01", "1984-01-08", "2000-01-02", "2000-01-09", "2018-09-02", "2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), nil, sqlString, query, validKeys)
}

func TestGetYear(t *testing.T) {
	testVal := 2020

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE year in ('%v') ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?year=%v", testVal)
	validKeys := []string{"2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), testVal, sqlString, query, validKeys)
}

func TestGetMonth(t *testing.T) {
	testVal := 1

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE month in ('%v') ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?month=%v", testVal)
	validKeys := []string{"1984-01-01", "1984-01-08", "2000-01-02", "2000-01-09"}

	RunTest(t, t.Name(), testVal, sqlString, query, validKeys)
}

func TestGetGt(t *testing.T) {
	testVal := 405.68

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average > %v ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?gt=%v", testVal)
	validKeys := []string{"2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validKeys)
}

func TestGetGte(t *testing.T) {
	testVal := 405.68

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average >= %v ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?gte=%v", testVal)
	validKeys := []string{"2018-09-02", "2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validKeys)
}

func TestGetLt(t *testing.T) {
	testVal := 344.19

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average < %v ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?lt=%v", testVal)
	validKeys := []string{"1974-05-19", "1974-05-26", "1984-01-08"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validKeys)
}

func TestGetLte(t *testing.T) {
	testVal := 344.19

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average <= %v ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?lte=%v", testVal)
	validKeys := []string{"1974-05-19", "1974-05-26", "1984-01-01", "1984-01-08"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validKeys)
}

func TestGetLimit(t *testing.T) {
	testVal := 2

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT %v`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?limit=%v", testVal)
	validKeys := []string{"1974-05-19", "1974-05-26"}

	RunTest(t, t.Name(), testVal, sqlString, query, validKeys)
}

func TestGetOffset(t *testing.T) {
	testVal := 4

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day OFFSET %v`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?offset=%v", testVal)
	validKeys := []string{"2000-01-02", "2000-01-09", "2018-09-02", "2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), testVal, sqlString, query, validKeys)
}

func TestGetPage(t *testing.T) {
	page := 2
	limit := 2

	//offset := (limit * (page - 1))
	offset := (limit * (page))

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT %v OFFSET %v`, limit, offset))
	query := fmt.Sprintf("/v1/co2/weekly?limit=%v&page=%v", limit, page)
	validKeys := []string{"1984-01-01", "1984-01-08"}

	RunTest(t, t.Name(), offset, sqlString, query, validKeys)
}

func TestGetCombo(t *testing.T) {
	years := []int{1984, 2000}
	month := 1
	gt := 344.19
	gte := 343.89
	lt := 369.03
	lte := 368.89

	// This regex will match the SELECT query with any arbitrary ordering of the WHERE clauses. This is needed because the order that the server concatenates WHERE clauses is semi-random
	sqlString := `SELECT \* FROM public\.co2_weekly_mlo WHERE (average [<>=]+ [\d\.]+( AND )*|year in \(('[\d]+'(,)?[ ]*)*\)( AND )*|month in \(('[\d]+'(,)?[ ]*)*\)( AND )*)* ORDER BY year,month,day`
	query := fmt.Sprintf("/v1/co2/weekly?year=%v,%v&month=%v&gt=%v&gte=%v&lt=%v&lte=%v", years[0], years[1], month, gt, gte, lt, lte)
	validKeys := []string{"1984-01-08", "2000-01-02"}

	RunTest(t, t.Name(), []float32{343.89, 368.89}, sqlString, query, validKeys)
}

func TestErrors(t *testing.T) {
	testVals := []string{
		"/v1/co2/weekly?year=2020a",
		"/v1/co2/weekly?year=20200",
		"/v1/co2/weekly?month=1a",
		"/v1/co2/weekly?month=14",
		"/v1/co2/weekly?gt=400a",
		"/v1/co2/weekly?lt=400a",
		"/v1/co2/weekly?gte=400a",
		"/v1/co2/weekly?lte=400a",
		"/v1/co2/weekly?gt=300,400",
		"/v1/co2/weekly?lt=300,400",
		"/v1/co2/weekly?lte=300,400",
		"/v1/co2/weekly?gt=300&gt=400",
		"/v1/co2/weekly?gt=40000",
		"/v1/co2/weekly?gt=-1",
	}

	sqlString := ``
	validValues := []string{"400"} // The http response code we're expecting

	for _, v := range testVals {
		query := fmt.Sprintf("%v", v)
		RunTest(t, t.Name(), nil, sqlString, query, validValues)
	}
}

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
		case "TestGetAll":
			// Add all entries to mock database response
			rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
		case "TestGetYear":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Only add 2020 entries to mock database response
			if v.Year == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetMonth":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Only add January entries to mock database response
			if v.Month == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetGt":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Only add January entries to mock database response
			if v.Average > testVal.(float32) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetGte":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Only add January entries to mock database response
			if v.Average >= testVal.(float32) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetLt":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Only add January entries to mock database response
			if v.Average < testVal.(float32) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetLte":
			if _, ok := testVal.(float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type float32.", testVal, testName)
			}

			// Only add January entries to mock database response
			if v.Average <= testVal.(float32) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetLimit":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Only add January entries to mock database response
			if i < testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetOffset":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Only add January entries to mock database response
			if i+1 > testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetPage":
			if _, ok := testVal.(int); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type int.", testVal, testName)
			}

			// Only add January entries to mock database response
			if i+1 > testVal.(int) && i+1 < 5 {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetCombo":
			if _, ok := testVal.([]float32); !ok {
				return fmt.Errorf("Test value '%v' for test '%v' is not of type []float32.", testVal, testName)
			}

			// Only add January entries to mock database response
			if v.Average == testVal.([]float32)[0] || v.Average == testVal.([]float32)[1] {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestErrors":
			if testVal != nil {
				return fmt.Errorf("Test value '%v' for test '%v' is not nil.", testVal, testName)
			}
		default:
			return fmt.Errorf("A test named '%s' has not been implemented in the 'RunTest' function.", testName)
		}
	}
	return nil
}

func validateResponse(t *testing.T, body []byte, validKeys []string) {
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)

	if len(result) != len(validKeys) {
		t.Error("Incorrect number of values returned from query.")
	}

	for _, v := range validKeys {
		if _, ok := result[v]; !ok {
			t.Errorf("Entry with key: '%s' was not present in JSON response.", v)
		}
	}
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
