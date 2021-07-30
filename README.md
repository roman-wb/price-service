[![Build Status](https://www.travis-ci.com/roman-wb/price-service.svg?branch=master)](https://www.travis-ci.com/roman-wb/price-service)
![Go Report](https://goreportcard.com/badge/github.com/roman-wb/price-service)
![Repository Top Language](https://img.shields.io/github/languages/top/roman-wb/price-service)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/roman-wb/price-service)
![Github Repository Size](https://img.shields.io/github/repo-size/roman-wb/price-service)
![Lines of code](https://img.shields.io/tokei/lines/github/roman-wb/price-service)
![License](https://img.shields.io/badge/license-MIT-green)
![GitHub last commit](https://img.shields.io/github/last-commit/roman-wb/price-service)

# GRPC Price service

## Features

- gRPC Service with MongoDB storage
- Method Fetch(url) for export CVS file with list of product from URL
  - Format file PRODUCT_NAME;PRICE
  - Last price should be saved in DB with request date
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

### Development

```bash
# Clone repo
git clone github.com/roman-wb/price-service
cd price-service
# Run Docker with Mongo
make docker-up
# Run gRPC (foreground)
make run-service
# Run server with generator CSV (foreground)
make run-static-server
```

### Local server on localhost:8080

```bash
make server
```

### Docker

```bash
make docker-build
make docker-run
```

### Mapping to the internet with ngrok

Note: Require installed [ngrok](https://ngrok.com)

```bash
make ngrok
```

## Inspired

https://github.com/gorilla/websocket/tree/master/examples/chat

## License

MIT

### TODO

- Run dev mode
- Run manual test mode
- Run prod
- Run test
- TravisCI + Coverage

- Balancer
- Readme
- e2e tests

### Options

- Security input URL (hack with local request)?
- Redis or another queue messages?
- Limit count prices in csv?

### Request with grpcurl

grpcurl -plaintext -d '{"url": "http://localhost:3000/generator.csv?count=10000"}' localhost:50051 proto.Price/Fetch

grpcurl -plaintext -d '{"skip": 0, "limit": 100, "order_by": "changes", "order_type": -1}' localhost:50051 proto.Price/List
