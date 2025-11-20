APP_NAME := app
DOCKER_COMPOSE := docker compose
SERVICE ?= app

.PHONY: help build run stop clean logs test up ps

help:
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

build: ## Build the Docker image
	$(DOCKER_COMPOSE) build

up: ## Start all services in background
	$(DOCKER_COMPOSE) up -d

down: ## Stop all services
	$(DOCKER_COMPOSE) down

logs: ## Show service logs
	$(DOCKER_COMPOSE) logs -f $(SERVICE)

clean: ## Remove containers and volumes
	$(DOCKER_COMPOSE) down -v

test: ## Run tests
	go test ./tests/... -count=1

restart: ## Restart services
	$(DOCKER_COMPOSE) restart

ps: ## Show running services
	$(DOCKER_COMPOSE) ps

rebuild : down build up ## Stop all services. Rebuild app and Start all services again 