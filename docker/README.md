# Docker

## Docker Compose

### Compose up

```
docker compose -f docker/compose.yaml up -d
```

### Compose down

```
docker compose -f docker/compose.yaml down
```

## Managing Migrations

### Migrate up

Run this to apply all up migrations
```
docker run --rm \
    -v './migrations:/migrations' \
    --network host \
    migrate/migrate -path /migrations/ -database 'postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:$POSTGRES_PORT/budgit?sslmode=disable' up
```

### Migrate down

Run this to apply all down migrations
```
docker run --rm \
    -v './migrations:/migrations' \
    --network host \
    migrate/migrate -path /migrations/ -database 'postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:$POSTGRES_PORT/budgit?sslmode=disable' down --all
```

### Migration as a docker-compose container

Add the following container to your docker compose:
```
migrate:
  image: migrate/migrate
  depends_on:
    postgres:
      condition: service_healthy
  env_file: ".env"
  volumes:
    - ./../migrations:/migrations
  command: ["-source", "file://migrations", "-database",  "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@postgres:$POSTGRES_PORT/budgit?sslmode=disable", "up"]
```