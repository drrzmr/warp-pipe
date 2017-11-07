all: test

linter:
	gometalinter --config=metalinter.json ./...

test:
	go test -v -race -covermode=atomic ./...

dep:
	dep ensure -v

docker-test:
	docker-compose run --rm golang make test

docker-linter:
	docker-compose run --rm golang make linter

docker-dep:
	docker-compose run --rm golang make dep


.PHONY: \
	all \
	linter \
	test \
	dep \
	docker-linter \
	docker-dep \
	docker-test
