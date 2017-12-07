all: test

linter:
	gometalinter --config=metalinter.json ./...

test:
	go test -v -race -short -covermode=atomic ./...

test-full:
	go test -v -race -covermode=atomic ./...

dep:
	dep ensure -v

format:
	go fmt ./...
	goimports -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

docker-test:
	docker-compose run --rm golang make test

docker-linter:
	docker-compose run --rm golang make linter

docker-dep:
	docker-compose run --rm golang make dep

docker-format:
	docker-compose run --rm golang make format

.PHONY: \
	all \
	linter \
	test \
	dep \
	format \
	docker-linter \
	docker-dep \
	docker-format \
	docker-test
