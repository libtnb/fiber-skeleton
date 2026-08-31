# platform

Infrastructure assembly, invisible to business modules:

- `bootstrap` — providers for logger, database, crypter, event bus, cron, migrator
- `conf` — configuration loading and validation
- `server` — HTTP server, router, middleware, health endpoints

Only `internal/app` may import these packages; the architecture test enforces it.
