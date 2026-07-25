package controller

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/yueli-official/foundation/go/authorization"

	"platform/products/nav/api/internal/navauthz"
	"platform/products/nav/api/internal/naverr"
)

type authorizationContextKey struct{}

func AuthorizationMiddleware(service *navauthz.Service) ghttp.HandlerFunc {
	return func(request *ghttp.Request) {
		ctx := context.WithValue(request.Context(), authorizationContextKey{}, service)
		correlationID := strings.TrimSpace(request.Header.Get("X-Trace-Id"))
		if correlationID == "" {
			correlationID = strings.TrimSpace(request.Header.Get("X-Request-Id"))
		}
		ctx = authorization.WithRequestMetadata(ctx, authorization.RequestMetadata{CorrelationID: correlationID})
		request.SetCtx(ctx)
		request.Middleware.Next()
	}
}

func authorizationService(ctx context.Context) *navauthz.Service {
	service, _ := ctx.Value(authorizationContextKey{}).(*navauthz.Service)
	return service
}

func requireCapability(
	ctx context.Context,
	capability authorization.CapabilityKey,
	scopeID authorization.ScopeID,
	resource authorization.ResourceFacts,
) error {
	service := authorizationService(ctx)
	if service == nil {
		return naverr.AuthorizationUnavailable()
	}
	decision, err := service.Decide(ctx, capability, scopeID, resource)
	if err != nil {
		if authorization.Is(err, authorization.ErrorUnavailable) {
			return naverr.AuthorizationUnavailable()
		}
		return naverr.Forbidden()
	}
	if !decision.Allowed {
		return naverr.Forbidden()
	}
	return nil
}

func requireAdmin(ctx context.Context) error {
	return requireCapability(ctx, authorization.CapabilityManage, navauthz.RootScopeID, authorization.ResourceFacts{})
}

func ensureCategoryScope(ctx context.Context, categoryID string) error {
	service := authorizationService(ctx)
	if service == nil || service.EnsureCategoryScope(ctx, categoryID) != nil {
		return naverr.AuthorizationUnavailable()
	}
	return nil
}

func ensureGroupScope(ctx context.Context, groupID, categoryID string) error {
	service := authorizationService(ctx)
	if service == nil || service.EnsureGroupScope(ctx, groupID, categoryID) != nil {
		return naverr.AuthorizationUnavailable()
	}
	return nil
}

func ensureLinkScope(ctx context.Context, linkID, groupID, categoryID string) error {
	service := authorizationService(ctx)
	if service == nil || service.EnsureLinkScope(ctx, linkID, groupID, categoryID) != nil {
		return naverr.AuthorizationUnavailable()
	}
	return nil
}

func reparentGroupScope(ctx context.Context, groupID, categoryID string) error {
	service := authorizationService(ctx)
	if service == nil || service.ReparentGroupScope(ctx, groupID, categoryID) != nil {
		return naverr.AuthorizationUnavailable()
	}
	return nil
}

func reparentLinkScope(ctx context.Context, linkID, groupID, categoryID string) error {
	service := authorizationService(ctx)
	if service == nil || service.ReparentLinkScope(ctx, linkID, groupID, categoryID) != nil {
		return naverr.AuthorizationUnavailable()
	}
	return nil
}
