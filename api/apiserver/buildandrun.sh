#! /usr/bin/env bash

podman build --rm -t apiserver:latest .
podman run --name apiserver --rm --rm -p 8080:8080 -d apiserver:latest
