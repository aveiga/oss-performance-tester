# oss-performance-tester

## Pre-requisites

- Docker

## How to

### Spin up test data services

- `docker compose up -d`

### Run RabbitMQ stress-testing

- `go run main.go rabbitmq --help`
- Load test RabbitMQ with 10 publishers, 5 consumers, a 1MiB message payload and on a queue named "throughput-test": `go run main.go rabbitmq -c 'amqp://localhost' -x 10 -y 5 -s 1000000 -q 'throughput-test'`

### Access the local RabbitMQ UI

- `http://localhost:15672`
