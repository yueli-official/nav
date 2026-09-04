// Package navcause contains domain failures independent of HTTP representation.
package navcause

type Kind string

const (
	KindNotFound         Kind = "not_found"
	KindConflict         Kind = "conflict"
	KindNotInitialized   Kind = "not_initialized"
	KindRevisionConflict Kind = "revision_conflict"
	KindValidation       Kind = "validation"
)

type Cause struct {
	Cause      error
	Violations []Violation
	Kind       Kind
	Value      string
	Resource   string
	Field      string
	Rule       string
	Params     map[string]any
}
type Violation struct {
	Pointer string
	Rule    string
}

func (c *Cause) Unwrap() error { return c.Cause }

func (c *Cause) Error() string { return "nav: " + string(c.Kind) }
func (c *Cause) Is(target error) bool {
	other, ok := target.(*Cause)
	return ok && c.Kind == other.Kind
}
func NotFound(id string) error { return &Cause{Kind: KindNotFound, Value: id} }
func FaviconNotFound(id string) error {
	return &Cause{Kind: KindNotFound, Value: id, Resource: "favicon"}
}
func Conflict(id string) error { return &Cause{Kind: KindConflict, Value: id} }
func NotInitialized(resource string) error {
	return &Cause{Kind: KindNotInitialized, Resource: resource}
}
func RevisionConflict() error { return &Cause{Kind: KindRevisionConflict} }
func Validation(field, rule string, params map[string]any) error {
	return &Cause{Kind: KindValidation, Field: field, Rule: rule, Params: params}
}
