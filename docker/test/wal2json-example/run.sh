#!/bin/bash -eux

docker-compose up -d
docker-compose exec postgres-client bash -c 'psql < /scripts/sql/wal2json-example/create.sql'
docker-compose exec postgres-client bash -c 'psql < /scripts/sql/wal2json-example/events.sql'
docker-compose exec postgres-client bash -c 'psql < /scripts/sql/wal2json-example/drop.sql'
docker-compose down
