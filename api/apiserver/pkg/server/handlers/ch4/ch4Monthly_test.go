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

package ch4

import (
	"apiserver/pkg/server/handlers"
	"fmt"
	"regexp"
	"testing"
)

var handlerConfig = &handlers.ApiHandlerConfig{
	SortBy: "average",
}

func TestCh4GetAll(t *testing.T) {
	sqlString := regexp.QuoteMeta(`SELECT * FROM public.ch4_mm_gl ORDER BY year,month LIMIT 10`)
	query := "/v1/ch4/monthly"
	validDates := []string{"1983.542", "1983.625", "1990.042", "1990.125", "2000.042", "2000.125", "2020.792", "2020.875"}

	RunTest(t, t.Name(), nil, sqlString, query, validDates, handlerConfig)
}

func TestCh4GetYear(t *testing.T) {
	testVal := 2020

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.ch4_mm_gl WHERE year in ('%v') ORDER BY year,month LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/ch4/monthly?year=%v", testVal)
	validDates := []string{"2020.792", "2020.875"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates, handlerConfig)
}

func TestCh4GetMonth(t *testing.T) {
	testVal := 1

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.ch4_mm_gl WHERE month in ('%v') ORDER BY year,month LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/ch4/monthly?month=%v", testVal)
	validDates := []string{"1990.042", "2000.042"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates, handlerConfig)
}

func TestCh4GetGt(t *testing.T) {
	testVal := 1890.1

	sqlString := `SELECT \* FROM public\.ch4_mm_gl WHERE average > [\d\.]+ ORDER BY year,month LIMIT 10`
	query := fmt.Sprintf("/v1/ch4/monthly?gt=%v", testVal)
	validDates := []string{"2020.875"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates, handlerConfig)
}

func TestCh4GetGte(t *testing.T) {
	testVal := 1890.1

	sqlString := `SELECT \* FROM public\.ch4_mm_gl WHERE average >= [\d\.]+ ORDER BY year,month LIMIT 10`
	query := fmt.Sprintf("/v1/ch4/monthly?gte=%v", testVal)
	validDates := []string{"2020.792", "2020.875"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates, handlerConfig)
}

func TestCh4GetLt(t *testing.T) {
	testVal := 1627.5

	sqlString := `SELECT \* FROM public\.ch4_mm_gl WHERE average < [\d\.]+ ORDER BY year,month LIMIT 10`
	query := fmt.Sprintf("/v1/ch4/monthly?lt=%v", testVal)
	validDates := []string{"1983.542"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates, handlerConfig)
}

func TestCh4GetLte(t *testing.T) {
	testVal := 1627.5

	sqlString := `SELECT \* FROM public\.ch4_mm_gl WHERE average <= [\d\.]+ ORDER BY year,month LIMIT 10`
	query := fmt.Sprintf("/v1/ch4/monthly?lte=%v", testVal)
	validDates := []string{"1983.542", "1983.625"}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validDates, handlerConfig)
}

func TestCh4GetLimit(t *testing.T) {
	testVal := 2

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.ch4_mm_gl ORDER BY year,month LIMIT %v`, testVal))
	query := fmt.Sprintf("/v1/ch4/monthly?limit=%v", testVal)
	validDates := []string{"1983.542", "1983.625"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates, handlerConfig)
}

func TestCh4GetOffset(t *testing.T) {
	testVal := 4

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.ch4_mm_gl ORDER BY year,month LIMIT 10 OFFSET %v`, testVal))
	query := fmt.Sprintf("/v1/ch4/monthly?offset=%v", testVal)
	validDates := []string{"2000.042", "2000.125", "2020.792", "2020.875"}

	RunTest(t, t.Name(), testVal, sqlString, query, validDates, handlerConfig)
}

func TestCh4GetPage(t *testing.T) {
	page := 2
	limit := 2

	offset := (limit * (page - 1))

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.ch4_mm_gl ORDER BY year,month LIMIT %v OFFSET %v`, limit, offset))
	query := fmt.Sprintf("/v1/ch4/monthly?limit=%v&page=%v", limit, page)
	validDates := []string{"1990.042", "1990.125"}

	RunTest(t, t.Name(), offset, sqlString, query, validDates, handlerConfig)
}

func TestCh4GetCombo(t *testing.T) {
	years := []int{1990, 2000}
	month := []int{1, 2}
	gt := 1712.1
	gte := 1713.5
	lt := 1776.1
	lte := 1776

	// This regex will match the SELECT query with any arbitrary ordering of the WHERE clauses. This is needed because the order that the server concatenates WHERE clauses is semi-random
	sqlString := `SELECT \* FROM public\.ch4_mm_gl WHERE (average [<>=]+ [\d\.]+( AND )*|year in \(('[\d]+'(,)?[ ]*)*\)( AND )*|month in \(('[\d]+'(,)?[ ]*)*\)( AND )*)* ORDER BY year,month LIMIT 10`
	query := fmt.Sprintf("/v1/ch4/monthly?year=%v,%v&month=%v,%v&gt=%v&gte=%v&lt=%v&lte=%v", years[0], years[1], month[0], month[1], gt, gte, lt, lte)
	validDates := []string{"1990.125", "2000.125"}

	RunTest(t, t.Name(), []float32{1776.00, 1713.5}, sqlString, query, validDates, handlerConfig)
}

func TestCh4GetNull(t *testing.T) {
	testVal := 500.00

	sqlString := regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM public.ch4_mm_gl WHERE average < %.2f ORDER BY year,month LIMIT 10`, testVal))
	query := fmt.Sprintf("/v1/ch4/monthly?lt=%v", testVal)
	validValues := []string{}

	RunTest(t, t.Name(), float32(testVal), sqlString, query, validValues, handlerConfig)
}

func TestCh4Errors(t *testing.T) {
	testVals := []string{
		"/v1/ch4/monthly?year=2020a",
		"/v1/ch4/monthly?year=20200",
		"/v1/ch4/monthly?month=1a",
		"/v1/ch4/monthly?month=14",
		"/v1/ch4/monthly?gt=400a",
		"/v1/ch4/monthly?lt=400a",
		"/v1/ch4/monthly?gte=400a",
		"/v1/ch4/monthly?lte=400a",
		"/v1/ch4/monthly?gt=40000",
		"/v1/ch4/monthly?gt=-1",
	}

	sqlString := ``
	validValues := []string{"400"} // The http response code we're expecting

	for _, v := range testVals {
		query := fmt.Sprintf("%v", v)
		RunTest(t, t.Name(), nil, sqlString, query, validValues, handlerConfig)
	}
}
