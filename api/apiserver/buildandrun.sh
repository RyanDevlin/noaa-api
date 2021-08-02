#! /usr/bin/env bash

podman build -t apiserver:latest .
podman run --name apiserver --rm -p 8080:8080 -d apiserver:latest
