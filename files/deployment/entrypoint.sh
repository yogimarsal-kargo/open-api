#!/bin/sh
# Entrypoint script for application
set -e

source ./files/etc/dbmate_config/dbmate.sh
dbmate --migrations-dir $DBMATE_DEFAULT_MIGRATIONS_DIR --migrations-table $DBMATE_DEFAULT_MIGRATIONS_TABLE migrate
if [ "$AUTO_MIGRATE_POSTDEPLOYMENT_MIGRATIONS" = "true" ]; then
    dbmate --migrations-dir $DBMATE_POSTDEPLOYMENT_MIGRATIONS_DIR --migrations-table $DBMATE_POSTDEPLOYMENT_MIGRATIONS_TABLE migrate
fi

exec "$@"
