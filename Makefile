test_files = `go list ./... | grep -v /mocks | grep -v /proto | grep -v /cli/static-server`

docker-dev-up:
	docker-compose -f deployments/docker-compose.dev.yml up --force-recreate --remove-orphans

run-dev-service:
	go run cli/service/main.go

run-dev-static-server:
	cd cli/static-server && go run main.go

docker-local-up:
	docker-compose -f deployments/docker-compose.local.yml up --build --force-recreate --remove-orphans

docker-prod-up:
	docker-compose -p price-service-prod -f deployments/docker-compose.prod.yml up --scale service=2 --build --force-recreate --remove-orphans

create-migration:
	docker run --rm -it -v `pwd`/migrations:/migrations --network host migrate/migrate create -ext json -dir=/migrations $(name)

generate:
	go generate ./...

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/proto/price.proto

lint:
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint golangci-lint run -v

unit-test:
	go test -short -count 1 -race -coverprofile=coverage.out $(test_files)

test:
	go test -count 1 -race -coverprofile=coverage.out $(test_files)

cover-total: test
	go tool cover -func=coverage.out

cover: test
	go tool cover -html=coverage.out