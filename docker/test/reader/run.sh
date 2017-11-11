#!/bin/bash -eux

docker-compose up -d
docker-compose exec golang bash -c 'psql < /scripts/sql/reader/start.sql'
docker-compose exec golang bash -c 'psql < /scripts/sql/reader/events.sql'
docker-compose exec golang bash -c 'pg_recvlogical -d reader -h postgres-server -S td -v --start -f - | go run main.go reader' || \
docker-compose exec golang bash -c 'psql < /scripts/sql/reader/stop.sql'
docker-compose down
