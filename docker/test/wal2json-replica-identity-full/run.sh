#!/bin/bash -eux

docker-compose up -d
docker-compose exec postgres-client bash -c 'psql < /scripts/sql/wal2json-replica-identity-full/create.sql'
docker-compose exec postgres-client bash -c 'psql < /scripts/sql/wal2json-replica-identity-full/replication-init.sql'
docker-compose exec postgres-client bash -c 'psql < /scripts/sql/wal2json-replica-identity-full/events.sql'
docker-compose exec postgres-client bash -c 'psql < /scripts/sql/wal2json-replica-identity-full/replication-show.sql'
docker-compose exec postgres-client bash -c 'psql < /scripts/sql/wal2json-replica-identity-full/replication-stop.sql'
docker-compose exec postgres-client bash -c 'psql < /scripts/sql/wal2json-replica-identity-full/drop.sql'
docker-compose down
