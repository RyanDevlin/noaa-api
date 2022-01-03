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

	// Pretty controls whether or not to pretty-print json responses.
	Pretty bool
}

// Query queries the database according to the supplied DBQuery.
// It loads a supplied dataObject with the requested data.
func (database *Database) Query(query DBQuery, dataObject models.DataObject) error {
	if err := database.ProbeConnection(); err != nil {
		return err
	}

	rows, err := database.DB.Query(query.ToString())
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		err := dataObject.Load(rows, query.Simple)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
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

// NewQuery returns an initialized DBQuery to be used by a handler to
// setup a new database request.
func NewQuery(table string, cols []string, orderBy string) DBQuery {
	return DBQuery{
		Table:   table,
		Cols:    cols,
		OrderBy: orderBy,
		Limit:   10,
		Offset:  0,
		Page:    0,
		Simple:  false,
		Pretty:  true,
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
