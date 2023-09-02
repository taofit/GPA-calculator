#===================#
#== Env Variables ==#
#===================#
DOCKER_COMPOSE_FILE ?= docker-compose.yml
# Version
V?=

#========================#
#== DATABASE MIGRATION ==#
#========================#

start:
	docker compose down
	docker compose up

migrate-up: ## Run migrations UP
migrate-up:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate up

migrate-down: ## Rollback migrations against non test DB
migrate-down:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate down

migrate-create: ## Create a DB migration files e.g `make migrate-create name=migration-name`
migrate-create:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate create -ext sql -seq -dir /migrations $(name)

migrate-force-version: ## forece to a version command: make migrate-force-version V=2. V is the forced version number
migrate-force-version:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate force $(V)

shell-db: ## Enter to database console
shell-db:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec database psql -U admin -d gpa_calculator