#! /usr/bin/env sh

echo "Initializing database"
psql -U postgres -c "CREATE DATABASE food;" template1
psql -U postgres -d food -a -f /data/schema.sql
psql -U postgres -d food -a -f /data/import.sql
