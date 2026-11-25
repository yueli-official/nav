package controller

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/yueli-official/foundation/go/authorization"
	"github.com/yueli-official/foundation/go/siteprofile"
	v1 "github.com/yueli-official/nav/api/api/v1"
	"github.com/yueli-official/nav/api/internal/navauthz"
)

var navBrandingDefaults = siteprofile.BrandImageURLs{LogoURL: "/favicon.svg", FaviconURL: "/favicon.svg"}

func (c *Public) GetSiteBranding(ctx context.Context, _ *v1.GetSiteBrandingReq) (*v1.GetSiteBrandingRes, error) {
	site, err := c.service.PublicSite(ctx)
	if err != nil {
		return nil, err
	}
	setPublicProfileHeaders(ctx, site.ETag)
	g.RequestFromCtx(ctx).Response.Header().Set("Access-Control-Allow-Origin", "*")
	urls := site.BrandImageURLs.WithDefaults(navBrandingDefaults)
	return &v1.GetSiteBrandingRes{LogoURL: urls.LogoURL, FaviconURL: urls.FaviconURL}, nil
}

func (c *Admin) PatchSiteBranding(ctx context.Context, req *v1.PatchSiteBrandingReq) (*v1.PatchSiteBrandingRes, error) {
	if err := requireCapability(ctx, navauthz.CapabilitySettingsManage, navauthz.RootScopeID, authorization.ResourceFacts{}); err != nil {
		return nil, err
	}
	current, err := c.service.AdminSiteSettings(ctx)
	if err != nil {
		return nil, err
	}
	profile := current.Snapshot.Profile
	if req.BrandImagePatch.Apply(&profile.Branding, navBrandingDefaults) {
		current, err = c.service.SaveAdminSiteSettings(ctx, current.Snapshot.Revision, current.RuntimeRevision, profile, current.SearchPlaceholder)
		if err != nil {
			return nil, err
		}
	}
	setAdminProfileHeaders(ctx, current.ETag)
	urls := current.Snapshot.Profile.Branding.BrandImageURLs.WithDefaults(navBrandingDefaults)
	return &v1.PatchSiteBrandingRes{LogoURL: urls.LogoURL, FaviconURL: urls.FaviconURL}, nil
}
