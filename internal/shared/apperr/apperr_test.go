package apperr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/libtnb/assert/must"

	"github.com/libtnb/fiber-skeleton/internal/shared/apperr"
)

var errSentinel = errors.New("sentinel")

func TestKindAndCodeSurviveWrapping(t *testing.T) {
	err := apperr.Conflict("user.name_taken", "name already taken").In("user").Wrap(errSentinel)

	must.Equal(t, apperr.KindOf(err), apperr.KindConflict)
	must.Equal(t, apperr.CodeOf(err), "user.name_taken")
	must.ErrorIs(t, err, errSentinel)

	wrapped := fmt.Errorf("placing order: %w", err)
	must.Equal(t, apperr.KindOf(wrapped), apperr.KindConflict)
	must.Equal(t, apperr.CodeOf(wrapped), "user.name_taken")
}

func TestPlainErrorsCarryNoKind(t *testing.T) {
	must.Equal(t, apperr.KindOf(errors.New("boom")), apperr.Kind(""))
	must.Empty(t, apperr.CodeOf(errors.New("boom")))
	must.Equal(t, apperr.KindOf(nil), apperr.Kind(""))
}

func TestHelpersSetTheirKinds(t *testing.T) {
	for helper, kind := range map[string]apperr.Kind{
		"invalid":       apperr.KindInvalid,
		"unauthorized":  apperr.KindUnauthorized,
		"forbidden":     apperr.KindForbidden,
		"not_found":     apperr.KindNotFound,
		"conflict":      apperr.KindConflict,
		"unprocessable": apperr.KindUnprocessable,
	} {
		var err error
		switch helper {
		case "invalid":
			err = apperr.Invalid("c", "p").Errorf("x")
		case "unauthorized":
			err = apperr.Unauthorized("c", "p").Errorf("x")
		case "forbidden":
			err = apperr.Forbidden("c", "p").Errorf("x")
		case "not_found":
			err = apperr.NotFound("c", "p").Errorf("x")
		case "conflict":
			err = apperr.Conflict("c", "p").Errorf("x")
		case "unprocessable":
			err = apperr.Unprocessable("c", "p").Errorf("x")
		}
		must.Equal(t, apperr.KindOf(err), kind, helper)
	}
}
