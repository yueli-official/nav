package main

import (
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/yueli-official/foundation/go/authorization"
	authorizationpostgres "github.com/yueli-official/foundation/go/authorization/postgres"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	"platform/gokit/authsetup"
	"platform/gokit/observability"
	"platform/gokit/openapiexport"
	"platform/gokit/postgresdb"
	"platform/products/nav/api/internal/appconfig"
	"platform/products/nav/api/internal/catalog"
	"platform/products/nav/api/internal/dao"
	"platform/products/nav/api/internal/navaudit"
	"platform/products/nav/api/internal/navauthz"
	"platform/products/nav/api/internal/navprofile"
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
		openAPICatalog := catalog.New(nil, siteMeta("月离导航"))
		openAPICatalog.SetSiteProfile(navprofile.NewMemory())
		server.Configure(httpServer, server.Deps{Catalog: openAPICatalog})
		if handled, exportErr := openapiexport.ExportIfRequested(httpServer); handled {
			if exportErr != nil {
				panic(exportErr)
			}
			return
		}
	}

	service := catalog.New(dao.NewPG(g.DB()), siteMeta(appconfig.SiteBrand(ctx)))
	profileDB, err := postgresdb.OpenDefault(ctx)
	if err != nil {
		panic(err)
	}
	defer profileDB.Close()
	auditJournal, err := navaudit.New(ctx, profileDB, appconfig.SiteSlug(ctx))
	if err != nil {
		panic(err)
	}
	service.SetAudit(auditJournal)
	profiles, err := navprofile.NewPostgres(profileDB)
	if err != nil {
		panic(err)
	}
	service.SetSiteProfile(profiles)
	if err := service.EnsureSiteProfile(ctx); err != nil {
		panic(err)
	}
	definition, err := authorization.Compile(navauthz.Definition())
	if err != nil {
		panic(err)
	}
	bootstrapSubs := appconfig.BootstrapAdministratorSubs(ctx)
	protected := make([]authorization.SubjectRef, 0, len(bootstrapSubs))
	for _, sub := range bootstrapSubs {
		if sub != "" {
			protected = append(protected, authorization.SubjectRef{Kind: authorization.SubjectUser, ID: sub})
		}
	}
	authz, err := authorizationpostgres.New(ctx, definition, authorizationpostgres.Options{
		DB: profileDB, InstanceKey: "nav:" + appconfig.SiteSlug(ctx),
		Memory: authorization.MemoryOptions{
			RootScopeID: navauthz.RootScopeID, ProtectedSubjects: protected,
			Constraints: navauthz.ConstraintEvaluators(),
			Predicates:  navauthz.PredicateEvaluators(),
		},
	})
	if err != nil {
		panic(err)
	}
	if authz.InstanceWasCreated() && len(protected) == 0 {
		panic("nav authorization bootstrap requires at least one administrator subject")
	}
	if err := navauthz.SyncResourceScopes(ctx, profileDB, authz); err != nil {
		panic(err)
	}
	authorizationService := navauthz.New(authz, profileDB)
	jwks := appconfig.LoadJWKS(ctx)
	verifier, err := authsetup.NewRemoteVerifier(authsetup.RemoteVerifierConfig{
		JWKSURL: jwks.URL, Issuer: jwks.Issuer, Audience: jwks.Audience,
		AllowLoopbackHTTP: jwks.AllowLoopbackHTTP,
	})
	if err != nil {
		panic(err)
	}

	server.Configure(httpServer, server.Deps{
		Verifier: verifier, Catalog: service, Authorization: authorizationService,
	})
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
