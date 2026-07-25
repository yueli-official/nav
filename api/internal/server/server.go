package server

import (
	"github.com/gogf/gf/v2/net/ghttp"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	"platform/gokit/authhttp"
	"platform/gokit/ghttpx"
	"platform/gokit/healthcheck"
	"platform/products/nav/api/internal/catalog"
	"platform/products/nav/api/internal/controller"
	"platform/products/nav/api/internal/navauthz"
)

type Deps struct {
	Verifier      *foundationauth.Verifier
	Catalog       *catalog.Service
	Authorization *navauthz.Service
	ReadyChecks   map[string]healthcheck.Check
}

func Configure(server *ghttp.Server, deps Deps) {
	apiMiddleware := ghttpx.NewMiddleware(ghttpx.MustRateLimiterFromEnvironment(), ghttpx.ForwardedClientIPKey)
	readyChecks := deps.ReadyChecks
	if readyChecks == nil {
		readyChecks = map[string]healthcheck.Check{"database": healthcheck.Database}
	}
	server.Use(ghttpx.TraceRouteMiddleware)
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(apiMiddleware)
		group.GET("/healthz", controller.Healthz)
		group.GET("/readyz", healthcheck.Handler(readyChecks))
	})
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(apiMiddleware, authhttp.Optional(deps.Verifier), controller.AuthorizationMiddleware(deps.Authorization))
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
		group.Middleware(apiMiddleware, authhttp.Required(deps.Verifier), controller.AuthorizationMiddleware(deps.Authorization))
		group.Bind(controller.NewAuthorization())
		group.Bind(controller.NewAdmin(deps.Catalog))
	})
}
