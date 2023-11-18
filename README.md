# Guardians of The Cluster

A dead simple multi-instance supported rate-limiting service that just works.
(Multi-instance supported part is not implemented correctly yet)

## How to start the server

1. Clone the project
2. Run `go mod tidy`
3. Run `make start`
4. Done! The server is running on port 3333

## How to run unit tests

1. Clone the project
2. Run `go mod tidy`
3. Run `make test`
4. Done! Unit tests should begin to run

## How to send requests?

You can use some basic curl commands to make requests to the server:

**GET request:**

```bash
curl localhost:3333/this-does-not-matter -H "X-User-Id: user-id"
```

**POST request:**

```bash
curl -X POST localhost:3333/this-does-not-matter -d "data" -H "X-User-Id: user-id"
```
