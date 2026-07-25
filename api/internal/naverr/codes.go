package naverr

import (
	"net/http"
	"strings"

	"github.com/yueli-official/foundation/go/problem"
	"platform/gokit/errs"
)

var (
	CodeNotFound                 = errs.Register("nav.not_found", http.StatusNotFound)
	CodeForbidden                = errs.Register("nav.forbidden", http.StatusForbidden)
	CodeConflict                 = errs.Register("nav.conflict", http.StatusConflict)
	CodeNotInitialized           = errs.Register("nav.not_initialized", 503)
	CodeAuthorizationUnavailable = errs.Register("nav.authorization_unavailable", 503)
	CodeRevisionConflict         = errs.Register("nav.site_profile_revision_conflict", 412)
	CodePreconditionRequired     = errs.Register("nav.site_profile_precondition_required", 428)
)

func NotFound(id string) *errs.Coded {
	return errs.New(CodeNotFound, "navigation link not found", map[string]any{"id": id})
}

func RevisionConflict() *errs.Coded {
	return errs.New(CodeRevisionConflict, "site profile revision does not match", nil)
}

func PreconditionRequired() *errs.Coded {
	return errs.New(CodePreconditionRequired, "If-Match is required for site profile updates", nil)
}

func FaviconNotFound(id string) *errs.Coded {
	return errs.New(CodeNotFound, "navigation favicon not found", map[string]any{"id": id})
}

func Forbidden() *errs.Coded {
	return errs.New(CodeForbidden, "forbidden", nil)
}

func AuthorizationUnavailable() *errs.Coded {
	return errs.New(CodeAuthorizationUnavailable, "navigation authorization is unavailable", nil)
}

func Validation(field, code string, params map[string]any) *errs.Coded {
	return errs.New(errs.CommonValidationFailed, "validation failed", map[string]any{
		"details": []problem.Violation{{Pointer: "/" + strings.ReplaceAll(strings.ReplaceAll(field, "~", "~0"), "/", "~1"), Code: "validation." + code, Params: problem.Parameters(params)}},
	})
}

func Conflict(id string) *errs.Coded {
	return errs.New(CodeConflict, "navigation resource conflicts with existing or referenced data", map[string]any{"id": id})
}

func NotInitialized(resource string) *errs.Coded {
	return errs.New(CodeNotInitialized, "navigation site configuration is not initialized", map[string]any{"resource": resource})
}
