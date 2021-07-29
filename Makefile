test_files = `go list ./... | grep -v /mocks | grep -v /proto | grep -v /cli/static-server`

run-service:
	go run cli/service/main.go

run-static-server:
	cd cli/static-server && go run main.go

generate:
	go generate ./...

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/proto/price.proto

unit-test:
	go test -short -count 1 -race -coverprofile=coverage.out $(test_files)

test:
	go test -count 1 -race -coverprofile=coverage.out $(test_files)

cover-total: test
	go tool cover -func=coverage.out

cover: test
	go tool cover -html=coverage.out

docker-up:
	docker-compose -f docker-compose.test.yml up --force-recreate --remove-orphans

create-migration:
	docker run --rm -it -v `pwd`/migrations:/migrations --network host migrate/migrate create -ext json -dir=/migrations $(name)