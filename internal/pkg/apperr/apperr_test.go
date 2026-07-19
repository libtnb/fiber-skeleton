package apperr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/pkg/apperr"
)

var errSentinel = errors.New("sentinel")

func TestKindAndCodeSurviveWrapping(t *testing.T) {
	err := apperr.Conflict("user.name_taken", "name already taken").In("user").Wrap(errSentinel)

	require.Equal(t, apperr.KindConflict, apperr.KindOf(err))
	require.Equal(t, "user.name_taken", apperr.CodeOf(err))
	require.ErrorIs(t, err, errSentinel)

	wrapped := fmt.Errorf("placing order: %w", err)
	require.Equal(t, apperr.KindConflict, apperr.KindOf(wrapped))
	require.Equal(t, "user.name_taken", apperr.CodeOf(wrapped))
}

func TestPlainErrorsCarryNoKind(t *testing.T) {
	require.Equal(t, apperr.Kind(""), apperr.KindOf(errors.New("boom")))
	require.Empty(t, apperr.CodeOf(errors.New("boom")))
	require.Equal(t, apperr.Kind(""), apperr.KindOf(nil))
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
		require.Equal(t, kind, apperr.KindOf(err), helper)
	}
}
