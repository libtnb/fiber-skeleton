# cmd

One directory per binary:

- `app` — the HTTP server
- `cli` — management commands (`migrate`, `user`, ...)
- `gen` — generates CRUD modules and migrations

Entry points stay minimal: call the generated Wire initializer and run.
