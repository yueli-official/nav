package server

import (
	"github.com/gogf/gf/v2/net/ghttp"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	"platform/gokit/authhttp"
	"platform/gokit/ghttpx"
	"platform/gokit/healthcheck"
	"platform/products/nav/api/internal/catalog"
	"platform/products/nav/api/internal/controller"
)

type Deps struct {
	Verifier    *foundationauth.Verifier
	Catalog     *catalog.Service
	ReadyChecks map[string]healthcheck.Check
}

func Configure(server *ghttp.Server, deps Deps) {
	readyChecks := deps.ReadyChecks
	if readyChecks == nil {
		readyChecks = map[string]healthcheck.Check{"database": healthcheck.Database}
	}
	server.Use(ghttpx.TraceRouteMiddleware)
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttpx.Middleware)
		group.GET("/healthz", controller.Healthz)
		group.GET("/readyz", healthcheck.Handler(readyChecks))
	})
	if deps.Catalog == nil {
		return
	}
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttpx.Middleware)
		group.Bind(controller.NewPublic(deps.Catalog))
	})
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttpx.Middleware, authhttp.Required(deps.Verifier))
		group.Bind(controller.NewAdmin(deps.Catalog))
	})
}
