package service_test

import (
	"testing"

	"github.com/libtnb/validator"
	"github.com/stretchr/testify/assert"

	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
	"github.com/libtnb/fiber-skeleton/internal/user/service"
)

// TestCheckRules catches invalid validate tags at test time; register custom
// rules here if a request uses them.
func TestCheckRules(t *testing.T) {
	v := validator.NewValidator()

	for _, req := range []any{
		transport.Paginate{},
		service.UserID{},
		service.UserAdd{},
		service.UserUpdate{},
	} {
		assert.NoError(t, v.CheckRules(req), "%T has an invalid validate tag", req)
	}
}
