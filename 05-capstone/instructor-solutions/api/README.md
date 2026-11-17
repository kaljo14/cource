# API

A standalone Go module. main.go connects to PostgreSQL and initializes the notes table; api.go handles HTTP validation; postgres.go performs parameterized SQL; server.go handles HTTP serving and shutdown.

From this directory: `go test ./...` and `docker build --build-arg VERSION=v1 -t campus-api:v1 .`.

Configuration: PGHOST, PGDATABASE, PGUSER and PGPASSWORD. The API listens on port 8080. See [startup instructions](../docs/START.md) for Kubernetes credentials, deployment and testing.
