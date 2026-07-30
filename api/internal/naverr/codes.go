// Package naverr declares Nav's immutable public Problem contract.
package naverr

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/yueli-official/foundation/go/problem"
)

const (
	CodeNotFound                 = "nav.not_found"
	CodeForbidden                = "nav.forbidden"
	CodeConflict                 = "nav.conflict"
	CodeNotInitialized           = "nav.not_initialized"
	CodeAuthorizationUnavailable = "nav.authorization_unavailable"
	CodeRevisionConflict         = "nav.site_profile_revision_conflict"
	CodePreconditionRequired     = "nav.site_profile_precondition_required"
)

var (
	DescriptorRateLimited = descriptor("common.rate_limited", http.StatusTooManyRequests)
	DescriptorValidation  = descriptor("common.validation_failed", http.StatusBadRequest)
	DescriptorInternal    = descriptor("common.internal", http.StatusInternalServerError)

	descriptors = map[string]problem.Descriptor{
		CodeNotFound:                 descriptor(CodeNotFound, http.StatusNotFound),
		CodeForbidden:                descriptor(CodeForbidden, http.StatusForbidden),
		CodeConflict:                 descriptor(CodeConflict, http.StatusConflict),
		CodeNotInitialized:           descriptor(CodeNotInitialized, http.StatusServiceUnavailable),
		CodeAuthorizationUnavailable: descriptor(CodeAuthorizationUnavailable, http.StatusServiceUnavailable),
		CodeRevisionConflict:         descriptor(CodeRevisionConflict, http.StatusPreconditionFailed),
		CodePreconditionRequired:     descriptor(CodePreconditionRequired, http.StatusPreconditionRequired),
	}
)

func descriptor(code string, status int) problem.Descriptor {
	return problem.MustDescriptor(
		problem.MustKind(code, status),
		"https://errors.yueli.dev/problems/"+code,
	)
}

func DescriptorForCode(code string) (problem.Descriptor, bool) {
	value, ok := descriptors[code]
	return value, ok
}

type CatalogEntry struct {
	Code   string `json:"code"`
	Status int    `json:"status"`
}

func Catalog() []CatalogEntry {
	result := make([]CatalogEntry, 0, len(descriptors))
	for code, value := range descriptors {
		result = append(result, CatalogEntry{Code: code, Status: value.Kind().Status()})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Code < result[j].Code })
	return result
}

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

func Validation(field, code string, params map[string]any) error {
	violation := problem.Violation{
		Pointer: "/" + strings.ReplaceAll(strings.ReplaceAll(field, "~", "~0"), "/", "~1"),
		Code:    "validation." + code,
		Params:  problem.Parameters(params),
	}
	result, err := problem.NewError(DescriptorValidation, nil, violation)
	if err != nil {
		return fmt.Errorf("nav validation error: %w", err)
	}
	return result
}

func Conflict(id string) error {
	return mapped(CodeConflict, map[string]any{"id": id})
}

func NotInitialized(resource string) error {
	return mapped(CodeNotInitialized, map[string]any{"resource": resource})
}
