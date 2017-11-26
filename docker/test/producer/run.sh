#!/bin/bash -eux

docker-compose up -d
docker-compose exec golang bash -c 'go run main.go producer'
docker-compose exec kafka kafka-topics --describe --topic test --zookeeper zookeeper:32181
docker-compose exec kafka kafka-console-consumer --bootstrap-server kafka:29092 --topic test --new-consumer --from-beginning --max-messages 1
docker-compose down
