#!/bin/sh

set -e

createdb --username "$POSTGRES_USER" atiradorfrequente
psql --username "$POSTGRES_USER" atiradorfrequente < ../atiradorfrequente.sql