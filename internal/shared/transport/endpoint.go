package transport

import (
	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator/contrib/openapi"
)

// Documentation adds an endpoint to an OpenAPI generator.
type Documentation func(
	g *openapi.Generator,
	method string,
	path string,
	summary string,
	tags []string,
) error

// Describe documents an endpoint with JSON request and response types.
func Describe[Request, Response any](status int) Documentation {
	return func(
		g *openapi.Generator,
		method string,
		path string,
		summary string,
		tags []string,
	) error {
		return g.Add[Request](method, path,
			openapi.WithSummary(summary),
			openapi.WithTags(tags...),
			openapi.WithResponse[Response](status),
		)
	}
}

// DescribeNoBody documents an endpoint whose successful response has no body.
func DescribeNoBody[Request any](status int) Documentation {
	return Describe[Request, openapi.NoBody](status)
}

// Endpoint declares one HTTP endpoint; without Document it stays undocumented.
type Endpoint struct {
	Method   string
	Path     string
	Handler  fiber.Handler
	Summary  string
	Tags     []string
	Document Documentation
}

// Endpoints is one module's route contribution.
type Endpoints []Endpoint
