package naverr

import (
	"errors"
	"github.com/yueli-official/foundation/go/problem"
	"github.com/yueli-official/nav/api/internal/navcause"
)

// MapCause owns the Nav application mapping. The original cause is retained
// for errors.Is/As and never copied into the public parameters.
func MapCause(err error) error {
	var cause *navcause.Cause
	if !errors.As(err, &cause) {
		return err
	}
	var mappedError error
	switch cause.Kind {
	case navcause.KindNotFound:
		if cause.Resource == "favicon" {
			mappedError = FaviconNotFound(cause.Value)
		} else {
			mappedError = NotFound(cause.Value)
		}
	case navcause.KindConflict:
		mappedError = Conflict(cause.Value)
	case navcause.KindNotInitialized:
		mappedError = NotInitialized(cause.Resource)
	case navcause.KindRevisionConflict:
		mappedError = RevisionConflict()
	case navcause.KindValidation:
		if len(cause.Violations) == 0 {
			mappedError = Validation(cause.Field, cause.Rule, cause.Params)
		} else {
			violations := make([]problem.Violation, 0, len(cause.Violations))
			for _, violation := range cause.Violations {
				violations = append(violations, problem.Violation{Pointer: violation.Pointer, Code: "validation." + violation.Rule})
			}
			mappedError = mapped(CodeInvalidInput, nil, violations...)
		}
	default:
		return err
	}
	value, ok, resolveErr := problem.FromError(mappedError, "application-mapping")
	if resolveErr != nil || !ok {
		return mappedError
	}
	descriptor, ok := DescriptorForCode(value.Code)
	if !ok {
		return mappedError
	}
	result, wrapErr := problem.WrapError(descriptor, err, value.Params, value.Violations...)
	if wrapErr != nil {
		return wrapErr
	}
	return result
}
