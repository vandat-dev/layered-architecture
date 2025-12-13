# name app
APP_NAME = server
BINARY_NAME = app

# Build the application using vendor
build:
	go build -mod=vendor -o ${BINARY_NAME} ./cmd/${APP_NAME}/

# Run the application (build first, then run)
run: build
	./${BINARY_NAME}

# Run without building (direct go run)
dev:
	go run ./cmd/${APP_NAME}/

# Build vendor dependencies
vendor:
	go mod tidy
	go mod vendor

# Clean built binaries
clean:
	rm -f ${BINARY_NAME}

swag:
	swag init -g ./cmd/${APP_NAME}/main.go -o docs

# Docker commands for development
docker-compose-up:
	docker-compose -f deployment/dev/docker-compose.yml up -d

docker-compose-down:
	docker-compose -f deployment/dev/docker-compose.yml down

docker-compose-logs:
	docker-compose -f deployment/dev/docker-compose.yml logs -f

docker-compose-build:
	docker-compose -f deployment/dev/docker-compose.yml build

dev-up: docker-compose-build docker-compose-up

dev-down: docker-compose-down
