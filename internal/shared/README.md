# shared

Contracts shared by every module:

- `transport` — request binding, response envelopes, endpoint declarations
- `apperr` — typed application errors
- `event` — the event bus interface
- `registry` — typed Wire contribution collections (routes, commands, jobs, subscriptions, health checks)
- `job` — the scheduler contribution type

Packages here depend only on each other; the architecture test enforces it.
