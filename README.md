# Guardian

A dead simple rate-limiting service that just works.

## How to run

1. Clone the project
2. Run `go mod tidy`
3. Run `make start`
4. Done! The server is running on port 3333

## How to send requests?

You can use some basic curl commands to make requests to the server:

**GET request:**

```bash
curl localhost:3333/this-does-not-matter
```

**POST request:**

```bash
curl -X POST localhost:3333/this-does-not-matter -d "the-data-string-that-you-want-to-send-with-your-request-why-you-are-still-reading-this?"
```
