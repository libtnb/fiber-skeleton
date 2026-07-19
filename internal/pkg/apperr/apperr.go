// Package apperr standardizes client-facing errors: transports map the closed
// set of kinds — never individual codes — to a status, so modules add codes
// without touching shared files.
package apperr

import (
	"github.com/samber/oops"
)

// Kind classifies an application error for transport mapping.
type Kind string

const (
	// KindInvalid rejects malformed input (HTTP 400).
	KindInvalid Kind = "invalid"
	// KindUnauthorized requires authentication (HTTP 401).
	KindUnauthorized Kind = "unauthorized"
	// KindForbidden denies access (HTTP 403).
	KindForbidden Kind = "forbidden"
	// KindNotFound reports a missing resource (HTTP 404).
	KindNotFound Kind = "not_found"
	// KindConflict reports a state conflict, e.g. a taken name (HTTP 409).
	KindConflict Kind = "conflict"
	// KindUnprocessable rejects a semantically impossible request, e.g. a
	// reference to a missing row (HTTP 422).
	KindUnprocessable Kind = "unprocessable"
)

// kindKey stores the Kind in the oops error context.
const kindKey = "apperr.kind"

// New starts an error of the given kind: code is the stable identifier sent
// to clients, public the message safe to expose. Finish with .Wrap/.Errorf.
func New(kind Kind, code, public string) oops.OopsErrorBuilder {
	return oops.Code(code).Public(public).With(kindKey, kind)
}

func Invalid(code, public string) oops.OopsErrorBuilder { return New(KindInvalid, code, public) }
func Unauthorized(code, public string) oops.OopsErrorBuilder {
	return New(KindUnauthorized, code, public)
}
func Forbidden(code, public string) oops.OopsErrorBuilder { return New(KindForbidden, code, public) }
func NotFound(code, public string) oops.OopsErrorBuilder  { return New(KindNotFound, code, public) }
func Conflict(code, public string) oops.OopsErrorBuilder  { return New(KindConflict, code, public) }
func Unprocessable(code, public string) oops.OopsErrorBuilder {
	return New(KindUnprocessable, code, public)
}

// KindOf reports the error's Kind, or "" when it carries none.
func KindOf(err error) Kind {
	if oopsErr, ok := oops.AsError[oops.OopsError](err); ok {
		if kind, ok := oopsErr.Context()[kindKey].(Kind); ok {
			return kind
		}
	}
	return ""
}

// CodeOf reports the error's machine-readable code, or "" when absent.
func CodeOf(err error) string {
	if oopsErr, ok := oops.AsError[oops.OopsError](err); ok {
		code, _ := oopsErr.Code().(string)
		return code
	}
	return ""
}
