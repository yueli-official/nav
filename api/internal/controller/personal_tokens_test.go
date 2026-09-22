package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	"github.com/yueli-official/foundation/go/authorization"
	v1 "github.com/yueli-official/nav/api/api/v1"
	"github.com/yueli-official/nav/api/internal/navauthz"
)

func navPersonalContext(t *testing.T, user string, capabilities ...string) context.Context {
	t.Helper()
	scopes := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		scope, err := foundationauth.PersonalScope("nav-yueli-web", capability)
		if err != nil {
			t.Fatal(err)
		}
		scopes = append(scopes, scope)
	}
	endpoint := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		if err := json.NewEncoder(response).Encode(map[string]any{"userKey": user, "scopes": scopes}); err != nil {
			t.Error(err)
		}
	}))
	defer endpoint.Close()
	verifier, err := foundationauth.NewPersonalTokenVerifier(endpoint.URL, "nav-yueli-web", nil)
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), "pat_test")
	if err != nil {
		t.Fatal(err)
	}
	return foundationauth.NewContext(context.Background(), principal)
}

func navAuthorization(t *testing.T) (*navauthz.Service, *authorization.Memory) {
	t.Helper()
	runtime, err := authorization.NewMemory(authorization.MustCompile(navauthz.Definition()), authorization.MemoryOptions{
		RootScopeID:       navauthz.RootScopeID,
		ProtectedSubjects: []authorization.SubjectRef{{Kind: authorization.SubjectUser, ID: "TestA123"}},
		Constraints:       navauthz.ConstraintEvaluators(),
		Predicates:        navauthz.PredicateEvaluators(),
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, command := range []authorization.RegisterScopeCommand{
		{ID: navauthz.CategoryScopeID("dev"), Type: navauthz.ScopeCategory, ParentID: navauthz.RootScopeID},
		{ID: navauthz.GroupScopeID("go"), Type: navauthz.ScopeGroup, ParentID: navauthz.CategoryScopeID("dev")},
		{ID: navauthz.LinkScopeID("foundation"), Type: navauthz.ScopeLink, ParentID: navauthz.GroupScopeID("go")},
	} {
		if _, err := runtime.RegisterScope(context.Background(), command); err != nil {
			t.Fatal(err)
		}
	}
	return navauthz.New(runtime, nil), runtime
}

func TestPersonalScopeIntersectsCurrentNavRights(t *testing.T) {
	service, runtime := navAuthorization(t)
	admin := navPersonalContext(t, "TestA123", string(navauthz.CapabilityStructureManage))
	decision, err := service.Decide(admin, navauthz.CapabilityStructureManage, navauthz.RootScopeID, authorization.ResourceFacts{})
	if err != nil || !decision.Allowed {
		t.Fatalf("selected administrator capability denied: decision=%#v error=%v", decision, err)
	}
	decision, err = service.Decide(admin, navauthz.CapabilitySettingsManage, navauthz.RootScopeID, authorization.ResourceFacts{})
	if err != nil || decision.Allowed {
		t.Fatalf("unselected administrator capability allowed: decision=%#v error=%v", decision, err)
	}
	ordinary := navPersonalContext(t, "TestB234", string(navauthz.CapabilityStructureManage))
	decision, err = service.Decide(ordinary, navauthz.CapabilityStructureManage, navauthz.RootScopeID, authorization.ResourceFacts{})
	if err != nil || decision.Allowed {
		t.Fatalf("scope granted administrator capability: decision=%#v error=%v", decision, err)
	}

	adminSubject := authorization.SubjectRef{Kind: authorization.SubjectUser, ID: "TestA123"}
	curatorSubject := authorization.SubjectRef{Kind: authorization.SubjectUser, ID: "TestB234"}
	grant, err := runtime.Grant(context.Background(), authorization.GrantCommand{
		Actor: adminSubject, Target: curatorSubject, Role: navauthz.RoleCurator,
		ScopeID: navauthz.CategoryScopeID("dev"), Source: authorization.GrantSourceDirect,
	})
	if err != nil {
		t.Fatal(err)
	}
	curator := navPersonalContext(t, "TestB234", string(navauthz.CapabilityLinkUpdate))
	resource := navauthz.LinkResource("foundation", "TestB234")
	decision, err = service.Decide(curator, navauthz.CapabilityLinkUpdate, navauthz.LinkScopeID("foundation"), resource)
	if err != nil || !decision.Allowed {
		t.Fatalf("selected own-link capability denied: decision=%#v error=%v", decision, err)
	}
	if _, err := runtime.Revoke(context.Background(), authorization.RevokeCommand{Actor: adminSubject, GrantID: grant.ID}); err != nil {
		t.Fatal(err)
	}
	decision, err = service.Decide(curator, navauthz.CapabilityLinkUpdate, navauthz.LinkScopeID("foundation"), resource)
	if err != nil || decision.Allowed {
		t.Fatalf("revoked curator retained access: decision=%#v error=%v", decision, err)
	}
}

func TestManageLinkFilterRejectsUnselectedPersonalScopeBeforeQuery(t *testing.T) {
	service, _ := navAuthorization(t)
	access, err := service.ManageLinkFilter(navPersonalContext(t, "TestA123", string(navauthz.CapabilityLinkSubmit)))
	if err != nil {
		t.Fatal(err)
	}
	if access.All || access.OwnedAll || len(access.UnrestrictedGroups) != 0 || len(access.OwnedGroups) != 0 {
		t.Fatalf("unselected link.update returned access: %#v", access)
	}
}

func TestPersonalDirectoryRequiresTrustedIdentity(t *testing.T) {
	service, _ := navAuthorization(t)
	controller := NewPersonalPermissions("nav-yueli-web", service)
	trusted := foundationauth.NewContext(context.Background(), &foundationauth.Principal{
		SubjectKind: foundationauth.SubjectClient, ClientID: "identity-svc", Scopes: []string{foundationauth.PersonalPermissionsScope},
	})
	result, err := controller.GetPersonalPermissions(trusted, &v1.PersonalPermissionsReq{UserKey: "TestA123"})
	if err != nil {
		t.Fatal(err)
	}
	keys := permissionKeys(result.Items)
	for _, key := range []string{
		string(navauthz.CapabilityStructureManage), string(navauthz.CapabilityHealthCheckRun), string(navauthz.CapabilitySettingsManage),
	} {
		if !slices.Contains(keys, key) {
			t.Errorf("administrator directory missing %q: %v", key, keys)
		}
	}
	untrusted := foundationauth.NewContext(context.Background(), &foundationauth.Principal{Subject: "TestA123", SubjectKind: foundationauth.SubjectUser})
	if _, err := controller.GetPersonalPermissions(untrusted, &v1.PersonalPermissionsReq{UserKey: "TestA123"}); err == nil {
		t.Fatal("untrusted directory request accepted")
	}
}

func TestPersonalRoutesAreExplicitAndCapabilityBound(t *testing.T) {
	if !allowsPersonalRoute(navPersonalContext(t, "TestA123", string(navauthz.CapabilityLinkSubmit)), "GET", "/api/v1/nav/catalog") {
		t.Fatal("public catalog denied to valid personal token")
	}
	for _, test := range []struct{ method, path, capability string }{
		{"GET", "/api/v1/admin/nav/links", string(navauthz.CapabilityLinkUpdate)},
		{"POST", "/api/v1/admin/nav/links", string(navauthz.CapabilityLinkSubmit)},
		{"DELETE", "/api/v1/admin/nav/links/link", string(navauthz.CapabilityLinkModerate)},
		{"GET", "/api/v1/admin/nav/structure", string(navauthz.CapabilityStructureManage)},
		{"POST", "/api/v1/admin/nav/checks/run", string(navauthz.CapabilityHealthCheckRun)},
		{"PUT", "/api/v1/admin/nav/settings", string(navauthz.CapabilitySettingsManage)},
	} {
		if !allowsPersonalRoute(navPersonalContext(t, "TestA123", test.capability), test.method, test.path) {
			t.Errorf("selected capability denied: %s %s", test.method, test.path)
		}
	}
	all := navPersonalContext(t, "TestA123",
		string(navauthz.CapabilityLinkSubmit), string(navauthz.CapabilityLinkUpdate), string(navauthz.CapabilityLinkModerate),
		string(navauthz.CapabilityStructureManage), string(navauthz.CapabilityHealthCheckRun), string(navauthz.CapabilitySettingsManage),
	)
	for _, path := range []string{
		"/api/v1/me", "/api/v1/admin/nav/members", "/api/v1/internal/personal-token/permissions",
		"/api/v1/authorization/manage/console", "/api/v1/admin/nav/links//", "/api/v1/admin/nav/groups/../settings",
	} {
		for _, method := range []string{"GET", "POST", "PATCH", "DELETE", "PUT"} {
			if allowsPersonalRoute(all, method, path) {
				t.Errorf("unexpected route allowed: %s %s", method, path)
			}
		}
	}
}

func permissionKeys(items []foundationauth.PersonalPermission) []string {
	keys := make([]string, 0, len(items))
	for _, item := range items {
		keys = append(keys, item.Key)
	}
	return keys
}
