package main

import (
	"net/http"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/yueli-official/foundation/go/authorization"
	authorizationpostgres "github.com/yueli-official/foundation/go/authorization/postgres"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	"github.com/yueli-official/nav/api/internal/appconfig"
	"github.com/yueli-official/nav/api/internal/catalog"
	"github.com/yueli-official/nav/api/internal/dao"
	"github.com/yueli-official/nav/api/internal/identityclient"
	"github.com/yueli-official/nav/api/internal/navaudit"
	"github.com/yueli-official/nav/api/internal/navauthz"
	"github.com/yueli-official/nav/api/internal/navmember"
	"github.com/yueli-official/nav/api/internal/navprofile"
	"github.com/yueli-official/nav/api/internal/runtime"
	"github.com/yueli-official/nav/api/internal/server"
)

func main() {
	if err := runtime.EnableEnvironmentConfig(); err != nil {
		panic(err)
	}
	ctx := gctx.New()
	shutdown, err := runtime.StartTelemetry(ctx, "nav-api")
	if err != nil {
		panic(err)
	}
	defer runtime.ShutdownTelemetry(shutdown)

	httpServer := g.Server()
	if runtime.OpenAPIRequested() {
		openAPICatalog := catalog.New(nil, siteMeta("月离导航"))
		openAPICatalog.SetSiteProfile(navprofile.NewMemory())
		server.Configure(httpServer, server.Deps{Catalog: openAPICatalog})
		if handled, exportErr := runtime.ExportOpenAPIIfRequested(httpServer); handled {
			if exportErr != nil {
				panic(exportErr)
			}
			return
		}
	}

	service := catalog.New(dao.NewPG(g.DB()), siteMeta(appconfig.SiteBrand(ctx)))
	profileDB, err := runtime.OpenDefaultPostgres(ctx)
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
	membershipService := navmember.New(profileDB, identityclient.NewHTTP(
		appconfig.IdentityBaseURL(ctx),
		runtime.TelemetryHTTPClient(&http.Client{Timeout: 2 * time.Second}),
	))
	jwks := appconfig.LoadJWKS(ctx)
	verifier, err := runtime.NewRemoteVerifier(runtime.RemoteVerifierConfig{
		JWKSURL: jwks.URL, Issuer: jwks.Issuer, Audience: jwks.Audience,
		AllowLoopbackHTTP: jwks.AllowLoopbackHTTP,
	})
	if err != nil {
		panic(err)
	}

	server.Configure(httpServer, server.Deps{
		Verifier: verifier, Catalog: service, Authorization: authorizationService, Membership: membershipService,
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
