# Gin Backend Template Application

This is a backend application built using the [Gin Web Framework](https://github.com/gin-gonic/gin) in Go. The application uses PostgreSQL as the database, GORM for ORM, Air for hotreload for the development.

## Docker-compose services:

1. api - Gin application.
2. postgresql - PostgreSQL server.
3. pgadmin - PGAdmin to inspect DB stracture and data.
4. mailhog - fake SMTP server for mail interception during the development.
5. stunnel - Proxy to add encryption for mailhog. gomail need activated TLS.

## Features

- Docker Compose for easy setup. Just one command (`docker compose up`) to launch all services.
- Environment variable-based configuration (`.env`).
- Air (github.com/air-verse/air) for hot reload of backend.
- Two GORM ORM models: User and Item.
- CRUD operations for items and users.
- User registration and authentication with JWT.
- Swagger documentation for API.
- Email sending feature with SMTP, configurable for Mailhog in development environment.
- Custom validators.
- Tests backend. Set up of test is done with following features:
- - Tests use separated DB on the PostgreSQL server.
- - Factories for User and Item to faciltate creation of data.
- - API Test Client with or without authentication for API testing.

## Getting starting.

You need installed Docker v24 (or greater).

1. Copy .env.example to .env.
2. Run services using docker compose:
   `docker compose up`

Once started, the following services are available:

1. http://localhost:8083/ - Backend
2. http://localhost:8083/swagger/index.html - API Documentation Swagger
3. http://localhost:8085/ - MailHog
4. http://localhost:8086/ - PGAdmin4

Attention. Documentation Swagger is not generated automatically. You need generate it when annotations are updated (see command below in "Useful commands")

## Useful commands

Build (or rebuild) the containers and launch all services:

```
docker-compose up --build --force-recreate
```

Execute a command inside a running container for service todo-my-service (here are a running a /bin/sh inside todo-api container):

```
docker-compose exec todo-api /bin/sh
```

Run tests inside todo-api container:

```
docker-compose exec todo-api go test ./internal/tests/api_tests/
```

Generate documentation Swagger:

```
cd backend
swag init
```

Or generate documentation Swagger in running api container:

```
docker-compose exec todo-api swag init
```
