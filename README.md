# oss-performance-tester

## Pre-requisites

- Docker

## How to

### Spin up test data services

- `docker compose up -d`

### Run RabbitMQ stress-testing

- `go run main.go rabbitmq --help`
- Load test RabbitMQ with 10 publishers, 5 consumers, a 1MiB message payload and on a queue named "throughput-test": `go run main.go rabbitmq -c 'amqp://localhost' -x 10 -y 5 -s 1000000 -q 'throughput-test'`
- Load test with 10 publishers, 5 consumers and Exchange and Queue details declared in an input file: `go run main.go rabbitmq -c 'amqp://localhost' -x 10 -y 5 -f ./examples/quorum.json`

### Access the local RabbitMQ UI

- `http://localhost:15672`

## References

- [Cobra User Guide](https://github.com/spf13/cobra/blob/main/user_guide.md)
- [Cobra Generator](https://github.com/spf13/cobra-cli/blob/main/README.md)
