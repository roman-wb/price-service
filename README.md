[![Build Status](https://www.travis-ci.com/roman-wb/price-service.svg?branch=master)](https://www.travis-ci.com/roman-wb/price-service)
![Go Report](https://goreportcard.com/badge/github.com/roman-wb/price-service)
![Repository Top Language](https://img.shields.io/github/languages/top/roman-wb/price-service)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/roman-wb/price-service)
![Github Repository Size](https://img.shields.io/github/repo-size/roman-wb/price-service)
![Lines of code](https://img.shields.io/tokei/lines/github/roman-wb/prices-service)
![License](https://img.shields.io/badge/license-MIT-green)
![GitHub last commit](https://img.shields.io/github/last-commit/roman-wb/price-service)

# gRPC Price service

## Features

- gRPC Service with MongoDB storage
- Method Fetch(url) - request CVS file from URL with list of products
  - Format file PRODUCT_NAME;PRICE
  - Last price should be saved in storage with request date
  - Save count changes price for every product
- Method List(<paging_params>,<sorting_params>) get list products
  - Fields: name, price, changes, updated_at
  - All variant orders (example infinty scroll)
- Server run with 2+ instances (every in Docker container) + wall with balancer
- Future run in test environment

## Underhood

- grpc / protobuf
- mongo-driver
- golang-migrate
- zap logger

## Get Started

### Development (environment: dev)

```bash
# Clone repo
git clone github.com/roman-wb/price-service
cd price-service
# Run docker with MongoDB on localhost:27017 (foreground)
make docker-dev-up
# Run gRPC service on localhost:50051 (foreground)
make run-service
# [optional] Run static server with generator CSV on localhost:3000 (foreground)
make run-static-server
```

### Local playground (environment: prod)

Setup:

```bash
# Clone repo
git clone github.com/roman-wb/price-service
cd price-service
# Run docker-compose with dependencies
# mongo listen on localhost:27017
# service listen on localhost:50051
# static-server listen on localhost:3000
make docker-local-up
```

Play use [grpcurl](https://github.com/fullstorydev/grpcurl):

```bash

# Request file
grpcurl -plaintext -d '{"url": "http://loalhost:3000/generator.csv?count=100"}' localhost:50051 proto.Price/Fetch
# Get List products
grpcurl -plaintext -d '{"skip": 0, "limit": 1, "order_by": "price", "order_type": -1}' localhost:50051 proto.Price/List
```

### Production (environment: prod)

Service scaled to 2 instances and available via nginx

```bash
# Clone repo
git clone github.com/roman-wb/price-service
cd price-service
# Run docker-compose with dependencies
# mongo unvailable
# service listen on localhost:50051
make docker-prod-up
```

## Test server for generate CSV `static-server`

Generate and return 100 prices (`count=100`) with plain header `plain=true` (see in browser)
`http://localhost:3000/generator.csv?count=100&plain=true`

Generate and return 1000 prices as file generator.csv
`http://localhost:3000/generator.csv?count=1000`

## Makefile commands

- `make docker-dev-up` - Run docker with MongoDB
- `make run-dev-service` - Run service [mode dev]
- `make run-dev-static-server` - Run CSV generator [mode dev]
- `make docker-local-up` - run docker-compose with all configured system (mongo, service, static-server) [mode prod]
- `make docker-prod-up` - run docker-compose with all configured system (nginx, mongo, service) [mode prod]
- `make create-migration` - Create migration in dir /migrations
- `make generate` - Go generate (mocks, etc)
- `make protoc` - Generate proto files
- `make lint` - Run `golangci-lint`
- `make unit-test` - Run unit tests
- `make test` - Run unit + integrations tests (require `make docker-dev-up`)
- `make cover-total` - Run code coverage
- `make cover` - Run `make test` and open browser with code coverage

### TODO

- Readme
- e2e tests

### Options

- Security input URL (hack with local request)?
- Redis or another message queue?
- Limit count prices in csv?

## License

MIT
