package catalog

import (
 "errors"
 "testing"
 "github.com/yueli-official/foundation/go/problem"
 "github.com/yueli-official/foundation/go/siteprofile"
 "github.com/yueli-official/nav/api/internal/naverr"
)

func TestProfileFailuresKeepAllPointersAndUnderlyingCause(t *testing.T) {
 original := &siteprofile.ValidationError{Diagnostics: []siteprofile.Diagnostic{
  {Path:"identity.name",Code:"required",Message:"internal diagnostic"},
  {Path:"footer.legal[0].href",Code:"link_invalid",Message:"provider text"},
 }}
 result := naverr.MapCause(mapSiteProfileError(original))
 var typed *siteprofile.ValidationError
 if !errors.As(result, &typed) || typed != original { t.Fatal("provider cause lost") }
 wire, ok, err := problem.FromError(result, "nav-test")
 if !ok || err != nil || len(wire.Violations) != 2 { t.Fatalf("wire: %#v, %v", wire, err) }
 if wire.Violations[0].Pointer != "/profile/identity/name" || wire.Violations[1].Pointer != "/profile/footer/legal/0/href" { t.Fatalf("pointers: %#v", wire.Violations) }
}
