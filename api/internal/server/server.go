package server

import (
	"github.com/gogf/gf/v2/net/ghttp"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	goframeauth "github.com/yueli-official/foundation/go/goframe/auth"
	"github.com/yueli-official/nav/api/internal/catalog"
	"github.com/yueli-official/nav/api/internal/controller"
	"github.com/yueli-official/nav/api/internal/navauthz"
	"github.com/yueli-official/nav/api/internal/navmember"
	"github.com/yueli-official/nav/api/internal/runtime"
)

type Deps struct {
	Verifier         *foundationauth.Verifier
	PersonalVerifier *foundationauth.PersonalTokenVerifier
	PersonalSite     string
	Catalog          *catalog.Service
	Authorization    *navauthz.Service
	Membership       navmember.Directory
	ReadyChecks      map[string]runtime.ReadinessCheck
}

func Configure(server *ghttp.Server, deps Deps) {
	apiMiddleware := runtime.MustAPIMiddleware(runtime.MustRateLimiterFromEnvironment()).Handle
	readyChecks := deps.ReadyChecks
	if readyChecks == nil {
		readyChecks = map[string]runtime.ReadinessCheck{"database": runtime.DatabaseReadiness}
	}
	server.Use(runtime.TraceRouteMiddleware)
	var verifier goframeauth.TokenVerifier
	if deps.Verifier != nil {
		verifier = foundationauth.CompositeVerifier{JWT: deps.Verifier, Personal: deps.PersonalVerifier}
	}
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(apiMiddleware, controller.CauseMappingMiddleware)
		group.GET("/healthz", controller.Healthz)
		group.GET("/readyz", runtime.ReadinessHandler(readyChecks))
	})
	server.Group("/", func(group *ghttp.RouterGroup) {
		middlewares := []ghttp.HandlerFunc{apiMiddleware, controller.CauseMappingMiddleware}
		if verifier != nil {
			middlewares = append(middlewares, runtime.OptionalAuth(verifier), controller.PersonalTokenRoutes)
		}
		middlewares = append(middlewares, controller.AuthorizationMiddleware(deps.Authorization), controller.MembershipMiddleware(deps.Membership, deps.Authorization, false))
		group.Middleware(middlewares...)
		group.Bind(controller.NewMe())
	})
	if deps.Catalog == nil {
		return
	}
	server.Group("/", func(group *ghttp.RouterGroup) {
		middlewares := []ghttp.HandlerFunc{apiMiddleware, controller.CauseMappingMiddleware}
		if verifier != nil {
			middlewares = append(middlewares, runtime.OptionalAuth(verifier), controller.PersonalTokenRoutes)
		}
		group.Middleware(middlewares...)
		group.Bind(controller.NewPublic(deps.Catalog))
	})
	server.Group("/", func(group *ghttp.RouterGroup) {
		middlewares := []ghttp.HandlerFunc{apiMiddleware, controller.CauseMappingMiddleware}
		if verifier != nil {
			middlewares = append(middlewares, runtime.RequiredAuth(verifier), controller.PersonalTokenRoutes)
		}
		middlewares = append(middlewares, controller.AuthorizationMiddleware(deps.Authorization), controller.MembershipMiddleware(deps.Membership, deps.Authorization, true))
		group.Middleware(middlewares...)
		group.Bind(controller.NewAuthorization())
		group.Bind(controller.NewAdmin(deps.Catalog))
		group.Bind(controller.NewMembers(deps.Membership))
		group.Bind(controller.NewPersonalPermissions(deps.PersonalSite, deps.Authorization))
	})
}
