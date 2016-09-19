#!/bin/sh

set -e

createdb --username "$POSTGRES_USER" atiradorfrequente
createuser --username "$POSTGRES_USER" atiradorfrequente

psql --username "$POSTGRES_USER" -c "ALTER USER atiradorfrequente WITH PASSWORD '$POSTGRES_PASSWORD';"
psql --username "$POSTGRES_USER" -c "GRANT CONNECT ON DATABASE atiradorfrequente to atiradorfrequente;"
psql --username "$POSTGRES_USER" atiradorfrequente < /tmp/atiradorfrequente.sql
psql --username "$POSTGRES_USER" -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO atiradorfrequente;" atiradorfrequente
psql --username "$POSTGRES_USER" -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO atiradorfrequente;" atiradorfrequente