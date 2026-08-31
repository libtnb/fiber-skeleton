// Package server assembles the HTTP layer from the modules' route contributions.
package server

import (
	"regexp"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator"
	"github.com/libtnb/validator/contrib/openapi"

	"github.com/libtnb/fiber-skeleton/internal/shared/registry"
)

// Version is the build version, injected by main; the OpenAPI document carries it.
type Version string

// NewVersion normalizes the build-time version for generated documentation.
func NewVersion(version string) Version {
	if version == "" {
		return "dev"
	}
	return Version(version)
}

// HTTP registers every route contribution on r.
func HTTP(groups registry.Routes, r fiber.Router) {
	for _, endpoints := range groups {
		for _, e := range endpoints {
			r.Add([]string{e.Method}, e.Path, e.Handler)
		}
	}
}

var pathParams = regexp.MustCompile(`:([A-Za-z0-9_]+)`)

// SpecJSON assembles the OpenAPI 3.1 document from every documented endpoint.
func SpecJSON(title string, version Version, validate *validator.Validator, groups registry.Routes) ([]byte, error) {
	g, err := openapi.New(title, string(version),
		openapi.WithValidator(validate),
		openapi.WithSchema[time.Time](&openapi.Schema{Type: "string", Format: "date-time"}),
	)
	if err != nil {
		return nil, err
	}
	for _, endpoints := range groups {
		for _, e := range endpoints {
			if e.Document == nil {
				continue
			}
			if err := e.Document(g, e.Method, pathParams.ReplaceAllString(e.Path, "{$1}"), e.Summary, e.Tags); err != nil {
				return nil, err
			}
		}
	}

	return g.JSON()
}
