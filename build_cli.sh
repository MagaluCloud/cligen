#!/bin/bash

cd tmp-cli

# Build the CLI
go build -o cli -ldflags "-X 'main.RawVersion=v0.0.0-$(date +%y%m%d_%H%M%S)'" -v

if [ $? -ne 0 ]; then
    echo "Failed to build the CLI"
    exit 1
fi

# Run the CLI
./cli --version