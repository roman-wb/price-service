run-service:
	go run cli/service/main.go

run-static-server:
	go run cli/static-server/main.go

generate:
	go generate ./...

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/price.proto

unit-test:
	go test -short -count 1 -race -coverprofile=coverage.out ./...

test:
	go test -count 1 -race -coverprofile=coverage.out ./...

cover: test
	go tool cover -html=coverage.out§

docker-up:
	docker-compose -f docker-compose.test.yml up --force-recreate --remove-orphans

create-migration:
	docker run --rm -it -v `pwd`/migrations:/migrations --network host migrate/migrate create -ext json -dir=/migrations $(name)