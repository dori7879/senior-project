#!/usr/bin/env bash

host="$1"
shift
cmd="$@"

until mysql -h "$host" -u ${DB_USER} -p${DB_PASS} ${DB_NAME} -e 'select 1'; do
    >&2 echo "MariaDB is unavailable - sleeping"
    sleep 1
done

>&2 echo "MariaDB is up - executing command"
exec $cmd