#!/bin/bash

echo "build geocoding.linux ..."
GOOS=linux GOARCH=amd64 go build -o geocoding.linux geo/geocoding.go
echo "geocoding.linux created"
