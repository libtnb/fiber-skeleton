# fiber-skeleton

Unlike [chi-skeleton](https://github.com/libtnb/chi-skeleton), this skeleton uses the incredibly fast [Fiber](https://gofiber.io/) framework, which is generally recommended.

## Features

- **Fiber v3** with sensible server hardening (timeouts, body/header limits) and a global middleware stack
- **Dependency injection** via [samber/do](https://github.com/samber/do): lazy generic container, reverse-order shutdown of resources, health checks feeding `/readyz` — and no code generation
- **Strongly-typed configuration** ([koanf](https://github.com/knadh/koanf)) with `APP_*` environment overrides, validated at startup
- **Graceful shutdown** on SIGINT/SIGTERM (drains requests and cron jobs) and **zero-downtime upgrades** on SIGHUP ([graceful](https://github.com/libtnb/graceful))
- **Structured logging** with [slog](https://pkg.go.dev/log/slog) on a rotating file writer ([logrotate](https://github.com/libtnb/logrotate)), stdout, or both
- **Request binding + validation** ([validator](https://github.com/libtnb/validator)) with a boolean rule DSL and i18n messages; each service holds its validator — no package-global state
- **Typed application errors** (`internal/pkg/apperr`): a closed set of error kinds maps to HTTP statuses in one place, so a module can add error codes without touching any shared file
- **Scheduled jobs** ([cron](https://github.com/libtnb/cron)) with panic recovery and overlap skipping; modules contribute jobs through `registry` without importing the boot wiring
- **rio + SQLite** ([rio](https://github.com/go-rio/rio), swap in MySQL/PostgreSQL freely) with versioned migrations ([migrate](https://github.com/go-rio/migrate)): schema as Go code, automatic rollbacks, `migrate status`/`rollback` commands
- **Code generator** (`cmd/gen`) that scaffolds a full CRUD module in one command
- **OpenAPI 3.1 docs from validate tags** — schemas and constraints generated from the validator rules, served with a Scalar UI at `/docs`
- **Tests included**: handler tests against mocked repos ([mockery](https://github.com/vektra/mockery)), data-layer tests on a real migrated SQLite, validate-tag linting, a container wiring test, and an architecture test that fails CI when a module crosses another module's boundary

## Quick start

Requires Go 1.25+.

```bash
git clone https://github.com/libtnb/fiber-skeleton my-app && cd my-app
make init   # copies config/config.example.yml to config/config.yml
make run    # or `make dev` for hot reload via air
```

The API listens on `:3000` by default: `curl localhost:3000/users`.

## Design

* `cmd` stores the entry point of each application, one directory per binary (`app`, `cli`, `gen`)
* `config` stores the configuration files
* `docs` stores hand-written documentation; the OpenAPI document is generated at runtime
* `internal` stores the application code: one directory per business module plus the shared layers below
* `internal/pkg` stores the contracts shared by every module (transport helpers, apperr, event bus, registry, job)
* `mocks` stores the generated mocks, one package per module (`mocks/user/biz`, `mocks/order/biz`)
* `storage` stores files generated while the application runs (logs, the SQLite database)
* `web` stores the front-end code of the application
* go.mod and go.sum manage dependencies — including the pinned `tool` directives (mockery)

Each business module (`internal/user`, `internal/order`, ...) follows the three-layer design of [Kratos](https://go-kratos.dev/):

* **biz** holds domain models, repository interfaces and **usecases** — transport-independent business logic
* **data** implements the repositories against the database
* **service** adapts HTTP: binds/validates requests, delegates to usecases, shapes responses

Because usecases are transport-independent, the HTTP handlers, the CLI commands (each module's `service/command.go`) and the cron jobs all share the same business logic instead of each talking to the database on their own.

Wiring follows a contribution model: every package exposes a `Package` list of lazy providers, and transports (routes, CLI commands, jobs, subscribers) are registered under naming conventions (`routes:*`, `commands:*`, `jobs:*`, `subscribers:*`) that assemblers collect at startup — adding a module never touches shared files beyond one line per Package list.

The boundaries are enforced, not aspirational: `TestModuleBoundaries` (`internal/app/arch_test.go`) parses the import graph and fails when a module reaches another module past its `biz` package, or any module imports the composition layers. Cross-module needs are expressed as interfaces in the consumer's biz package (see `order/biz.Users`) and adapted over the other module's public usecase in `data` — swap that adapter for an RPC client and the module splits into a service without touching its business logic.

## Configuration

`config/config.yml` is loaded first (override the path with `APP_CONFIG`), then any `APP_*` environment variable wins over the file. A double underscore separates nesting levels:

```bash
APP_HTTP__ADDRESS=:8080 APP_LOG__OUTPUT=stdout ./app
```

Configuration is parsed into a struct and validated at startup — a missing key or a bad value fails fast instead of panicking mid-request.

## Scheduled jobs

Add a job where it belongs — in the module that owns it: one `job.Fn` (`internal/pkg/job`) contribution per job, registered with one `do.LazyNamed(registry.JobPrefix+"name", ...)` line in the module's `Package` list (see `bootstrap.Heartbeat` for the shape). Specs support an optional seconds field, `@every 30s` descriptors and per-entry timezones. Jobs receive a `context.Context` that is cancelled on shutdown; panics are recovered and overlapping runs are skipped.

## Code generation

```bash
make gen name=article    # or: go run ./cmd/gen article
```

generates the biz entity + repo interface, data repository, service handlers, route contribution, request structs and a migration for a new module, then prints the remaining wiring: one line per Package list.

## Development

```bash
make help       # list all targets
make generate   # regenerate mocks after changing interfaces
make lint       # golangci-lint
make test       # go test -race with coverage
make build      # static binaries in bin/ with the version injected
```

A `Dockerfile` is included; mount `config/` and `storage/` when running.

## OpenAPI documentation

Every documented endpoint declares `Request`/`Response` samples in its route
contribution; schemas, parameters and constraints are generated from the very
same `validate` tags that enforce them ([validator/contrib/openapi](https://github.com/libtnb/validator/tree/main/contrib/openapi)) — `min:3 && max:255` becomes `minLength`/`maxLength`, `in:a,b` becomes an enum, and the two can never drift apart. With `http.docs: true` the app serves the OpenAPI 3.1 document at `/openapi.json` and a [Scalar](https://github.com/scalar/scalar) UI at `/docs`.

## Observability

- `/healthz` (liveness) and `/readyz` (readiness, pings the DB) are wired for containers and load balancers; the Dockerfile ships a matching `HEALTHCHECK`.
- Access logs and application logs share one slog logger, one format and one `request_id`, so a request can be traced across both.
- Set `http.debug_address` (e.g. `127.0.0.1:6060`) to serve `net/http/pprof` and `expvar` on a **separate private port** — profiling in production without exposing it on the API port.
- Errors returned by handlers and by the framework itself (404, 405, 413, panics) all leave through one handler in the same JSON shape; 5xx details go to the log, not the client.

## Error model

A usecase creates client-facing errors through `internal/pkg/apperr`:

```go
apperr.Conflict("user.name_taken", "name already taken").In("user").Wrap(ErrNameTaken)
```

The **kind** (conflict, not_found, invalid, ...) is a closed set that `transport.ErrorFrom` maps to an HTTP status; the **code** and public message travel to the client; everything else — stack trace, domain, attributes — goes to the log. Adding a module adds codes, never a new case in shared code. Errors without a kind are unexpected: the client sees a bare 500 and the details stay in the log.

## Serving a frontend

Put your built frontend under `web/` and serve it with fiber's static middleware in `internal/server/server.go` (`NewRouter`):

```go
r.Get("/*", static.New("./web/dist"))
```

## Graceful lifecycle

| Signal | Behavior |
|---|---|
| SIGINT / SIGTERM | stop accepting connections, drain in-flight requests and cron jobs (30s cap), close DB and log writer |
| SIGHUP (non-Windows) | zero-downtime binary upgrade via [graceful](https://github.com/libtnb/graceful) |

## Credits

The development of this project refers to the following projects, I would like to express my gratitude:

* [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
* [Kratos](https://go-kratos.dev/)
* [Goravel](https://github.com/goravel/goravel)
* [Fiber backend template](https://github.com/create-go-app/fiber-go-template)
* [GinSkeleton](https://github.com/qifengzhang007/GinSkeleton)
* [gin-layout](https://github.com/wannanbigpig/gin-layout)
