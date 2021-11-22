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
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// RequestID represents a request's unique ID value. This is used with a golang
// context to trace the request through log messages. It should be used as the key
// when calling context.WithValue() and should always be set to "RequestID". This
// allows the server logic to look up "RequestID" in the request context and obtain
// the ID value.
type RequestID string

// RequestIdDefaultKey is the default value used as a "tag" to extract the RequestIdDefaultKey
// from a request's context.
const RequestIdDefaultKey = RequestID("RequestID")

// SetReqId initializes a new reuest with a unique ID value in a safe manner.
// Setting the requestId key/value pair without this function will usually result
// in server errors.
func SetReqId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		id, err := newRequestId()
		if err != nil {
			log.Errorf("Failed to generate new Request ID: %v", err)
		}

		log.Tracef("Initializing new request with ID: %v\n", id)

		ctx = context.WithValue(ctx, RequestIdDefaultKey, id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func newRequestId() (string, error) {
	binaryId, err := uuid.New().MarshalBinary()
	if err != nil {
		return "", err
	}

	short, err := shinkBinary(binaryId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", short), nil
}

// shinkBinary is used to shrink a UUID by half. This function XOR's the first half of
// the UUID with the second. This process shrinks the UUID and preserves it's
// distribution characteristics.
func shinkBinary(uuid []byte) ([]byte, error) {
	length := len(uuid)

	// Split the array into two halves
	left, right := uuid[:length/2], uuid[length/2:]
	if len(left) != len(right) || length == 0 {
		return nil, fmt.Errorf("something went wrong, UUID byte slice should be an even length and nonzero")
	}

	for i := 0; i < len(left); i++ {
		left[i] = left[i] ^ right[i] // This copies the new values into the left slice to save some mem
	}
	return left, nil
}

// GetReqId extracts a request's unique ID value and returns it as a string.
// If extraction fails, an error is returned.
func GetReqId(r *http.Request) (string, error) {
	ctx := r.Context()
	id, ok := ctx.Value(RequestIdDefaultKey).(string)
	if !ok {
		return "", fmt.Errorf("request ID key was not set with SetReqId function")
	}
	return id, nil
}
