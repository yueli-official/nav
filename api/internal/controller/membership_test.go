package controller

import (
	"context"
	"testing"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	"github.com/yueli-official/foundation/go/authorization"
	v1 "github.com/yueli-official/nav/api/api/v1"
	"github.com/yueli-official/nav/api/internal/navauthz"
	"github.com/yueli-official/nav/api/internal/navmember"
	"github.com/yueli-official/nav/api/internal/testidentity"
)

type membershipDirectoryFixture struct {
	page      navmember.Page
	lastQuery navmember.Query
}

func (fixture *membershipDirectoryFixture) Ensure(context.Context, string) (navmember.EnsureResult, error) {
	return navmember.EnsureResult{}, nil
}
func (fixture *membershipDirectoryFixture) MarkJoinReconciled(context.Context, string) error {
	return nil
}
func (fixture *membershipDirectoryFixture) Get(context.Context, string) (navmember.Member, error) {
	return navmember.Member{}, nil
}
func (fixture *membershipDirectoryFixture) List(_ context.Context, query navmember.Query) (navmember.Page, error) {
	fixture.lastQuery = query
	return fixture.page, nil
}
func (fixture *membershipDirectoryFixture) SetStatus(context.Context, navmember.SetStatusCommand) (navmember.Member, error) {
	return navmember.Member{}, nil
}

func TestListMembersComposesMembershipGrantAndApplicationWithoutConflatingThem(t *testing.T) {
	const adminKey = "TestA123"
	const memberKey = "TestB234"
	admin := authorization.SubjectRef{Kind: authorization.SubjectUser, ID: adminKey}
	module, err := authorization.NewMemory(
		authorization.MustCompile(navauthz.Definition()),
		authorization.MemoryOptions{
			RootScopeID: navauthz.RootScopeID, ProtectedSubjects: []authorization.SubjectRef{admin},
			Constraints: navauthz.ConstraintEvaluators(), Predicates: navauthz.PredicateEvaluators(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	memberSubject := authorization.SubjectRef{Kind: authorization.SubjectUser, ID: memberKey}
	if _, err := module.Apply(context.Background(), authorization.ApplyCommand{
		Actor: memberSubject, Role: navauthz.RoleCurator, ScopeID: navauthz.RootScopeID,
		Reason: "wants a scoped assignment",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := module.Grant(context.Background(), authorization.GrantCommand{
		Actor: admin, Target: memberSubject, Role: navauthz.RoleCurator,
		ScopeID: navauthz.RootScopeID, Source: authorization.GrantSourceDirect,
	}); err != nil {
		t.Fatal(err)
	}
	directory := &membershipDirectoryFixture{page: navmember.Page{
		Members: []navmember.Member{{UserKey: memberKey, Status: navmember.StatusActive}},
		Counts:  navmember.Counts{All: 1, Active: 1}, Total: 1, Page: 1, Size: 20,
	}}
	service := navauthz.New(module, nil)
	ctx := foundationauth.NewContext(context.Background(), testidentity.User(t, adminKey, nil, nil))
	ctx = context.WithValue(ctx, authorizationContextKey{}, service)

	response, err := NewMembers(directory).ListMembers(ctx, &v1.AdminListMembersReq{
		Role: string(navauthz.RoleCurator), Page: 1, Size: 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !directory.lastQuery.ConstrainUserKeys || len(directory.lastQuery.UserKeys) != 1 || directory.lastQuery.UserKeys[0] != memberKey {
		t.Fatalf("membership query = %#v", directory.lastQuery)
	}
	if len(response.Members) != 1 || len(response.Members[0].Roles) != 1 ||
		response.Members[0].Roles[0].Key != string(navauthz.RoleCurator) ||
		response.Members[0].PendingApplications != 1 {
		t.Fatalf("member response = %#v", response.Members)
	}
}
