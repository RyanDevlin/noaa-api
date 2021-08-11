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
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestGetAll(t *testing.T) {
	sqlString := regexp.QuoteMeta(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day`)
	query := "/v1/co2/weekly"
	validKeys := []string{"1974-5-19", "1974-5-26", "1984-1-1", "1984-1-8", "2000-1-2", "2000-1-9", "2018-9-2", "2018-10-7", "2020-2-2", "2020-5-24"}

	RunTest(t, t.Name(), nil, sqlString, query, validKeys)
}

func TestGetYear(t *testing.T) {
	testVal := 2020

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE year in ('%v') ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?year=%v", testVal)
	validKeys := []string{"2020-2-2", "2020-5-24"}

	RunTest(t, t.Name(), testVal, sqlString, query, validKeys)
}

func TestGetMonth(t *testing.T) {
	testVal := 1

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE month in ('%v') ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?month=%v", testVal)
	validKeys := []string{"1984-1-1", "1984-1-8", "2000-1-2", "2000-1-9"}

	RunTest(t, t.Name(), testVal, sqlString, query, validKeys)
}

func TestGetGt(t *testing.T) {
	testVal := 405.68

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average > %v ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?gt=%v", testVal)
	validKeys := []string{"2018-10-7", "2020-2-2", "2020-5-24"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validKeys)
}

func TestGetGte(t *testing.T) {
	testVal := 405.68

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average >= %v ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?gte=%v", testVal)
	validKeys := []string{"2018-9-2", "2018-10-7", "2020-2-2", "2020-5-24"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validKeys)
}

func TestGetLt(t *testing.T) {
	testVal := 344.19

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average < %v ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?lt=%v", testVal)
	validKeys := []string{"1974-5-19", "1974-5-26", "1984-1-8"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validKeys)
}

func TestGetLte(t *testing.T) {
	testVal := 344.19

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average <= %v ORDER BY year,month,day`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?lte=%v", testVal)
	validKeys := []string{"1974-5-19", "1974-5-26", "1984-1-1", "1984-1-8"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validKeys)
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
	validKeys := []string{"1984-1-8", "2000-1-2"}

	RunTest(t, t.Name(), []float32{343.89, 368.89}, sqlString, query, validKeys)
}

// TODO: Add error tests

func RunTest(t *testing.T, testName string, testVal interface{}, sqlString string, query string, validKeys []string) {
	db, mock, rows, data, err := test.NewMockCo2Db()
	if err != nil {
		t.Errorf("error generating mock database: %s", err.Error())
	}
	defer db.Close()

	for _, v := range data {
		switch testName {
		case "TestGetAll":
			// Add all entries to mock database response
			rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
		case "TestGetYear":
			// Only add 2020 entries to mock database response
			if v.Year == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetMonth":
			// Only add January entries to mock database response
			if v.Month == testVal.(int) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetGt":
			// Only add January entries to mock database response
			if v.Average > testVal.(float32) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetGte":
			// Only add January entries to mock database response
			if v.Average >= testVal.(float32) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetLt":
			// Only add January entries to mock database response
			if v.Average < testVal.(float32) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetLte":
			// Only add January entries to mock database response
			if v.Average <= testVal.(float32) {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		case "TestGetCombo":
			// Only add January entries to mock database response
			if v.Average == testVal.([]float32)[0] || v.Average == testVal.([]float32)[1] {
				rows.AddRow(v.Year, v.Month, v.Day, v.Date_decimal, v.Average, v.Ndays, v.One_year_ago, v.Ten_years_ago, v.Increase_since_1800, v.YYYYMMDD)
			}
		default:
			t.Errorf("A test named '%s' has not been implemented in the 'RunTest' function.", testName)
			return
		}
	}

	mock.ExpectQuery(sqlString).WillReturnRows(rows)

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
		test.ErrorLog(t, err)
		test.HttpJsonError(w, err)
	}

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if verbose {
		t.Logf("SQL query: '%s'", sqlString)
		test.PrintServerResponse(t, resp, body)
	}

	validateResponse(t, body, validKeys)

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
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
