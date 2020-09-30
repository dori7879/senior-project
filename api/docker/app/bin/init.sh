#!/usr/bin/env bash
echo "Running migrations..."
/api/migrate up > /dev/null 2>&1 &

echo 'Deleting mysql-client...'
apk del mysql-client

echo 'Starting application...'
/api/app