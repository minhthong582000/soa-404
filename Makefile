BINARY=service

proto:
	buf dep update  
	buf generate

test: 
	go test -v -cover -covermode=atomic ./...

service:
	go build -o bin/${BINARY} cmd/random-service/main.go

client:
	go build -o bin/client cmd/client/main.go

unittest:
	go test -short  ./...

clean:
	if [ -f bin/${BINARY} ] ; then rm -rf bin/${BINARY} ; fi

docker:
	docker build -t ${BINARY} -f deploy/docker/Dockerfile .
	docker build -t client -f deploy/docker/Dockerfile.client .

run:
	docker compose -f deploy/docker/docker-compose.yaml up -d --build

stop:
	docker compose -f deploy/docker/docker-compose.yaml down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: clean install unittest build docker run stop lint-prepare lint proto
