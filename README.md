# Subscriptions Service

REST API for managing user online subscriptions.

## Run (Docker Compose)

```bash
docker compose up --build
```

API: `http://localhost:8080`  
Swagger UI: `http://localhost:8080/docs/index.html`

## Makefile

Regenerate Swagger docs after updating handler annotations:

```bash
make swagger
```
