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
	v := validator.MustNew()

	assert.NoError(t, v.Check[transport.Paginate]())
	assert.NoError(t, v.Check[service.UserID]())
	assert.NoError(t, v.Check[service.UserAdd]())
	assert.NoError(t, v.Check[service.UserUpdate]())
}
