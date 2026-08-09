package controller

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	foundationauth "github.com/yueli-official/foundation/go/auth"
	"github.com/yueli-official/foundation/go/authorization"

	v1 "github.com/yueli-official/nav/api/api/v1"
	"github.com/yueli-official/nav/api/internal/navauthz"
	"github.com/yueli-official/nav/api/internal/naverr"
	"github.com/yueli-official/nav/api/internal/navidentity"
	"github.com/yueli-official/nav/api/internal/navmember"
)

type membershipContextKey struct{}

func MembershipMiddleware(directory navmember.Directory, authorizationService *navauthz.Service, requireActive bool) ghttp.HandlerFunc {
	return func(request *ghttp.Request) {
		principal, authenticated := foundationauth.FromContext(request.Context())
		kind, _ := principal.Claim("subject_kind")
		if !authenticated || principal == nil || kind != "user" {
			request.Middleware.Next()
			return
		}
		userKey, ok := navidentity.PublicUserKey(principal)
		if !ok || directory == nil {
			request.SetError(naverr.MembershipUnavailable())
			return
		}
		result, err := directory.Ensure(request.Context(), userKey)
		if err != nil {
			request.SetError(naverr.MembershipUnavailable())
			return
		}
		if result.NeedsJoinReconcile {
			if authorizationService == nil || authorizationService.ReconcileNewMember(request.Context()) != nil ||
				directory.MarkJoinReconciled(request.Context(), userKey) != nil {
				request.SetError(naverr.MembershipUnavailable())
				return
			}
		}
		if requireActive && result.Member.Status == navmember.StatusSuspended {
			request.SetError(naverr.MembershipSuspended())
			return
		}
		request.SetCtx(context.WithValue(request.Context(), membershipContextKey{}, result.Member))
		request.Middleware.Next()
	}
}

func membershipFromContext(ctx context.Context) (navmember.Member, bool) {
	member, ok := ctx.Value(membershipContextKey{}).(navmember.Member)
	return member, ok
}

type Members struct {
	directory navmember.Directory
}

func NewMembers(directory navmember.Directory) *Members {
	return &Members{directory: directory}
}

func (controller *Members) ListMembers(ctx context.Context, req *v1.AdminListMembersReq) (*v1.AdminListMembersRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	roles, grants, applications, err := memberAuthorizationState(ctx)
	if err != nil {
		return nil, err
	}
	query := navmember.Query{
		Search: req.Q, Status: navmember.Status(req.Status), Page: req.Page, Size: req.Size,
	}
	if req.Role != "" {
		known := false
		for _, role := range roles {
			known = known || string(role.Key) == req.Role
		}
		if !known {
			return nil, naverr.Validation("role", "unknown", map[string]any{"role": req.Role})
		}
		query.ConstrainUserKeys = true
		for _, grant := range grants {
			if string(grant.Role) == req.Role && grant.Target.Kind == authorization.SubjectUser {
				query.UserKeys = append(query.UserKeys, grant.Target.ID)
			}
		}
	}
	page, err := controller.directory.List(ctx, query)
	if err != nil {
		return nil, mapMembershipError(err)
	}
	roleNames := make(map[string]string, len(roles))
	for _, role := range roles {
		roleNames[string(role.Key)] = role.DisplayName
	}
	memberRoles := make(map[string][]v1.MemberRoleView)
	seenRoles := make(map[string]map[string]struct{})
	for _, grant := range grants {
		if grant.Target.Kind != authorization.SubjectUser {
			continue
		}
		key := string(grant.Role) + "\x00" + string(grant.Source)
		if seenRoles[grant.Target.ID] == nil {
			seenRoles[grant.Target.ID] = map[string]struct{}{}
		}
		if _, exists := seenRoles[grant.Target.ID][key]; exists {
			continue
		}
		seenRoles[grant.Target.ID][key] = struct{}{}
		memberRoles[grant.Target.ID] = append(memberRoles[grant.Target.ID], v1.MemberRoleView{
			Key: string(grant.Role), DisplayName: roleNames[string(grant.Role)], Source: string(grant.Source),
		})
	}
	pending := make(map[string]int)
	for _, application := range applications {
		if application.Subject.Kind == authorization.SubjectUser {
			pending[application.Subject.ID]++
		}
	}
	views := make([]v1.MemberView, len(page.Members))
	for index, member := range page.Members {
		views[index] = memberView(member, memberRoles[member.UserKey], pending[member.UserKey])
	}
	return &v1.AdminListMembersRes{
		Members: views,
		Counts:  v1.MembershipStatusCountsView{All: page.Counts.All, Active: page.Counts.Active, Suspended: page.Counts.Suspended},
		Roles:   authorizationRoleViews(roles), Total: page.Total, Page: page.Page, Size: page.Size,
	}, nil
}

