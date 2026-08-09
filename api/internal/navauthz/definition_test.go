package navauthz_test

import (
	"context"
	"testing"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	"github.com/yueli-official/foundation/go/authorization"
	"github.com/yueli-official/nav/api/internal/navauthz"
	"github.com/yueli-official/nav/api/internal/testidentity"
)

func TestServiceReconcilesEnabledAutomaticCuratorOnlyForNewMemberEvent(t *testing.T) {
	admin := authorization.SubjectRef{Kind: authorization.SubjectUser, ID: "TestA123"}
	module, err := authorization.NewMemory(
		authorization.MustCompile(navauthz.Definition()),
		authorization.MemoryOptions{
			RootScopeID: navauthz.RootScopeID, ProtectedSubjects: []authorization.SubjectRef{admin},
			Constraints: navauthz.ConstraintEvaluators(), Predicates: navauthz.PredicateEvaluators(),
		},
	)
	if err != nil {
		t.Fatalf("NewMemory() error = %v", err)
	}
	ctx := context.Background()
	draft, err := module.CreatePolicyDraft(ctx, authorization.CreatePolicyDraftCommand{
		Actor: admin, ScopeID: navauthz.RootScopeID, ExpectedActiveRevision: 1,
	})
	if err != nil {
		t.Fatalf("CreatePolicyDraft() error = %v", err)
	}
	if _, err := module.SetAutomaticRuleEnabled(ctx, authorization.SetAutomaticRuleEnabledCommand{
		Actor: admin, Revision: draft.Number,
		Rule: navauthz.AutomaticRegistrationCuratorKey, Enabled: true,
	}); err != nil {
		t.Fatalf("SetAutomaticRuleEnabled() error = %v", err)
	}
	if _, err := module.ActivatePolicy(ctx, authorization.ActivatePolicyCommand{
		Actor: admin, Revision: draft.Number, ExpectedActiveRevision: 1,
	}); err != nil {
		t.Fatalf("ActivatePolicy() error = %v", err)
	}
	service := navauthz.New(module, nil)
	userContext := foundationauth.NewContext(ctx, testidentity.User(t, "TestB234", nil, nil))
	access, err := service.EffectiveAccess(userContext)
	if err != nil {
		t.Fatalf("EffectiveAccess() error = %v", err)
	}
	for _, grant := range access.Grants {
		if grant.Role == navauthz.RoleCurator {
			t.Fatalf("ordinary access unexpectedly reconciled grant: %#v", access.Grants)
		}
	}
	if err := service.ReconcileNewMember(userContext); err != nil {
		t.Fatalf("ReconcileNewMember() error = %v", err)
	}
	access, err = service.EffectiveAccess(userContext)
	if err != nil {
		t.Fatalf("EffectiveAccess() after join error = %v", err)
	}
	found := false
	for _, grant := range access.Grants {
		found = found || grant.Role == navauthz.RoleCurator && grant.Source == authorization.GrantSourceAutomatic
	}
	if !found {
		t.Fatalf("EffectiveAccess() grants = %#v, want automatic curator", access.Grants)
	}
}

