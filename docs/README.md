# docs

API documentation is generated automatically from the route declarations — no
separate codegen step. When `http.docs` is enabled (see
`config/config.example.yml`), the running server serves:

- `GET /openapi.json` — the OpenAPI 3.1 document, built from every endpoint's
  `Request`/`Response` samples and their `validate` tags
  (see `internal/server/endpoint.go`).
- `GET /docs` — a browsable UI for that document.

To document a new endpoint, set its `Request`/`Response` types on the module's
route contribution (e.g. `internal/user/service/route.go`); an endpoint with
neither stays out of the document.

This directory also holds any additional hand-written documents.
