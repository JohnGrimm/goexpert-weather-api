conf ?= .env
include $(conf)
export $(shell sed 's/=.*//' $(conf))



## ---------- UTILS
.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@rm -f coverage.*



## ----- BUILD & PUSH
.PHONY: build
build: ## Build the container image
	@docker build -t gcr.io/full-cycle-go-expert/goexpert-weather-api:v1 .

.PHONY: push
push: ## Push the container image to image registry
	@docker push gcr.io/full-cycle-go-expert/goexpert-weather-api:v1



## ----- COMPOSE
.PHONY: up
up: ## Put the compose containers up
	@docker-compose up -d --build

.PHONY: down
down: ## Put the compose containers down
	@docker-compose down



## ----- MAIN
.PHONY: serve
serve: ## Run the server
	@cd cmd/app && go run main.go
	@cd -

.PHONY: run
run: ## Make some requests
	@echo -n "422: "; curl -s "http://localhost:8080/cep/1234567"
	@echo -n "404: "; curl -s "http://localhost:8080/cep/12345678"
	@echo -n "200: "; curl -s "http://localhost:8080/cep/13330250"

.PHONY: test
test: ## Run the tests
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html coverage.out -o coverage.html
