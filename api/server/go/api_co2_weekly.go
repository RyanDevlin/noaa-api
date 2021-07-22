/*
 * Planet Pulse
 *
 * Planet Pulse is an API designed to serve climate data pulled from NOAA's Global Monitoring Laboratory FTP server. This API is based on the OpenAPI v3 specification.
 *
 * API version: 0.1.0
 * Contact: planetpulse.api@gmail.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"net/http"
)

func Co2WeeklyGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func Co2WeeklyIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
