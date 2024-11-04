test:
	go test -v ./...

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

build-and-run-docker:
	docker compose up --detach --force-recreate --build server --build client
	docker logs wow-server
	docker logs wow-client

run-docker:
	docker compose up --detach --force-recreate
	docker logs wow-server
	docker logs wow-client

lint:
	golangci-lint -v run ./...
