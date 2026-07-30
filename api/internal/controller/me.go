package controller

import (
	"context"
	"sort"

	v1 "platform/products/nav/api/api/v1"
	"platform/products/nav/api/internal/naverr"
)

type Me struct{}

func NewMe() *Me { return &Me{} }

func (*Me) Me(ctx context.Context, _ *v1.MeReq) (*v1.MeRes, error) {
	service := authorizationService(ctx)
	if service == nil {
		return nil, naverr.AuthorizationUnavailable()
	}
	subject := service.Subject(ctx)
	capabilityKeys, err := service.EffectiveManagementAccess(ctx)
	if err != nil {
		return nil, naverr.AuthorizationUnavailable()
	}
	capabilities := make([]string, len(capabilityKeys))
	for index, capability := range capabilityKeys {
		capabilities[index] = string(capability)
	}
	sort.Strings(capabilities)
	return &v1.MeRes{Me: &v1.MeView{
		Sub: string(subject.ID), Authenticated: subject.ID != "",
		IsAdministrator: service.IsAdministrator(ctx), Capabilities: capabilities,
	}}, nil
}
