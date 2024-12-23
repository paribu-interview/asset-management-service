#!/usr/bin/env bash
set -e

host="$1"
port="$2"
shift 2

>&2 echo "Waiting for $host:$port to become available..."
sleep 3

>&2 echo "$host:$port is available, proceeding..."
exec "$@"
