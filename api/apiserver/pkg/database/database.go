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

	// The blank import here is used to import the pq PostgreSQL drivers
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// Database represents a Postgres database and configuration parameters required to connect to it.
type Database struct {
	DB     *sql.DB
	Config *DBConfig
}

// DBConfig represents the configuration parameters required to establish a connection to the database.
// These parameters are loaded from environment variables and the config.yaml file (see config.go).
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

// DBQuery represents an SQL query
type DBQuery struct {
	// The name of the table to query
	Table string

	// The columns to select from
	Cols []string

	// A list of Boolean SQL expressions to be used as 'WHERE' clauses
	Where []string

	// An expression passed directly to the ORDER BY keyword. Usually should just be one or more cols.
	OrderBy string

	// Offset shifts the rows returned using the OFFSET clause of an SQL query
	Offset int

	// Limit limits the number rows returned
	Limit int

	// Page shifts the offset value to provide the next page of data
	Page int

	// Simple provides a way to tell the Query function that the data returned will be simplified.
	Simple bool
}

// Query querys the database according to the supplied DBQuery.
// It returns a Co2Table of the requested data.
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
		if !query.Simple {
			var co2entry models.Co2Entry
			if err := rows.Scan(&co2entry.Year, &co2entry.Month, &co2entry.Day, &co2entry.DateDecimal, &co2entry.Average, &co2entry.NumDays, &co2entry.OneYearAgo, &co2entry.TenYearsAgo, &co2entry.IncSincePreIndustrial, &co2entry.Timestamp); err != nil {
				return nil, err
			}

			// Use the unique date of measurement as the key to the co2table
			year, month, day := co2entry.Timestamp.Date()
			key := strconv.Itoa(year) + "-" + formatInt(int(month)) + "-" + formatInt(day)
			co2table[key] = co2entry
		} else {
			var co2entry models.Co2EntrySimple
			var year, month, day int

			if err := rows.Scan(&year, &month, &day, &co2entry.Average, &co2entry.IncSincePreIndustrial); err != nil {
				return nil, err
			}

			// Use the unique date of measurement as the key to the co2table
			key := strconv.Itoa(year) + "-" + formatInt(int(month)) + "-" + formatInt(day)
			co2table[key] = co2entry
		}

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return co2table, nil
}

// Connect establishes a database connection based on the DBConfig values.
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
// An error is returned when a connection cannot be established.
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

func NewQuery(table string, cols []string, orderBy string) DBQuery {
	return DBQuery{
		Table:   table,
		Cols:    cols,
		OrderBy: orderBy,
		Limit:   -1,
		Offset:  0,
		Page:    0,
		Simple:  false,
	}
}

// ToString marshalls a DBQuery object into a string query that can be sent to an SQL database.
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
		sqlString += "ORDER BY " + query.OrderBy + " "
	}

	offset := query.Offset
	if query.Limit >= 0 {
		offset += (query.Limit * query.Page)
		sqlString += "LIMIT " + strconv.Itoa(query.Limit) + " "
	}

	if offset > 0 {
		sqlString += "OFFSET " + strconv.Itoa(offset)
	}

	return sqlString
}

// formatInt pads integers lower than 10 with a leading '0'.
// This was implemented mainly to combat problems with the way the JSON
// Marshall() function sorts dictionary keys. If ints are not padded with 0
// when they are less than 10, the values will show up out of order in a
// server response.
func formatInt(val int) string {
	if val < 10 {
		return fmt.Sprintf("%02d", val)
	} else {
		return fmt.Sprintf("%d", val)
	}
}
