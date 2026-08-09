package controller

import (
	"context"
	"sort"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	v1 "github.com/yueli-official/nav/api/api/v1"
	"github.com/yueli-official/nav/api/internal/naverr"
	"github.com/yueli-official/nav/api/internal/navidentity"
	"github.com/yueli-official/nav/api/internal/navmember"
)

type Me struct{}

func NewMe() *Me { return &Me{} }

func (*Me) Me(ctx context.Context, _ *v1.MeReq) (*v1.MeRes, error) {
	service := authorizationService(ctx)
	if service == nil {
		return nil, naverr.AuthorizationUnavailable()
	}
	principal, authenticated := foundationauth.FromContext(ctx)
	member, hasMembership := membershipFromContext(ctx)
	if hasMembership && member.Status == navmember.StatusSuspended {
		return &v1.MeRes{Me: &v1.MeView{
			Sub: principal.Subject, UserKey: member.UserKey, Authenticated: true,
			Capabilities: []string{}, Membership: meMembershipView(member),
		}}, nil
	}
	capabilityKeys, err := service.EffectiveManagementAccess(ctx)
	if err != nil {
		return nil, naverr.AuthorizationUnavailable()
	}
	capabilities := make([]string, len(capabilityKeys))
	for index, capability := range capabilityKeys {
		capabilities[index] = string(capability)
	}
	sort.Strings(capabilities)
	userKey, _ := navidentity.FromContext(ctx)
	sub := ""
	if authenticated && principal != nil {
		sub = principal.Subject
	}
	return &v1.MeRes{Me: &v1.MeView{
		Sub: sub, UserKey: userKey, Authenticated: authenticated && principal != nil,
		IsAdministrator: service.IsAdministrator(ctx), Capabilities: capabilities,
		Membership: meMembershipView(member),
	}}, nil
}

func meMembershipView(member navmember.Member) *v1.MeMembershipView {
	if member.UserKey == "" {
		return nil
	}
	return &v1.MeMembershipView{
		Status: string(member.Status), JoinedAt: member.JoinedAt, LastSeenAt: member.LastSeenAt,
	}
}
