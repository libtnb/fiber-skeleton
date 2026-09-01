package service_test

import (
	"testing"

	"github.com/libtnb/assert/check"
	"github.com/libtnb/validator"

	"github.com/libtnb/fiber-skeleton/internal/shared/transport"
	"github.com/libtnb/fiber-skeleton/internal/user/service"
)

// TestCheckRules catches invalid validate tags at test time.
func TestCheckRules(t *testing.T) {
	v := validator.MustNew()

	check.NoError(t, v.Check[transport.Paginate]())
	check.NoError(t, v.Check[service.UserID]())
	check.NoError(t, v.Check[service.UserAdd]())
	check.NoError(t, v.Check[service.UserUpdate]())
}
