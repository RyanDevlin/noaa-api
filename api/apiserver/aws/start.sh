#! /usr/bin/env bash

path=$(pwd)
podman run --name apiserver --rm --rmi -v ${path}/config:/ -p 8080:8080 -it ghcr.io/ryandevlin/planetpulse/apiserver:latest
