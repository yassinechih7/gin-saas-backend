FROM postgres:12.13-alpine

# ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres

# COPY POSTGRES_INIT_DB
COPY database.sql /docker-entrypoint-initdb.d/init.sql