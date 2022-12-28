FROM postgres:12.13-alpine

# ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
# # ENV POSTGRES_DB golang_gin_db


COPY database.sql /docker-entrypoint-initdb.d/init.sql