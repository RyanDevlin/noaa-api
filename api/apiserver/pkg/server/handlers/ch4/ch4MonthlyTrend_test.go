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
	"fmt"
	"regexp"
	"testing"
)

func TestCh4IncreaseGetAll(t *testing.T) {
	sqlString := regexp.QuoteMeta(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT 10`)
	query := "/v1/co2/weekly/increase"
	validDates := []string{"1974-05-19", "1974-05-26", "1984-01-01", "1984-01-08", "2000-01-02", "2000-01-09", "2018-09-02", "2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), nil, sqlString, query, validDates)
}

func TestCh4IncreaseGetYear(t *testing.T) {
	testVal := 2020

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE year in ('%v') ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?year=%v", testVal)
	validDates := []string{"2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates)
}

func TestCh4IncreaseGetMonth(t *testing.T) {
	testVal := 1

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE month in ('%v') ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?month=%v", testVal)
	validDates := []string{"1984-01-01", "1984-01-08", "2000-01-02", "2000-01-09"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates)
}

func TestCh4IncreaseGetGt(t *testing.T) {
	testVal := 128.89

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE increase_since_1800 > %v ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?gt=%v", testVal)
	validDates := []string{"2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates)
}

func TestCh4IncreaseGetGte(t *testing.T) {
	testVal := 128.89

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE increase_since_1800 >= %v ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?gte=%v", testVal)
	validDates := []string{"2018-09-02", "2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates)
}

func TestCh4IncreaseGetLt(t *testing.T) {
	testVal := 64.53

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE increase_since_1800 < %v ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?lt=%v", testVal)
	validDates := []string{"1974-05-19", "1974-05-26", "1984-01-08"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates)
}

func TestCh4IncreaseGetLte(t *testing.T) {
	testVal := 64.53

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE increase_since_1800 <= %v ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?lte=%v", testVal)
	validDates := []string{"1974-05-19", "1974-05-26", "1984-01-01", "1984-01-08"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates)
}

func TestCh4IncreaseGetLimit(t *testing.T) {
	testVal := 2

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT %v`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?limit=%v", testVal)
	validDates := []string{"1974-05-19", "1974-05-26"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates)
}

func TestCh4IncreaseGetOffset(t *testing.T) {
	testVal := 4

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT 10 OFFSET %v`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?offset=%v", testVal)
	validDates := []string{"2000-01-02", "2000-01-09", "2018-09-02", "2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates)
}

func TestCh4IncreaseGetPage(t *testing.T) {
	page := 2
	limit := 2

	offset := (limit * (page - 1))

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT %v OFFSET %v`, limit, offset))
	query := fmt.Sprintf("/v1/co2/weekly/increase?limit=%v&page=%v", limit, page)
	validDates := []string{"1984-01-01", "1984-01-08"}

	RunTest(t, t.Name(), offset, sqlString, query, validDates)
}

func TestCh4IncreaseGetCombo(t *testing.T) {
	years := []int{1984, 2000}
	month := 1
	gt := 64.09
	gte := 64.53
	lt := 88.9
	lte := 88.88

	// This regex will match the SELECT query with any arbitrary ordering of the WHERE clauses. This is needed because the order that the server concatenates WHERE clauses is semi-random
	sqlString := `SELECT \* FROM public\.co2_weekly_mlo WHERE (increase_since_1800 [<>=]+ [\d\.]+( AND )*|year in \(('[\d]+'(,)?[ ]*)*\)( AND )*|month in \(('[\d]+'(,)?[ ]*)*\)( AND )*)* ORDER BY year,month,day LIMIT 10`
	query := fmt.Sprintf("/v1/co2/weekly/increase?year=%v,%v&month=%v&gt=%v&gte=%v&lt=%v&lte=%v", years[0], years[1], month, gt, gte, lt, lte)
	validDates := []string{"1984-01-01", "2000-01-09"}

	RunTest(t, t.Name(), []float32{64.53, 88.88}, sqlString, query, validDates)
}

func TestCh4IncreaseGetNull(t *testing.T) {
	testVal := 500.00

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE increase_since_1800 > %.2f ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?gt=%v", testVal)
	validValues := []string{}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validValues)
}

func TestCh4IncreaseErrors(t *testing.T) {
	testVals := []string{
		"/v1/co2/weekly/increase?year=2020a",
		"/v1/co2/weekly/increase?year=20200",
		"/v1/co2/weekly/increase?month=1a",
		"/v1/co2/weekly/increase?month=14",
		"/v1/co2/weekly/increase?gt=400a",
		"/v1/co2/weekly/increase?lt=400a",
		"/v1/co2/weekly/increase?gte=400a",
		"/v1/co2/weekly/increase?lte=400a",
		"/v1/co2/weekly/increase?gt=40000",
		"/v1/co2/weekly/increase?gt=-1",
	}

	sqlString := ``
	validValues := []string{"400"} // The http response code we're expecting

	for _, v := range testVals {
		query := fmt.Sprintf("%v", v)
		RunTest(t, t.Name(), nil, sqlString, query, validValues)
	}
}
