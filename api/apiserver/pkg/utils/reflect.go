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
	"reflect"
)

// This function is experimental and no guaranteed to work. It was implemented
// as a potential solution for working with interfaces down the road.
func GetFieldVal(obj interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(obj)
	field := reflect.ValueOf(val).FieldByName(fieldName)

	switch val.FieldByName(fieldName).Kind() {
	case reflect.Bool:
		return field.Bool()
	case reflect.Int:
		return field.Int()
	case reflect.Int8:
		return int8(field.Int())
	case reflect.Int16:
		return int16(field.Int())
	case reflect.Int32:
		return int32(field.Int())
	case reflect.Int64:
		return int32(field.Int())
	case reflect.Uint:
		return field.Uint()
	case reflect.Uint8:
		return uint8(field.Uint())
	case reflect.Uint16:
		return uint16(field.Uint())
	case reflect.Uint32:
		return uint32(field.Uint())
	case reflect.Uint64:
		return uint64(field.Uint())
	case reflect.Uintptr:
		return uintptr(field.Uint())
	case reflect.Float32:
		return float32(field.Float())
	case reflect.Float64:
		return field.Float()
	case reflect.Complex64:
		return complex64(field.Complex())
	case reflect.Complex128:
		return complex128(field.Complex())
	case reflect.Array:
		return field.Interface()
	case reflect.Chan:
		return field.Interface()
	case reflect.Func:
		return field.Interface()
	case reflect.Interface:
		return field.Interface()
	case reflect.Map:
		return field.Interface()
	case reflect.Ptr:
		return field.Interface()
	case reflect.Slice:
		return field.Interface()
	case reflect.String:
		return field.String()
	case reflect.Struct:
		return field.Interface()
	case reflect.UnsafePointer:
		return field.Interface()
	}
	return nil
}
