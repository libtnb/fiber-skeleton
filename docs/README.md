# docs

Hand-written documentation lives here. The API documentation is generated at
runtime: with `http.docs` enabled, the server serves the OpenAPI 3.1 document
at `/openapi.json` and a browsable UI at `/docs`.

To document an endpoint, attach `transport.Describe[Request, Response](status)`
or `transport.DescribeNoBody[Request](status)` to its route contribution (see
`internal/user/service/route.go`). Endpoints without a `Document` callback stay
out of the document.
