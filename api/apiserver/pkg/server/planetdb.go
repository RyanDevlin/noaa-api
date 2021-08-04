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

package server

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"

	"apiserver/pkg/v1/co2Weekly"

	_ "github.com/lib/pq"
)

func (apiserver *ApiServer) planetDBConnect() error {
	if !apiserver.configured {
		return fmt.Errorf("cannot connect to database, server is not configured")
	}

	conninfo := fmt.Sprintf("postgres://%s:%s@%s/postgres", url.PathEscape(apiserver.Config.DBUser), url.PathEscape(apiserver.Config.DBPass), apiserver.Config.DBHost)
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	// Validate conninfo args with ping
	if err = db.Ping(); err != nil {
		db.Close()
		return err
	}

	apiserver.Db = db
	return nil
}

func (apiserver *ApiServer) PlanetDBGetCo2Table() (co2Weekly.Co2Table, error) {
	sqlStatement := "SELECT * FROM public.co2_weekly_mlo"

	rows, err := apiserver.Db.Query(sqlStatement)
	if err != nil {
		return co2Weekly.Co2Table{}, err
	}

	co2table := co2Weekly.Co2Table{}
	defer rows.Close()
	for rows.Next() {
		var co2entry co2Weekly.Co2Entry
		if err := rows.Scan(&co2entry.Year, &co2entry.Month, &co2entry.Day, &co2entry.DateDecimal, &co2entry.Average, &co2entry.NumDays, &co2entry.OneYearAgo, &co2entry.TenYearsAgo, &co2entry.IncSincePreIndustrial, &co2entry.Timestamp); err != nil {
			return co2Weekly.Co2Table{}, err
		}

		// Use the unique date of measurement as the key to the co2table
		year, month, day := co2entry.Timestamp.Date()
		key := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
		co2table[key] = co2entry
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return co2table, nil
}
