package naverr

import (
 "errors"
 "testing"
 "github.com/yueli-official/foundation/go/problem"
 "github.com/yueli-official/nav/api/internal/navcause"
)

func TestMappingRetainsCauseAndRejectsUndeclaredDiagnosticParams(t *testing.T) {
 cause := navcause.Validation("url", "absolute_http_url", map[string]any{"message":"secret provider body", "max":2048})
 result := MapCause(cause)
 var typed *navcause.Cause
 if !errors.As(result, &typed) || !errors.Is(result, cause) { t.Fatal("domain cause lost") }
 wire, ok, err := problem.FromError(result, "nav-test")
 if !ok || err != nil || wire.Code != CodeInvalidInput { t.Fatalf("wire = %#v, %v", wire, err) }
 if len(wire.Violations) != 1 || wire.Violations[0].Pointer != "/url" { t.Fatalf("violations = %#v", wire.Violations) }
 if _, exists := wire.Violations[0].Params["message"]; exists { t.Fatal("diagnostic leaked") }
 if wire.Violations[0].Params["max"] != 2048 { t.Fatal("constraint missing") }
}
