FROM golang:1.15.2-stretch AS builder

WORKDIR /build
COPY . .

USER root
RUN go build -o api-server ./cmd/api-server/main.go

FROM ubuntu:20.04
COPY . .

EXPOSE 5432
EXPOSE 5000

RUN apt-get -y update && apt -y install postgresql-12
USER postgres

RUN  /etc/init.d/postgresql start &&\
    psql --command "CREATE USER test WITH SUPERUSER PASSWORD 'test';" &&\
    createdb -O forum_db forum_db &&\
    psql -f /scripts/postgresql/init_db.sql -d forum_db &&\
    /etc/init.d/postgresql stop


RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/12/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/12/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root
COPY --from=builder  /build/run /usr/bin
CMD /etc/init.d/postgresql start && api-server