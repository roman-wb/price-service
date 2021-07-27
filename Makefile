server:
	go run cli/server/main.go

generate:
	go generate ./...

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/price.proto

test:
	go test -count 1 -race -coverprofile=coverage.out ./...

cover: test
	go tool cover -html=coverage.out
