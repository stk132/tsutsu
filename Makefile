test:
	docker run -d --rm --name=test_fireworq -p 9090:8080 fireworq/fireworq --driver=in-memory --queue-default=default
	TEST_FIREWORQ_PORT=9090 go test
	docker stop test_fireworq