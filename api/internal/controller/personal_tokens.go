package controller

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	foundationauth "github.com/yueli-official/foundation/go/auth"
	"github.com/yueli-official/foundation/go/authorization"
	v1 "github.com/yueli-official/nav/api/api/v1"
	"github.com/yueli-official/nav/api/internal/navauthz"
	"github.com/yueli-official/nav/api/internal/naverr"
)

type personalPermission struct {
	foundationauth.PersonalPermission
	Capability authorization.CapabilityKey
}

var personalPermissions = []personalPermission{
	{foundationauth.PersonalPermission{Key: string(navauthz.CapabilityLinkSubmit), Label: "提交链接", Description: "创建链接草稿；发布或推荐仍需审核权限。"}, navauthz.CapabilityLinkSubmit},
	{foundationauth.PersonalPermission{Key: string(navauthz.CapabilityLinkUpdate), Label: "读取和编辑链接", Description: "读取可管理链接并编辑当前账号授权范围内的链接。"}, navauthz.CapabilityLinkUpdate},
	{foundationauth.PersonalPermission{Key: string(navauthz.CapabilityLinkModerate), Label: "治理链接", Description: "发布、归档、删除或批量治理链接；仍需当前账号拥有对应权限。"}, navauthz.CapabilityLinkModerate},
	{foundationauth.PersonalPermission{Key: string(navauthz.CapabilityStructureManage), Label: "管理导航结构", Description: "读取和管理分类、分组与标签；仍需当前管理员权限。"}, navauthz.CapabilityStructureManage},
	{foundationauth.PersonalPermission{Key: string(navauthz.CapabilityHealthCheckRun), Label: "管理链接健康检查", Description: "读取健康状态、运行检查并维护豁免；仍需当前管理员权限。"}, navauthz.CapabilityHealthCheckRun},
	{foundationauth.PersonalPermission{Key: string(navauthz.CapabilitySettingsManage), Label: "管理站点设置", Description: "读取和修改 Nav 公开站点设置；仍需当前管理员权限。"}, navauthz.CapabilitySettingsManage},
}

type PersonalPermissions struct {
	site    string
	service *navauthz.Service
}

func NewPersonalPermissions(site string, service *navauthz.Service) *PersonalPermissions {
	return &PersonalPermissions{site: site, service: service}
}

func (controller *PersonalPermissions) GetPersonalPermissions(ctx context.Context, req *v1.PersonalPermissionsReq) (*v1.PersonalPermissionsRes, error) {
	principal, _ := foundationauth.FromContext(ctx)
	if principal == nil || principal.SubjectKind != foundationauth.SubjectClient || principal.ClientID != "identity-svc" ||
		!principal.HasScope(foundationauth.PersonalPermissionsScope) || controller.site == "" {
		return nil, naverr.Forbidden()
	}
	if controller.service == nil {
		return nil, naverr.AuthorizationUnavailable()
	}
	userCtx := foundationauth.NewContext(ctx, &foundationauth.Principal{Subject: req.UserKey, SubjectKind: foundationauth.SubjectUser})
	capabilities, err := controller.service.EffectiveManagementAccess(userCtx)
	if err != nil {
		return nil, mapAuthorizationError(err)
	}
	allowed := make(map[authorization.CapabilityKey]struct{}, len(capabilities))
	for _, capability := range capabilities {
		allowed[capability] = struct{}{}
	}
	items := make([]foundationauth.PersonalPermission, 0, len(personalPermissions))
	for _, permission := range personalPermissions {
		if _, ok := allowed[permission.Capability]; ok {
			items = append(items, permission.PersonalPermission)
		}
	}
	return &v1.PersonalPermissionsRes{Site: controller.site, UserKey: req.UserKey, Items: items}, nil
}

type personalRoute struct {
	method     string
	path       string
	capability authorization.CapabilityKey
}

var personalRoutes = []personalRoute{
	{"GET", "/api/v1/nav/catalog", ""},
	{"GET", "/api/v1/nav/groups/{groupId}", ""},
	{"POST", "/api/v1/nav/links/{id}/click", ""},
	{"GET", "/api/v1/nav/links/{id}/favicon", ""},
	{"GET", "/api/v1/admin/nav/links", navauthz.CapabilityLinkUpdate},
	{"POST", "/api/v1/admin/nav/links", navauthz.CapabilityLinkSubmit},
	{"PATCH", "/api/v1/admin/nav/links/{id}", navauthz.CapabilityLinkUpdate},
	{"DELETE", "/api/v1/admin/nav/links/{id}", navauthz.CapabilityLinkModerate},
	{"POST", "/api/v1/admin/nav/links/bulk", navauthz.CapabilityLinkModerate},
	{"GET", "/api/v1/admin/nav/checks", navauthz.CapabilityHealthCheckRun},
	{"POST", "/api/v1/admin/nav/checks/run", navauthz.CapabilityHealthCheckRun},
	{"GET", "/api/v1/admin/nav/checks/jobs/{jobId}", navauthz.CapabilityHealthCheckRun},
	{"PUT", "/api/v1/admin/nav/checks/{id}/exemption", navauthz.CapabilityHealthCheckRun},
	{"GET", "/api/v1/admin/nav/structure", navauthz.CapabilityStructureManage},
	{"POST", "/api/v1/admin/nav/categories", navauthz.CapabilityStructureManage},
	{"PATCH", "/api/v1/admin/nav/categories/{id}", navauthz.CapabilityStructureManage},
	{"DELETE", "/api/v1/admin/nav/categories/{id}", navauthz.CapabilityStructureManage},
	{"POST", "/api/v1/admin/nav/groups", navauthz.CapabilityStructureManage},
	{"PATCH", "/api/v1/admin/nav/groups/{id}", navauthz.CapabilityStructureManage},
	{"DELETE", "/api/v1/admin/nav/groups/{id}", navauthz.CapabilityStructureManage},
	{"GET", "/api/v1/admin/nav/tags", navauthz.CapabilityStructureManage},
	{"POST", "/api/v1/admin/nav/tags/rename", navauthz.CapabilityStructureManage},
	{"POST", "/api/v1/admin/nav/tags/delete", navauthz.CapabilityStructureManage},
	{"GET", "/api/v1/admin/nav/settings", navauthz.CapabilitySettingsManage},
	{"PUT", "/api/v1/admin/nav/settings", navauthz.CapabilitySettingsManage},
}

func PersonalTokenRoutes(request *ghttp.Request) {
	principal, _ := foundationauth.FromContext(request.Context())
	if principal != nil && principal.IsPersonalToken() && !allowsPersonalRoute(request.Context(), request.Method, request.URL.Path) {
		request.SetError(naverr.Forbidden())
		return
	}
	request.Middleware.Next()
}

func allowsPersonalRoute(ctx context.Context, method, path string) bool {
	for _, route := range personalRoutes {
		if route.method != method || !matchesPersonalPath(route.path, path) {
			continue
		}
		return route.capability == "" || foundationauth.AllowsPersonalCapability(ctx, string(route.capability))
	}
	return false
}

func matchesPersonalPath(pattern, path string) bool {
	want, got := strings.Split(pattern, "/"), strings.Split(path, "/")
	if len(want) != len(got) {
		return false
	}
	for index, part := range want {
		if strings.HasPrefix(part, "{") {
			if got[index] == "" || got[index] == "." || got[index] == ".." {
				return false
			}
			continue
		}
		if part != got[index] {
			return false
		}
	}
	return true
}
