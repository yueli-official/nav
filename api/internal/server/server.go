package server

import (
	"github.com/gogf/gf/v2/net/ghttp"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	"github.com/yueli-official/nav/api/internal/catalog"
	"github.com/yueli-official/nav/api/internal/controller"
	"github.com/yueli-official/nav/api/internal/navauthz"
	"github.com/yueli-official/nav/api/internal/navmember"
	"github.com/yueli-official/nav/api/internal/runtime"
)

type Deps struct {
	Verifier      *foundationauth.Verifier
	Catalog       *catalog.Service
	Authorization *navauthz.Service
	Membership    navmember.Directory
	ReadyChecks   map[string]runtime.ReadinessCheck
}

func Configure(server *ghttp.Server, deps Deps) {
	apiMiddleware := runtime.MustAPIMiddleware(runtime.MustRateLimiterFromEnvironment()).Handle
	readyChecks := deps.ReadyChecks
	if readyChecks == nil {
		readyChecks = map[string]runtime.ReadinessCheck{"database": runtime.DatabaseReadiness}
	}
	server.Use(runtime.TraceRouteMiddleware)
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(apiMiddleware)
		group.GET("/healthz", controller.Healthz)
		group.GET("/readyz", runtime.ReadinessHandler(readyChecks))
	})
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			apiMiddleware, runtime.OptionalAuth(deps.Verifier), controller.AuthorizationMiddleware(deps.Authorization),
			controller.MembershipMiddleware(deps.Membership, deps.Authorization, false),
		)
		group.Bind(controller.NewMe())
	})
	if deps.Catalog == nil {
		return
	}
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(apiMiddleware)
		group.Bind(controller.NewPublic(deps.Catalog))
	})
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			apiMiddleware, runtime.RequiredAuth(deps.Verifier), controller.AuthorizationMiddleware(deps.Authorization),
			controller.MembershipMiddleware(deps.Membership, deps.Authorization, true),
		)
		group.Bind(controller.NewAuthorization())
		group.Bind(controller.NewAdmin(deps.Catalog))
		group.Bind(controller.NewMembers(deps.Membership))
	})
}
