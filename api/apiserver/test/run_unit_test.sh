# Copyright 2021 The PlanetPulse Authors.

# Planet Pulse is an API designed to serve climate data pulled from NOAA's
# Global Monitoring Laboratory FTP server. This API is based on the
# OpenAPI v3 specification.

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# A copy of the GNU General Public License can be found here:
# https://www.gnu.org/licenses/

# Contact: planetpulse.api@gmail.com

#!/usr/bin/env bash

while [[ $# -gt 0 ]]; do
  key="$1"

  case $key in
    -v|--verbose)
      VERBOSE=true
      shift # past argument
      shift # past value
      ;;
    *)    # unknown option
      echo "run_unit_tests.sh $1: unknown argument"
      cat << EOF
Usage: ./run_unit_tests.sh [OPTIONS]

Options: 
    -v, --verbose        run tests with verbosity on
EOF
      exit 1
      ;;
  esac
done

if [ "$VERBOSE" = true ] ; then
    go test -v $(go list apiserver/...|grep -ve '.*/test') -verbose
else
    go test -v $(go list apiserver/...|grep -ve '.*/test')
fi