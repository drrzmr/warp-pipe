all: test

test:
	go test -v -race -covermode=atomic ./...

docker-test:
	docker-compose run --rm golang make test

.PHONY: \
	all \
	test \
	docker-test
