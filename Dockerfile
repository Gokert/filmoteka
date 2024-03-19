FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . .
RUN go build ./cmd/main.go

FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get -y install postgresql postgresql-contrib ca-certificates
USER postgres
COPY /scripts /opt/scripts
RUN service postgresql start && \
        psql -c "CREATE USER admin WITH superuser login password 'admin';" && \
        psql -c "ALTER ROLE admin WITH PASSWORD 'admin';" && \
        createdb -O admin vk_filmoteka && \
        psql -f ./opt/scripts/sql/init_db.sql -d vk_filmoteka
VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /rest
COPY --from=builder /build/main .

COPY . .


CMD service postgresql start && ./main



