#!/bin/bash -eux

docker-compose up -d
docker-compose exec postgres-server bash -c 'psql < /scripts/sql/dump/create.sql'
