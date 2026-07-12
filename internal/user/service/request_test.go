package service_test

import (
	"testing"

	"github.com/libtnb/validator"
	"github.com/stretchr/testify/assert"

	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
	"github.com/libtnb/fiber-skeleton/internal/user/service"
)

// TestCheckRules catches unknown rules, DSL syntax errors and bad static
// arguments in validate tags at test time instead of at request time.
// Register custom rules (exists, not_exists) here if a request uses them.
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
