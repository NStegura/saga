.PHONY: up
up:
	docker-compose up --build --no-recreate --detach; sleep 5

.PHONY: down
down:
	docker-compose down --remove-orphans --rmi local