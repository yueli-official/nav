package naverr_test

import (
	"net/http"
	"testing"

	"github.com/yueli-official/foundation/go/problem"
	"github.com/yueli-official/nav/api/internal/naverr"
)

func TestCodesRegistered(t *testing.T) {
	cases := map[string]int{
		naverr.CodeNotFound:                 http.StatusNotFound,
		naverr.CodeForbidden:                http.StatusForbidden,
		naverr.CodeConflict:                 http.StatusConflict,
		naverr.CodeNotInitialized:           http.StatusServiceUnavailable,
		naverr.CodeAuthorizationUnavailable: http.StatusServiceUnavailable,
		naverr.CodeRevisionConflict:         http.StatusPreconditionFailed,
		naverr.CodePreconditionRequired:     http.StatusPreconditionRequired,
	}
	for code, want := range cases {
		descriptor, ok := naverr.DescriptorForCode(code)
		if !ok {
			t.Errorf("DescriptorForCode(%q) is missing", code)
			continue
		}
		if got := descriptor.Kind().Status(); got != want {
			t.Errorf("Status(%q) = %d, want %d", code, got, want)
		}
	}
}

func TestConstructorsCarryCode(t *testing.T) {
	cases := map[string]error{
		naverr.CodeNotFound:                 naverr.NotFound("link"),
		naverr.CodeForbidden:                naverr.Forbidden(),
		naverr.CodeConflict:                 naverr.Conflict("link"),
		naverr.CodeNotInitialized:           naverr.NotInitialized("site"),
		naverr.CodeAuthorizationUnavailable: naverr.AuthorizationUnavailable(),
		naverr.CodeRevisionConflict:         naverr.RevisionConflict(),
		naverr.CodePreconditionRequired:     naverr.PreconditionRequired(),
	}
	for want, err := range cases {
		value, ok, resolveErr := problem.FromError(err, "nav-error-inspection")
		if resolveErr != nil || !ok || value.Code != want {
			t.Errorf("FromError(%s) = %#v, %v, %v", want, value, ok, resolveErr)
		}
	}
}
