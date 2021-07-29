### TODO:

- balancer
- readme
- e2e

### Improves:

- security input URL (local requests...)
- redis as queue
- limit prices in csv

### Request with grpcurl

grpcurl -plaintext -d '{"url": "http://localhost:3000/generator.csv?count=10000"}' localhost:50051 proto.Price/Fetch

grpcurl -plaintext -d '{"skip": 0, "limit": 100, "order_by": "changes", "order_type": -1}' localhost:50051 proto.Price/List
