# Layered Architecture

A modern Golang backend application built with clean architecture principles, providing a solid foundation for building scalable web services.

## Features

- **Clean Architecture**: Structured using domain-driven design principles
- **RESTful API**: Built with Gin web framework
- **Dependency Injection**: Uses Google Wire for dependency management

- **Caching**: Redis for performance optimization
- **Message Broker**: Kafka integration (optional)
- **API Documentation**: Swagger for API documentation
- **Containerization**: Docker and Docker Compose support
- **Logging**: Structured logging with Zap

## Getting Started

### Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Make (optional, for convenience commands)

### Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd app
   ```

2. Install dependencies:
   ```
   go mod tidy
   go mod vendor
   ```
   Or simply:
   ```
   make vendor
   ```

3. Setup environment variables:
   ```
   cp example.env .env
   # Edit .env file with your configuration
   ```

4. Start required services using Docker:
   ```
   docker-compose up -d redis
   ```
5. Run database migrations manually:
   ```
   # See migrations/README.md for detailed instructions
   # See migrations/README.md for detailed instructions
   # Ensure you have PostgreSQL installed and running
   # psql -h 127.0.0.1 -p 5432 -U postgres -d kado -f migrations/initial_version_go.sql
   ```

### Running the Application

Run the application locally:
```
make run
```

Or using Docker:
```
docker-compose up
```

### API Testing

Test a sample endpoint:
```
curl http://localhost:8008/v1/user/test
```

Visit Swagger documentation:
```
http://localhost:8008/docs/index.html
```

## Development

### Dependency Injection

This project uses Google Wire for dependency injection:

1. Install Wire if not already installed:
   ```
   go install github.com/google/wire/cmd/wire@latest
   ```

2. Generate dependency injection code:
   ```
   cd internal/wire
   wire
   ```

### API Documentation
go install github.com/swaggo/swag/cmd/swag@latest
Generate Swagger documentation:
```
make swag
```
or
```
swag init -g ./cmd/server/main.go -o docs
```

## Configuration

The application now uses environment variables for configuration. All settings are loaded from `.env` file:

- Copy `example.env` to `.env`
- Modify the values according to your environment
- Environment variables take precedence over default values

## Project Structure

- `cmd/`: Application entry points
- `config/`: Legacy configuration files (kept for reference)
- `docs/`: API documentation
- `internal/`: Application core code
   - `initialize/`: Application initialization (DB, Redis, Router, etc.)
   - `middlewares/`: HTTP middlewares
   - `module/`: Domain modules
     - `user/`: User module
       - `constants/`: Module-specific constants
       - `controller/`: HTTP handlers
       - `dto/`: Data Transfer Objects
       - `model/`: Domain models
       - `repo/`: Data access layer
       - `router/`: Module routes
       - `service/`: Business logic
   - `wire/`: Dependency injection wiring
- `migrations/`: Database migrations (run manually)
- `pkg/`: Reusable packages
- `tests/`: Test files
- `example.env`: Example environment configuration

## Deployment

The application can be deployed using Docker:

```
docker build -t app .
docker run -p 8008:8008 app
```

Or using Docker Compose:

```
docker-compose up -d
```