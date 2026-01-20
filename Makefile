DATASTORE_VERSION ?= latest
TEST_TIMEOUT ?= 240s
DATASTORE_QUORUM_INSERT ?= 1

.PHONY: build test lint up down up-cluster down-cluster cli docker docker-push contributors staticcheck codegen

build:
	@go build ./...

test:
	@go install -race -v
	@DATASTORE_VERSION=$(DATASTORE_VERSION) DATASTORE_QUORUM_INSERT=$(DATASTORE_QUORUM_INSERT) go test -race -timeout $(TEST_TIMEOUT) -count=1 -v ./...

lint:
	golangci-lint run || :

staticcheck:
	staticcheck ./...

up:
	@docker compose up --wait

down:
	@docker compose down

up-cluster:
	@docker compose -f docker-compose.cluster.yml up

down-cluster:
	@docker compose -f docker-compose.cluster.yml down

cli:
	docker run -it --rm --net datastore-go_datastore --link datastore-server:datastore-server hanzoai/datastore-server:$(DATASTORE_VERSION) clickhouse-client --host datastore-server

docker:
	docker build -t hanzoai/datastore-server:$(DATASTORE_VERSION) -f Dockerfile .
	docker build -t hanzoai/datastore-proxy:$(DATASTORE_VERSION) -f Dockerfile.proxy .

docker-push: docker
	docker push hanzoai/datastore-server:$(DATASTORE_VERSION)
	docker push hanzoai/datastore-proxy:$(DATASTORE_VERSION)

contributors:
	@git log --pretty="%an <%ae>%n%cn <%ce>" | sort -u -t '<' -k 2,2 | LC_ALL=C sort | \
		grep -v "users.noreply.github.com\|GitHub <noreply@github.com>" \
		> contributors/list

codegen: contributors
	@go run lib/column/codegen/main.go
