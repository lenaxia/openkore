#!/bin/bash

# Script to run the server connection integration test

# Change to the implementation directory
cd "$(dirname "$0")/.."

# Compile and run the server connection test
go run server_connection.go

# Exit with the same status as the test
exit $?