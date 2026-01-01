.PHONY: all create-volume run

all: create-volume run

create-volume:
	docker network create shared_network
	docker volume inspect db-data >/dev/null 2>&1 || docker volume create db-data

run:
	docker-compose up
