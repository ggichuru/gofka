run-broker:
	docker-compose up -d

test-producer:
	go run src/producer/producer.go

test-consumer:
	go run src/consumer/consumer.go

run-sumulation:
	./simulate.sh