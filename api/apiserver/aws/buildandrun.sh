#! /usr/bin/env bash

podman build --rm -t apiserver:latest ../.
podman run --name apiserver --rm --rmi -v ${PWD}/config/config.yaml:/opt/apiserver/config.yaml:Z --env-file ./config/env.secret -p 8080:8080 -d ghcr.io/ryandevlin/planetpulse/apiserver:latest
