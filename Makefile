all: test

linter:
	gometalinter --config=metalinter.json ./...

test:
	go test -v -race -covermode=atomic ./...

docker-test:
	docker-compose run --rm golang make test

docker-linter:
	docker-compose run --rm golang make linter


.PHONY: \
	all \
	linter \
	test \
	docker-linter \
	docker-test
