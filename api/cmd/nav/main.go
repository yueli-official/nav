package main

import (
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	"platform/gokit/authsetup"
	"platform/gokit/observability"
	"platform/gokit/openapiexport"
	"platform/products/nav/api/internal/appconfig"
	"platform/products/nav/api/internal/catalog"
	"platform/products/nav/api/internal/dao"
	"platform/products/nav/api/internal/server"
)

func main() {
	ctx := gctx.New()
	shutdown, err := observability.StartFromEnvironment(ctx, "nav-api")
	if err != nil {
		panic(err)
	}
	defer observability.ShutdownWithTimeout(shutdown)

	httpServer := g.Server()
	if os.Getenv("PLATFORM_OPENAPI_OUTPUT") != "" {
		server.Configure(httpServer, server.Deps{Catalog: catalog.New(nil, siteMeta("月离导航"))})
		if handled, exportErr := openapiexport.ExportIfRequested(httpServer); handled {
			if exportErr != nil {
				panic(exportErr)
			}
			return
		}
	}

	service := catalog.New(dao.NewPG(g.DB()), siteMeta(appconfig.SiteBrand(ctx)))
	jwks := appconfig.LoadJWKS(ctx)
	verifier, err := authsetup.NewRemoteVerifier(authsetup.RemoteVerifierConfig{
		JWKSURL: jwks.URL, Issuer: jwks.Issuer, Audience: jwks.Audience,
		AllowLoopbackHTTP: jwks.AllowLoopbackHTTP,
	})
	if err != nil {
		panic(err)
	}

	server.Configure(httpServer, server.Deps{Verifier: verifier, Catalog: service})
	g.Log().Info(ctx, "nav service starting")
	httpServer.Run()
}

func siteMeta(brand string) catalog.Site {
	return catalog.Site{
		Name:              brand,
		Title:             "把常用互联网，整理成工作台",
		Description:       "为创作者与开发者整理的精选互联网入口，按任务浏览，也可以直接搜索名称、标签和域名。",
		SearchPlaceholder: "搜索工具、文档、社区或关键词",
		FooterTagline:     "月离导航，持续整理值得回访的互联网入口。",
	}
}
