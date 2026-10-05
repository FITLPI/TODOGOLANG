include .env
export

ifeq ($(OS),Windows_NT)
    PROJECT_ROOT := $(subst \,/,$(CURDIR))
else
    PROJECT_ROOT := $(shell pwd)
endif
export PROJECT_ROOT

env-up:
	@docker compose up -d todoapp-postgress

env-down:
	@docker compose down

env-cleanup:
	@read -p "Are you sure you want to clean up the environment? (y/n): " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgress && \
		rm -rf out/pgdata && \
		echo "Environment cleaned up."; \
	else \
		echo "Environment not cleaned up."; \
	fi

migrate-create:
#   Версия для linux и macOS
# 	@if [ -z "$(seq)" ]; then \
# 		echo "Error: Please provide a migration sequence using the seq variable."; \
# 		exit 1; \
# 	fi; \
# 	docker compose run --rm todoapp-postgress-migrate \
# 		create \
# 	 	-ext sql \
# 		-dir /migrations \
# 		-seq "$(seq)"
	@$(if $(strip $(seq)),,@$(error Please provide a migration sequence using the seq variable.))
	MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgress-migrate create -ext sql -dir /migrations -seq "$(seq)"


migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@$(if $(strip $(action)),,@$(error Please provide a migration sequence using the action variable.))
	@MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgress-migrate \
		-path /migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgress:5432/$(POSTGRES_DB)?sslmode=disable" \
		"$(action)"

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder
