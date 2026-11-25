package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/yueli-official/foundation/go/siteprofile"
)

type GetSiteBrandingReq struct {
	g.Meta `path:"/api/v1/site/branding" method:"get" tags:"site" summary:"Get public site branding"`
}

type GetSiteBrandingRes struct {
	LogoURL    string `json:"logoUrl"`
	FaviconURL string `json:"faviconUrl"`
}

type PatchSiteBrandingReq struct {
	g.Meta `path:"/api/v1/site/branding" method:"patch" tags:"site" summary:"Update site branding"`
	siteprofile.BrandImagePatch
}

type PatchSiteBrandingRes struct {
	LogoURL    string `json:"logoUrl"`
	FaviconURL string `json:"faviconUrl"`
}
