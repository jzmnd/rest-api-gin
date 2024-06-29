IMAGE := rest-api-gin:latest
CONTAINER := rest-api-gin-service

build:
	@docker build -t $(IMAGE) .

start-dev:
	@go run .

start:
	@docker run -d --rm --name $(CONTAINER) -p 8080:8080 $(IMAGE)

stop:
	@docker stop $(CONTAINER)

.PHONY: build start-dev start stop
