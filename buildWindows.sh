#!/bin/bash

echo "build geocoding.exe ..."
GOOS=windows GOARCH=386 go build -o geocoding.exe geo/geocoding.go
echo "geocoding.exe created"
