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

package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

func PlanetDBConnect() {
	conninfo := "user=planet_pulse dbname=planetpulse password=boilerup host=planetpulse.ch0g0ophcqsz.us-east-1.rds.amazonaws.com sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	minReconn := 10 * time.Second
	maxReconn := time.Minute
	listener := pq.NewListener(conninfo, minReconn, maxReconn, reportProblem)
	err = listener.Listen("getwork")
	if err != nil {
		panic(err)
	}
	// TODO: Change to logger
	fmt.Println("entering main loop")
	db.Close()
	/*for {
		// process all available work before waiting for notifications
		getWork(db)
		waitForNotification(listener)
	}*/
}
