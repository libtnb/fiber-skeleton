# fiber-skeleton

[![Test](https://img.shields.io/github/actions/workflow/status/libtnb/fiber-skeleton/test.yml?branch=main&label=test)](https://github.com/libtnb/fiber-skeleton/actions/workflows/test.yml)
[![Lint](https://img.shields.io/github/actions/workflow/status/libtnb/fiber-skeleton/lint.yml?branch=main&label=lint)](https://github.com/libtnb/fiber-skeleton/actions/workflows/lint.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/libtnb/fiber-skeleton)](go.mod)
[![License](https://img.shields.io/github/license/libtnb/fiber-skeleton)](LICENSE)

A modular monolith skeleton for Go web applications, built on
[Fiber](https://gofiber.io/) v3.
Prefer `net/http`? See [chi-skeleton](https://github.com/libtnb/chi-skeleton).

## Features

- **Modular architecture** — Kratos-style `biz` / `data` / `service` layers per module, boundaries enforced by an architecture test
- **Compile-time DI** — [libtnb/wire](https://github.com/libtnb/wire) generates plain constructors; no runtime container, no reflection
- **Database** — [rio](https://github.com/go-rio/rio) on SQLite (MySQL/PostgreSQL drop in) with versioned migrations written in Go
- **Validation** — request binding with a boolean rule DSL and i18n messages ([validator](https://github.com/libtnb/validator))
- **OpenAPI 3.1** — generated from the same `validate` tags, served with a Scalar UI at `/docs`
- **Typed errors** — a closed set of error kinds maps to HTTP statuses in one place (`internal/shared/apperr`)
- **Logging** — structured [slog](https://pkg.go.dev/log/slog) to a rotating file and/or stdout; access logs share the logger
- **Scheduled jobs** — cron with panic recovery and overlap skipping
- **Event bus** — in-process pub/sub; modules contribute subscribers
- **WebSocket** — example echo endpoint at `/ws`
- **Lifecycle** — graceful shutdown on SIGINT/SIGTERM, zero-downtime upgrade on SIGHUP ([graceful](https://github.com/libtnb/graceful))
- **Code generation** — scaffold a CRUD module or a migration with one command
- **Tests** — handler tests on mocked repos, data-layer tests on a real SQLite, and an architecture test

## Getting started

Requires Go 1.27.

```bash
git clone https://github.com/libtnb/fiber-skeleton my-app && cd my-app
make init   # create config/config.yml from the example
make run    # or: make dev (hot reload via air)
```

The API listens on `:3000`:

```bash
curl localhost:3000/users
```

## Project layout

```
cmd/            entry points: app (HTTP server), cli (management commands), gen (generator)
config/         configuration files
docs/           hand-written docs; the OpenAPI document is generated at runtime
internal/
  app/          composition root: combines modules into the app and cli injectors
  migrations/   schema history, one file per migration
  platform/     infrastructure assembly: bootstrap (providers), conf, server
  shared/       contracts shared by every module: transport, apperr, event, registry, job
  user/         business module
  order/        business module
mocks/          generated repository mocks
storage/        runtime files: logs, SQLite database
web/            frontend code
```

## Architecture

Each business module follows the three-layer design of [Kratos](https://go-kratos.dev/):

- **biz** — domain models, repository interfaces and usecases; no transport or database code
- **data** — repository implementations
- **service** — transport adapters: bind and validate the request, call the usecase, shape the response

HTTP handlers, CLI commands and cron jobs all call the same usecases.

Each module declares a Wire `Module` that provides its constructors and
contributes routes, commands, jobs, subscribers and health checks.
`internal/app/wire.go` combines the modules into the `app` and `cli` injectors;
`make generate` writes the constructor code to `wire_gen.go`.

`TestModuleBoundaries` (`internal/app/arch_test.go`) fails the build when:

- a module imports another module except through its `biz` package
- a module imports `app`, `platform` or `migrations`
- `shared` or `platform` imports anything above its layer

Everything under `internal/` that is not `app`, `migrations`, `platform` or
`shared` is a business module.

To use another module, declare an interface in your own `biz` package and adapt
it over the other module's usecase in `data` (see `order/biz.Users`). Swapping
that adapter for an RPC client turns the module into a separate service without
touching its business logic.

## Configuration

`config/config.yml` is loaded first (override the path with `APP_CONFIG`), then
any `APP_*` environment variable wins over the file; a double underscore
separates nesting levels:

```bash
APP_HTTP__ADDRESS=:8080 APP_LOG__OUTPUT=stdout ./app
```

All keys are listed in [config/config.example.yml](config/config.example.yml).
Configuration is parsed into a struct and validated at startup.

## Database

Query shapes are declared once at package level, validated with `.Must()` and
reused concurrently; runtime values go to the terminal operations:

```go
var userByNameQuery = rio.From[biz.User]().Where("name = ?").Must()

exists, err := userByNameQuery.Exists(ctx, db, name)
```

### Migrations

The schema lives in `internal/migrations` as Go code
([migrate](https://github.com/go-rio/migrate)), one file per migration, applied
in file-name order:

```bash
make gen-migration name=add_email_to_users_table   # scaffold a migration
go run ./cmd/cli migrate                           # apply pending migrations
go run ./cmd/cli migrate status                    # list applied and pending
go run ./cmd/cli migrate rollback --step 1         # undo the most recent
```

`create_*_table` names scaffold a create-table body; `add_*_to_<table>_table`
names scaffold an alter body for that table.

## Code generation

```bash
make gen name=article
```

scaffolds a full module: biz entity + repository interface, data repository,
service handlers, request structs, documented routes, a create-table migration
and the Wire module. Then:

1. add `article.Module` to `ApplicationModule.Include` in `internal/app/wire.go`
2. run `make generate`

## Scheduled jobs

A job is a `job.Fn` provider in the module that owns it, contributed with
`Contribute[registry.Jobs](NewJob)` (see `bootstrap.Heartbeat`). Specs support
an optional seconds field, `@every 30s` and per-entry timezones. Jobs receive a
context cancelled on shutdown; panics are recovered and overlapping runs are
skipped.

## Error handling

Usecases build client-facing errors with `internal/shared/apperr`:

```go
apperr.Conflict("user.name_taken", "name already taken").In("user").Wrap(ErrNameTaken)
```

The kind (invalid, not_found, conflict, ...) maps to an HTTP status in
`transport.ErrorFrom`; the code and message travel to the client, everything
else — stack trace, domain, attributes — goes to the log. Errors without a kind
return a bare 500 with the details kept in the log.

## Observability

- `/healthz` (liveness) and `/readyz` (readiness, pings the database); the Dockerfile ships a matching `HEALTHCHECK`
- access and application logs share one logger and one `request_id`
- `http.debug_address` serves pprof and expvar on a separate private port
- framework errors (404, 405, 413, panics) return the same JSON envelope as the API

## Serving a frontend

Put the built frontend under `web/` and serve it from `NewRouter`
(`internal/platform/server/server.go`):

```go
r.Get("/*", static.New("./web/dist"))
```

## Deployment

`make build` produces static `app` and `cli` binaries in `bin/` with the
version injected. A `Dockerfile` is included; mount `config/` and `storage/`.

| Signal           | Behavior                                                                  |
| ---------------- | ------------------------------------------------------------------------- |
| SIGINT / SIGTERM | drain requests and jobs (30s cap), then close resources in reverse order  |
| SIGHUP           | zero-downtime binary upgrade                                              |

## Development

```bash
make help       # list all targets
make generate   # regenerate Wire constructors and mocks
make lint       # golangci-lint
make test       # go test -race with coverage
```

## Credits

Inspired by
[Standard Go Project Layout](https://github.com/golang-standards/project-layout),
[Kratos](https://go-kratos.dev/),
[Goravel](https://github.com/goravel/goravel),
[Fiber backend template](https://github.com/create-go-app/fiber-go-template),
[GinSkeleton](https://github.com/qifengzhang007/GinSkeleton) and
[gin-layout](https://github.com/wannanbigpig/gin-layout).

## License

[MIT](LICENSE)
