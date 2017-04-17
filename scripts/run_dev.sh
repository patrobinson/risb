#!/bin/bash -e

go get
go build
curl -X POST "http://${INFLUX_HOSTNAME}:${INFLUX_PORT}/query" \
  --data-urlencode 'q=CREATE DATABASE risb'
./running_is_beautiful