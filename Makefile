.PHONY: start stop dev build help
.DEFAULT_GOAL := help

start:
	docker compose up -d

stop:
	docker compose down

dev:
	docker compose up

build:
	docker compose up --build
	
help:
	@echo "  \033[36mstart\033[0m  - Start services in background"
	@echo "  \033[36mstop\033[0m   - Stop services"
	@echo "  \033[36mdev\033[0m    - Start with logs in foreground"
	@echo "  \033[36mbuild\033[0m  - Rebuild and start"