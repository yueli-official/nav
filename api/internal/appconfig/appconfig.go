package appconfig

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type JWKS struct {
	URL               string
	Issuer            string
	Audience          string
	AllowLoopbackHTTP bool
}

func LoadJWKS(ctx context.Context) JWKS {
	return JWKS{
		URL:               g.Cfg().MustGet(ctx, "nav.jwks.url", "http://localhost:8081/oauth2/jwks.json").String(),
		Issuer:            g.Cfg().MustGet(ctx, "nav.jwks.issuer", "http://localhost:8081").String(),
		Audience:          g.Cfg().MustGet(ctx, "nav.jwks.audience", "").String(),
		AllowLoopbackHTTP: g.Cfg().MustGet(ctx, "nav.jwks.allowLoopbackHttp", false).Bool(),
	}
}

func SiteBrand(ctx context.Context) string {
	return g.Cfg().MustGet(ctx, "nav.brand", "月离导航").String()
}
