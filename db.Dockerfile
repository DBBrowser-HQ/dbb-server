FROM postgres:latest

COPY ./migrations/init.sql /docker-entrypoint-initdb.d/