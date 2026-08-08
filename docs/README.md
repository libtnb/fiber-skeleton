# docs

API documentation is generated automatically from the route declarations — no
separate codegen step. When `http.docs` is enabled (see
`config/config.example.yml`), the running server serves:

- `GET /openapi.json` — the OpenAPI 3.1 document, built from every endpoint's
  generic `Document` callback and the request type's `validate` tags
  (see `internal/pkg/transport/endpoint.go`).
- `GET /docs` — a browsable UI for that document.

To document a new endpoint, attach
`transport.Describe[Request, Response](status)` or
`transport.DescribeNoBody[Request](status)` to its route contribution (e.g.
`internal/user/service/route.go`). An endpoint without `Document` stays out of
the document.

This directory also holds any additional hand-written documents.
