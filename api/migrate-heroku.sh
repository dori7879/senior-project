#!/usr/bin/env bash

echo "Running migrations..."
./migrate up > /dev/null 2>&1 &