# Go REST API with Gin

Simple example using gin for a REST API connected to a postgres database using pgx.

## Environment

Set postgres database connection parameters:

```bash
export DB_HOST=...
export DB_USER=...
export DB_PASSWORD=...
export DB_NAME=...
```

## Dev mode

```bash
make start-dev
```

## Dockerized version

```bash
make build
make start
```
