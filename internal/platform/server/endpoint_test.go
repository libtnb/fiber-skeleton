package server

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/libtnb/validator"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/shared/registry"
	"github.com/libtnb/fiber-skeleton/internal/shared/transport"
)

type documentRequest struct {
	ID uint `uri:"id" validate:"required && number"`
}

type documentResponse struct {
	CreatedAt time.Time `json:"created_at"`
}

func TestSpecJSONUsesTypedSchemasAndNoBodyResponse(t *testing.T) {
	routes := registry.Routes{{
		{Method: http.MethodGet, Path: "/things/:id", Summary: "Get thing", Tags: []string{"thing"},
			Document: transport.Describe[documentRequest, documentResponse](http.StatusOK)},
		{Method: http.MethodDelete, Path: "/things/:id", Summary: "Delete thing", Tags: []string{"thing"},
			Document: transport.DescribeNoBody[documentRequest](http.StatusNoContent)},
	}}

	spec, err := SpecJSON("test", "v1", validator.MustNew(), routes)
	require.NoError(t, err)
	require.Contains(t, string(spec), `"format": "date-time"`)

	type response struct {
		Content map[string]any `json:"content"`
	}
	type operation struct {
		Responses map[string]response `json:"responses"`
	}
	var document struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
		Paths map[string]map[string]operation `json:"paths"`
	}
	require.NoError(t, json.Unmarshal(spec, &document))
	require.Equal(t, "v1", document.Info.Version)
	path := document.Paths["/things/{id}"]
	require.NotNil(t, path["get"].Responses["200"].Content)
	require.Nil(t, path["delete"].Responses["204"].Content)
}

func TestSpecJSONPropagatesGeneratorErrors(t *testing.T) {
	_, err := SpecJSON("", "v1", validator.MustNew(), nil)
	require.Error(t, err)

	routes := registry.Routes{{{
		Method:   http.MethodGet,
		Path:     "/things",
		Document: transport.Describe[documentRequest, documentResponse](0),
	}}}
	_, err = SpecJSON("test", "v1", validator.MustNew(), routes)
	require.Error(t, err)
}
