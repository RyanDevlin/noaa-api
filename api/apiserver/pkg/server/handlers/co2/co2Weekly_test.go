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
	"apiserver/pkg/server/handlers"
	"fmt"
	"regexp"
	"testing"
)

var handlerConfig = &handlers.ApiHandlerConfig{
	SortBy: "average",
}

func TestCo2GetAll(t *testing.T) {
	sqlString := regexp.QuoteMeta(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT 10`)
	query := "/v1/co2/weekly"
	validDates := []string{"1974-05-19", "1974-05-26", "1984-01-01", "1984-01-08", "2000-01-02", "2000-01-09", "2018-09-02", "2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), nil, sqlString, query, validDates, handlerConfig)
}

func TestCo2GetYear(t *testing.T) {
	testVal := 2020

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE year in ('%v') ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?year=%v", testVal)
	validDates := []string{"2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates, handlerConfig)
}

func TestCo2GetMonth(t *testing.T) {
	testVal := 1

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE month in ('%v') ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?month=%v", testVal)
	validDates := []string{"1984-01-01", "1984-01-08", "2000-01-02", "2000-01-09"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates, handlerConfig)
}

func TestCo2GetGt(t *testing.T) {
	testVal := 405.68

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average > %v ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?gt=%v", testVal)
	validDates := []string{"2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates, handlerConfig)
}

func TestCo2GetGte(t *testing.T) {
	testVal := 405.68

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average >= %v ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?gte=%v", testVal)
	validDates := []string{"2018-09-02", "2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates, handlerConfig)
}

func TestCo2GetLt(t *testing.T) {
	testVal := 344.19

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average < %v ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?lt=%v", testVal)
	validDates := []string{"1974-05-19", "1974-05-26", "1984-01-08"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates, handlerConfig)
}

func TestCo2GetLte(t *testing.T) {
	testVal := 344.19

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average <= %v ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?lte=%v", testVal)
	validDates := []string{"1974-05-19", "1974-05-26", "1984-01-01", "1984-01-08"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates, handlerConfig)
}

func TestCo2GetLimit(t *testing.T) {
	testVal := 2

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT %v`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?limit=%v", testVal)
	validDates := []string{"1974-05-19", "1974-05-26"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates, handlerConfig)
}

func TestCo2GetOffset(t *testing.T) {
	testVal := 4

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT 10 OFFSET %v`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly?offset=%v", testVal)
	validDates := []string{"2000-01-02", "2000-01-09", "2018-09-02", "2018-10-07", "2020-02-02", "2020-05-24"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates, handlerConfig)
}

func TestCo2GetPage(t *testing.T) {
	page := 2
	limit := 2

	offset := (limit * (page - 1))

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo ORDER BY year,month,day LIMIT %v OFFSET %v`, limit, offset))
	query := fmt.Sprintf("/v1/co2/weekly?limit=%v&page=%v", limit, page)
	validDates := []string{"1984-01-01", "1984-01-08"}

	RunTest(t, t.Name(), offset, sqlString, query, validDates, handlerConfig)
}

func TestCo2GetCombo(t *testing.T) {
	years := []int{1984, 2000}
	month := 1
	gt := 344.19
	gte := 343.89
	lt := 369.03
	lte := 368.89

	// This regex will match the SELECT query with any arbitrary ordering of the WHERE clauses. This is needed because the order that the server concatenates WHERE clauses is semi-random
	sqlString := `SELECT \* FROM public\.co2_weekly_mlo WHERE (average [<>=]+ [\d\.]+( AND )*|year in \(('[\d]+'(,)?[ ]*)*\)( AND )*|month in \(('[\d]+'(,)?[ ]*)*\)( AND )*)* ORDER BY year,month,day LIMIT 10`
	query := fmt.Sprintf("/v1/co2/weekly?year=%v,%v&month=%v&gt=%v&gte=%v&lt=%v&lte=%v", years[0], years[1], month, gt, gte, lt, lte)
	validDates := []string{"1984-01-08", "2000-01-02"}

	RunTest(t, t.Name(), []float32{343.89, 368.89}, sqlString, query, validDates, handlerConfig)
}

func TestCo2GetNull(t *testing.T) {
	testVal := 500.00

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.co2_weekly_mlo WHERE average > %.2f ORDER BY year,month,day LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/co2/weekly/increase?gt=%v", testVal)
	validValues := []string{}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validValues, handlerConfig)
}

func TestCo2Errors(t *testing.T) {
	testVals := []string{
		"/v1/co2/weekly?year=2020a",
		"/v1/co2/weekly?year=20200",
		"/v1/co2/weekly?month=1a",
		"/v1/co2/weekly?month=14",
		"/v1/co2/weekly?gt=400a",
		"/v1/co2/weekly?lt=400a",
		"/v1/co2/weekly?gte=400a",
		"/v1/co2/weekly?lte=400a",
		"/v1/co2/weekly?gt=40000",
		"/v1/co2/weekly?gt=-1",
	}

	sqlString := ``
	validValues := []string{"400"} // The http response code we're expecting

	for _, v := range testVals {
		query := fmt.Sprintf("%v", v)
		RunTest(t, t.Name(), nil, sqlString, query, validValues, handlerConfig)
	}
}
