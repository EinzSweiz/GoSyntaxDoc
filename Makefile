export COMPOSE_DOCKER_CLI_BUILD=1
export DOCKER_BUILDKIT=1

APP_CONTAINER = fiber-app
EXEC = docker exec -it
DC = docker-compose

# Start services in detached mode
up:
	${DC} -f ${DC}.yaml up -d

# Stop and remove services
down:
	${DC} -f ${DC}.yaml down

# View logs from services
logs:
	${DC} -f ${DC}.yaml logs -f

# Rebuild and start services
build:
	${DC} -f ${DC}.yaml up --build

# Stop services without removing containers
stop:
	${DC} -f ${DC}.yaml stop

# Remove containers, networks, and volumes
clean:
	${DC} -f ${DC}.yaml down --volumes

#Go inside backend container
exec:
	${EXEC} ${APP_CONTAINER} bash

#Logs of backend container
logs-docker:
	docker logs ${APP_CONTAINER}