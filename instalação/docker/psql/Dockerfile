FROM postgres:9.6

ENV PGDATA=/db

COPY atiradorfrequente.sql /
COPY entrypoint.sh /docker-entrypoint-initdb.d/

EXPOSE 5432