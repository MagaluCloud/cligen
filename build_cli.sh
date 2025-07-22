#!/bin/bash

cd tmp-cli

# Build the CLI
go build -o cli -ldflags "-X 'main.RawVersion=v0.0.0-$(date +%y%m%d_%H%M%S)'" -v

# Run the CLI
./cli --version