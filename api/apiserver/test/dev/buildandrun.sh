#! /usr/bin/env bash

build_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

go get -d -v ${build_path}/../../...
CGO_ENABLED=0 go build -a -o ${build_path}/planetpulse ${build_path}/../../cmd/planetpulse/

cd ${build_path}
./planetpulse