# Development

## Build from source

Wavekit server is built with `go` and `pnpm`. The web frontend is served by wavekit server as static files.
All static files are bundled into the server binary through `embed.go`.

1. Build frontend

```shell
make build-web
```

2. Build server

```shell
make build-server
```

## Note for development

1. `sql/migrations` and `sql/queries` are single source of truth for database schema and operation. Database related code will be generated from these files.

2. `wire/wire.go` helps to resolve dependency injection.

3. `web/api` is the single source of truth for API definition. Controller code in server and client code in web are generated from this file.

4. `internal/config` is the single source of truth for configuration. It will be used to load configuration from `config.yaml` or environment variables.

## Development flow

1. Generate code

```shell
make gen
```

2. Start development environment

```shell
make dev
```

3. Reload code 

```shell
make reload
```

4. Run unit tests

```shell
make ut
```

5. Run end-to-end test

```shell
make test
```
