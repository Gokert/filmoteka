FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . .
RUN go build ./cmd/main.go

FROM ubuntu:latest

FROM redis:latest

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

WORKDIR /filmoteka
COPY --from=builder /build/main .

COPY . .

EXPOSE 8081
EXPOSE 6379
EXPOSE 5432

ENV APPLICATION_PORT=8081
ENV PSX_PORT=6379
ENV REDIS_PORT=5432
ENV DB_USER=user
ENV DB_NAME=vk_filmoteka



#CMD ["docker-compose", "up"]
CMD service postgresql start && ./main
