// Package server assembles the HTTP layer from the modules' route
// contributions and serves the non-domain endpoints.
package server

import (
	"regexp"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator/contrib/openapi"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
)

// Version is the build version, injected by main; the OpenAPI document carries it.
type Version string

// HTTP registers every "routes:*" contribution on r.
func HTTP(i do.Injector, r fiber.Router) error {
	groups, err := registry.Collect[transport.Endpoints](i, registry.RoutePrefix)
	if err != nil {
		return err
	}
	for _, endpoints := range groups {
		for _, e := range endpoints {
			r.Add([]string{e.Method}, e.Path, e.Handler)
		}
	}

	return nil
}

var pathParams = regexp.MustCompile(`:([A-Za-z0-9_]+)`)

// SpecJSON assembles the OpenAPI 3.1 document from every documented endpoint.
func SpecJSON(i do.Injector, title string) ([]byte, error) {
	version := "dev"
	if v, err := do.Invoke[Version](i); err == nil && v != "" {
		version = string(v)
	}
	g := openapi.New(title, version,
		openapi.WithType(time.Time{}, &openapi.Schema{Type: "string", Format: "date-time"}),
	)

	groups, err := registry.Collect[transport.Endpoints](i, registry.RoutePrefix)
	if err != nil {
		return nil, err
	}
	for _, endpoints := range groups {
		for _, e := range endpoints {
			if e.Request == nil && e.Response == nil {
				continue
			}
			err := g.Add(e.Method, pathParams.ReplaceAllString(e.Path, "{$1}"), openapi.Op{
				Summary:  e.Summary,
				Tags:     e.Tags,
				Request:  e.Request,
				Response: e.Response,
				Status:   e.Status,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return g.JSON()
}
