#!/usr/bin/env bash

podman run \
  -p 9090:9090/tcp \
  prom/prometheus:latest