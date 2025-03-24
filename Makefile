.PHONY: all create-volume run-backend run-reminder clean

all: create-volume run-backend run-reminder

create-volume:
	docker volume inspect db-data >/dev/null 2>&1 || docker volume create db-data

run-backend:
	echo "Starting backend service..."
	cd backend && make run

run-reminder:
	echo "Starting reminder service..."
	cd reminder && docker-compose up -d

clean:
	echo "Cleaning up..."
	cd backend && make stop
	cd reminder && docker-compose down
	docker volume rm db-data