func (controller *Members) SetMemberStatus(ctx context.Context, req *v1.AdminSetMemberStatusReq) (*v1.AdminSetMemberStatusRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	actor, ok := navidentity.FromContext(ctx)
	if !ok {
		return nil, naverr.MembershipUnavailable()
	}
	member, err := controller.directory.SetStatus(ctx, navmember.SetStatusCommand{
		UserKey: req.UserKey, Status: navmember.Status(req.Status), ActorUserKey: actor, Reason: req.Reason,
	})
	if err != nil {
		return nil, mapMembershipError(err)
	}
	roles, grants, applications, err := memberAuthorizationState(ctx)
	if err != nil {
		return nil, err
	}
	roleNames := make(map[string]string, len(roles))
	for _, role := range roles {
		roleNames[string(role.Key)] = role.DisplayName
	}
	var roleViews []v1.MemberRoleView
	seen := map[string]struct{}{}
	for _, grant := range grants {
		if grant.Target.Kind != authorization.SubjectUser || grant.Target.ID != member.UserKey {
			continue
		}
		key := string(grant.Role) + "\x00" + string(grant.Source)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		roleViews = append(roleViews, v1.MemberRoleView{
			Key: string(grant.Role), DisplayName: roleNames[string(grant.Role)], Source: string(grant.Source),
		})
	}
	pending := 0
	for _, application := range applications {
		if application.Subject.Kind == authorization.SubjectUser && application.Subject.ID == member.UserKey {
			pending++
		}
	}
	return &v1.AdminSetMemberStatusRes{Member: memberView(member, roleViews, pending)}, nil
}

func memberAuthorizationState(ctx context.Context) ([]authorization.Role, []authorization.Grant, []authorization.Application, error) {
	service := authorizationService(ctx)
	if service == nil || service.Runtime() == nil {
		return nil, nil, nil, naverr.AuthorizationUnavailable()
	}
	actor := service.Subject(ctx)
	var roles []authorization.Role
	for offset := 0; ; {
		page, err := service.Runtime().ListRoles(ctx, authorization.RoleListQuery{
			Actor: actor, ScopeID: navauthz.RootScopeID, Offset: offset, Limit: 500,
		})
		if err != nil {
			return nil, nil, nil, mapAuthorizationError(err)
		}
		roles = append(roles, page.Roles...)
		offset += len(page.Roles)
		if offset >= page.Total || len(page.Roles) == 0 {
			break
		}
	}
	var grants []authorization.Grant
	for offset := 0; ; {
		page, err := service.Runtime().ListGrants(ctx, authorization.GrantListQuery{
			Actor: actor, ScopeID: navauthz.RootScopeID, ActiveOnly: true, Offset: offset, Limit: 500,
		})
		if err != nil {
			return nil, nil, nil, mapAuthorizationError(err)
		}
		grants = append(grants, page.Grants...)
		offset += len(page.Grants)
		if offset >= page.Total || len(page.Grants) == 0 {
			break
		}
	}
	var applications []authorization.Application
	for offset := 0; ; {
		page, err := service.Runtime().ListApplications(ctx, authorization.ApplicationListQuery{
			Actor: actor, ScopeID: navauthz.RootScopeID, State: authorization.ApplicationPending,
			Offset: offset, Limit: 500,
		})
		if err != nil {
			return nil, nil, nil, mapAuthorizationError(err)
		}
		applications = append(applications, page.Applications...)
		offset += len(page.Applications)
		if offset >= page.Total || len(page.Applications) == 0 {
			break
		}
	}
	return roles, grants, applications, nil
}

func memberView(member navmember.Member, roles []v1.MemberRoleView, pending int) v1.MemberView {
	if roles == nil {
		roles = []v1.MemberRoleView{}
	}
	sort.Slice(roles, func(i, j int) bool {
		if roles[i].DisplayName == roles[j].DisplayName {
			return roles[i].Source < roles[j].Source
		}
		return roles[i].DisplayName < roles[j].DisplayName
	})
	return v1.MemberView{
		UserKey: member.UserKey, Status: string(member.Status), DisplayName: member.DisplayName,
		Handle: member.Handle, AvatarMediaKey: member.AvatarMediaKey,
		JoinedAt: member.JoinedAt, LastSeenAt: member.LastSeenAt,
		SuspendedAt: member.SuspendedAt, SuspendedBy: member.SuspendedBy,
		SuspensionReason: member.SuspensionReason, SubmissionCount: member.SubmissionCount,
		PendingApplications: pending, Roles: roles,
	}
}

func mapMembershipError(err error) error {
	switch {
	case errors.Is(err, navmember.ErrNotFound):
		return naverr.NotFound("membership")
	case errors.Is(err, navmember.ErrInvalidUserKey), errors.Is(err, navmember.ErrInvalidStatus):
		return naverr.Validation("membership", "invalid", nil)
	case errors.Is(err, navmember.ErrSelfSuspend):
		return naverr.Conflict("cannot-suspend-current-administrator")
	case errors.Is(err, navmember.ErrReasonRequired):
		return naverr.Validation("reason", "required", nil)
	case strings.Contains(err.Error(), "500 characters"):
		return naverr.Validation("reason", "length", map[string]any{"max": 500})
	default:
		return naverr.MembershipUnavailable()
	}
}
