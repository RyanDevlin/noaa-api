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

package database

import (
	"apiserver/pkg/database/models"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	DB     *sql.DB
	Config *DBConfig
}

type DBConfig struct {
	// The database endpoint
	DBHost string

	// The database username
	DBUser string

	// The database password
	DBPass string

	// (OPTIONAL) The port the database listens on
	DBPort int

	// (OPTIONAL) The connection timeout in seconds used when connecting to the database
	DBConnTimeout int
}

type DBQuery struct {
	// The name of the table to query
	Table string

	// The columns to select from
	Cols []string

	// A list of Boolean SQL expressions to be used as 'WHERE' clauses
	Where []string

	// An expression passed directly to the ORDER BY keyword. Usually should just be one or more cols.
	OrderBy string

	// Simple provides a way to tell the Query function that the data returned will be simplified.
	Simple bool
}

func (database *Database) Query(query DBQuery) (models.Co2Table, error) {
	if err := database.ProbeConnection(); err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(query.ToString())
	if err != nil {
		return nil, err
	}

	co2table := models.Co2Table{}
	defer rows.Close()
	for rows.Next() {
		if query.Simple {
			var co2entry models.Co2EntrySimple
			var year, month, day int

			if err := rows.Scan(&year, &month, &day, &co2entry.Average, &co2entry.IncSincePreIndustrial); err != nil {
				return nil, err
			}

			// Use the unique date of measurement as the key to the co2table
			key := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
			co2table[key] = co2entry
		} else {
			var co2entry models.Co2Entry
			if err := rows.Scan(&co2entry.Year, &co2entry.Month, &co2entry.Day, &co2entry.DateDecimal, &co2entry.Average, &co2entry.NumDays, &co2entry.OneYearAgo, &co2entry.TenYearsAgo, &co2entry.IncSincePreIndustrial, &co2entry.Timestamp); err != nil {
				return nil, err
			}

			// Use the unique date of measurement as the key to the co2table
			year, month, day := co2entry.Timestamp.Date()
			key := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
			co2table[key] = co2entry
		}

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return co2table, nil
}

func (database *Database) Connect() error {
	conninfo := fmt.Sprintf("postgres://%s:%s@%s/postgres?connect_timeout=%d", url.PathEscape(database.Config.DBUser), url.PathEscape(database.Config.DBPass), database.Config.DBHost, database.Config.DBConnTimeout)

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		return err
	}

	// Validate conninfo args with ping
	if err = db.Ping(); err != nil {
		db.Close()
		return err
	}

	database.DB = db
	return nil
}

// ProbeConnection provides a safe mechanism for checking the database connection.
// This function should be called before a query is made. It first detects if a database connection
// has not been initialized. If this is the case a new connection attempt will be made.
// An error is returned
func (database *Database) ProbeConnection() error {
	// If database failed to initialize, apiserver.Db will be nil
	if database.DB == nil {
		log.Error("Database connection has not been established.")

		log.Info("Retrying database connection....")
		status := database.Connect()
		if status != nil {
			return status
		}
		log.Info("Database connection successfully established.")
	}

	if err := database.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (query DBQuery) ToString() string {
	sqlString := "SELECT "

	// Append columns to select from
	for i, col := range query.Cols {
		if i == len(query.Cols)-1 {
			sqlString += col + " "
			break
		}
		sqlString += col + ", "
	}

	sqlString += "FROM " + query.Table + " "

	if len(query.Where) >= 1 {
		sqlString += "WHERE "
	}

	for i, expr := range query.Where {
		if i == len(query.Where)-1 {
			sqlString += expr + " "
			break
		}
		sqlString += expr + " AND "
	}

	if query.OrderBy != "" {
		sqlString += "ORDER BY " + query.OrderBy
	}

	return sqlString
}
