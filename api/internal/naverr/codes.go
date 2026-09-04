// Package naverr declares Nav's immutable public Problem contract.
package naverr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yueli-official/foundation/go/problem"
)

var (
	DescriptorRateLimited = descriptor("common.rate_limited", http.StatusTooManyRequests)
	DescriptorValidation  = descriptor("common.validation_failed", http.StatusBadRequest)
	DescriptorInternal    = descriptor("common.internal", http.StatusInternalServerError)
)

func mapped(code string, params problem.Parameters, violations ...problem.Violation) error {
	value, ok := DescriptorForCode(code)
	if !ok {
		return fmt.Errorf("nav public error code is not declared: %s", code)
	}
	result, err := problem.NewError(value, params, violations...)
	if err != nil {
		return fmt.Errorf("nav public error %s: %w", code, err)
	}
	return result
}

func NotFound(id string) error {
	return mapped(CodeNotFound, map[string]any{"id": id})
}

func RevisionConflict() error {
	return mapped(CodeRevisionConflict, nil)
}

func PreconditionRequired() error {
	return mapped(CodePreconditionRequired, nil)
}

func FaviconNotFound(id string) error {
	return mapped(CodeNotFound, map[string]any{"id": id, "resource": "favicon"})
}

func Forbidden() error {
	return mapped(CodeForbidden, nil)
}

func AuthorizationUnavailable() error {
	return mapped(CodeAuthorizationUnavailable, nil)
}

func MembershipUnavailable() error {
	return mapped(CodeMembershipUnavailable, nil)
}

func MembershipSuspended() error {
	return mapped(CodeMembershipSuspended, nil)
}

func Validation(field, code string, params map[string]any) error {
	violation := problem.Violation{
		Pointer: "/" + strings.ReplaceAll(strings.ReplaceAll(field, "~", "~0"), "/", "~1"),
		Code:    "validation." + code,
		Params:  validationParams(params),
	}
	return mapped(CodeInvalidInput, nil, violation)
}

// Only declarative constraint metadata may cross the public boundary.
func validationParams(params map[string]any) problem.Parameters {
	safe := problem.Parameters{}
	for name, value := range params {
		switch name {
		case "min", "max":
			if n, ok := value.(int); ok && n >= 0 {
				safe[name] = n
			}
		case "allowed":
			if values, ok := value.([]string); ok && len(values) <= 32 {
				bounded := true
				for _, item := range values {
					if len(item) > 120 {
						bounded = false
					}
				}
				if bounded {
					safe[name] = append([]string(nil), values...)
				}
			}
		}
	}
	return safe
}

func Conflict(id string) error {
	return mapped(CodeConflict, map[string]any{"id": id})
}

func NotInitialized(resource string) error {
	return mapped(CodeNotInitialized, map[string]any{"resource": resource})
}