func TestDefinitionSupportsSubmitterQueriesAndScopedDelegation(t *testing.T) {
	admin := authorization.SubjectRef{Kind: authorization.SubjectUser, ID: "admin"}
	curator := authorization.SubjectRef{Kind: authorization.SubjectUser, ID: "curator"}
	other := authorization.SubjectRef{Kind: authorization.SubjectUser, ID: "other"}
	module, err := authorization.NewMemory(
		authorization.MustCompile(navauthz.Definition()),
		authorization.MemoryOptions{
			RootScopeID: navauthz.RootScopeID, ProtectedSubjects: []authorization.SubjectRef{admin},
			Constraints: navauthz.ConstraintEvaluators(),
			Predicates:  navauthz.PredicateEvaluators(),
		},
	)
	if err != nil {
		t.Fatalf("NewMemory() error = %v", err)
	}
	ctx := context.Background()
	for _, command := range []authorization.RegisterScopeCommand{
		{ID: navauthz.CategoryScopeID("dev"), Type: navauthz.ScopeCategory, ParentID: navauthz.RootScopeID},
		{ID: navauthz.GroupScopeID("go"), Type: navauthz.ScopeGroup, ParentID: navauthz.CategoryScopeID("dev")},
		{ID: navauthz.LinkScopeID("foundation"), Type: navauthz.ScopeLink, ParentID: navauthz.GroupScopeID("go")},
	} {
		if _, err := module.RegisterScope(ctx, command); err != nil {
			t.Fatalf("RegisterScope(%q) error = %v", command.ID, err)
		}
	}
	if _, err := module.Grant(ctx, authorization.GrantCommand{
		Actor: admin, Target: curator, Role: navauthz.RoleCurator,
		ScopeID: navauthz.CategoryScopeID("dev"), Source: authorization.GrantSourceDirect,
	}); err != nil {
		t.Fatalf("Grant(curator) error = %v", err)
	}
	resource := authorization.ResourceFacts{
		Type: "link", ID: "foundation", ScopeID: navauthz.LinkScopeID("foundation"),
		Relations: map[authorization.RelationKind][]authorization.SubjectRef{
			navauthz.RelationSubmitter: {curator},
		},
	}
	own, err := module.Decide(ctx, authorization.DecisionRequest{
		Subject: curator, Capability: navauthz.CapabilityLinkUpdate,
		ScopeID: navauthz.LinkScopeID("foundation"), Resource: resource,
	})
	if err != nil || !own.Allowed {
		t.Fatalf("Decide(own) = %#v, %v; want allow", own, err)
	}
	resource.Relations[navauthz.RelationSubmitter] = []authorization.SubjectRef{other}
	notOwn, err := module.Decide(ctx, authorization.DecisionRequest{
		Subject: curator, Capability: navauthz.CapabilityLinkUpdate,
		ScopeID: navauthz.LinkScopeID("foundation"), Resource: resource,
	})
	if err != nil || notOwn.Allowed {
		t.Fatalf("Decide(other) = %#v, %v; want deny", notOwn, err)
	}
	plan, err := module.Plan(ctx, authorization.QueryRequest{
		Subject: curator, Capability: navauthz.CapabilityLinkUpdate,
		ScopeID: navauthz.CategoryScopeID("dev"),
	})
	if err != nil || plan.Kind != authorization.QueryRelation || plan.Relation != navauthz.RelationSubmitter {
		t.Fatalf("Plan() = %#v, %v; want submitter relation", plan, err)
	}

	draft, err := module.CreatePolicyDraft(ctx, authorization.CreatePolicyDraftCommand{
		Actor: admin, ScopeID: navauthz.CategoryScopeID("dev"), ExpectedActiveRevision: 1,
	})
	if err != nil {
		t.Fatalf("CreatePolicyDraft() error = %v", err)
	}
	leadRole, err := module.CreateRole(ctx, authorization.CreateRoleCommand{
		Actor: admin, Revision: draft.Number, Key: "category-lead", DisplayName: "分类管理员",
		ScopeID: navauthz.CategoryScopeID("dev"),
		Capabilities: []authorization.CapabilityKey{
			authorization.CapabilityManage, navauthz.CapabilityLinkModerate,
		},
		Assignment: authorization.AssignmentPolicy{Sources: []authorization.GrantSource{authorization.GrantSourceDirect}},
	})
	if err != nil {
		t.Fatalf("CreateRole(category-lead) error = %v", err)
	}
	if _, err := module.ActivatePolicy(ctx, authorization.ActivatePolicyCommand{
		Actor: admin, Revision: draft.Number, ExpectedActiveRevision: 1,
	}); err != nil {
		t.Fatalf("ActivatePolicy() error = %v", err)
	}
	lead := authorization.SubjectRef{Kind: authorization.SubjectUser, ID: "lead"}
	if _, err := module.Grant(ctx, authorization.GrantCommand{
		Actor: admin, Target: lead, Role: leadRole.Key,
		ScopeID: navauthz.CategoryScopeID("dev"), Source: authorization.GrantSourceDirect,
	}); err != nil {
		t.Fatalf("Grant(lead) error = %v", err)
	}
	secondDraft, err := module.CreatePolicyDraft(ctx, authorization.CreatePolicyDraftCommand{
		Actor: lead, ScopeID: navauthz.GroupScopeID("go"), ExpectedActiveRevision: draft.Number,
	})
	if err != nil {
		t.Fatalf("delegated CreatePolicyDraft() error = %v", err)
	}
	if _, err := module.CreateRole(ctx, authorization.CreateRoleCommand{
		Actor: lead, Revision: secondDraft.Number, Key: "group-reviewer", DisplayName: "分组审核员",
		ScopeID:      navauthz.GroupScopeID("go"),
		Capabilities: []authorization.CapabilityKey{navauthz.CapabilityLinkModerate},
		Assignment:   authorization.AssignmentPolicy{Sources: []authorization.GrantSource{authorization.GrantSourceDirect}},
	}); err != nil {
		t.Fatalf("delegated CreateRole() error = %v", err)
	}
}
