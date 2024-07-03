#!/usr/bin/env bash

set -e

DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

export MIGRATION_SOURCE_URL="file://${DIR}/..//db/migrations"
export DB_PATH="${DIR}/../db/data/navigator.db"
export DEV="false"

(cd ${DIR}/frontend && npm run build)

go run .
