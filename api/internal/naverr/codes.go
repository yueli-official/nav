package naverr

import (
	"net/http"

	"platform/gokit/errs"
	"platform/gokit/response"
)

var (
	CodeNotFound  = errs.Register("nav.not_found", http.StatusNotFound)
	CodeForbidden = errs.Register("nav.forbidden", http.StatusForbidden)
	CodeConflict  = errs.Register("nav.conflict", http.StatusConflict)
)

func NotFound(id string) *errs.Coded {
	return errs.New(CodeNotFound, "navigation link not found", map[string]any{"id": id})
}

func Forbidden() *errs.Coded {
	return errs.New(CodeForbidden, "forbidden", nil)
}

func Validation(field, code string, params map[string]any) *errs.Coded {
	return errs.New(errs.CommonValidationFailed, "validation failed", map[string]any{
		"details": []response.ValidationDetail{{Field: field, Code: code, Params: params}},
	})
}

func Conflict(id string) *errs.Coded {
	return errs.New(CodeConflict, "navigation link already exists", map[string]any{"id": id})
}
