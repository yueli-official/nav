package v1

import (
	"context"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
)

func TestGetFaviconRequestBindsVersionFromQuery(t *testing.T) {
	t.Parallel()

	field, ok := reflect.TypeOf(GetFaviconReq{}).FieldByName("Version")
	if !ok {
		t.Fatal("GetFaviconReq.Version is missing")
	}
	if got := field.Tag.Get("p"); got != "v" {
		t.Fatalf("GetFaviconReq.Version p tag = %q, want v", got)
	}
}

func TestAdminRunChecksRequestAcceptsFifteenIDs(t *testing.T) {
	t.Parallel()

	req := &AdminRunChecksReq{IDs: []string{
		"adobe-color",
		"ambientcg",
		"artstation",
		"bangumi",
		"bbc-learning-english",
		"behance",
		"bigjpg",
		"blender",
		"blueprintue",
		"cambridge-dictionary",
		"chatgpt",
		"codeforces",
		"codepen",
		"convertio",
		"coursera",
	}}

	if err := g.Validator().Data(req).Run(context.Background()); err != nil {
		t.Fatalf("15 ids should satisfy the request contract: %v", err)
	}
}

func TestAdminRunChecksDeclaresAcceptedJobResponse(t *testing.T) {
	t.Parallel()

	meta, ok := reflect.TypeOf(AdminRunChecksRes{}).FieldByName("Meta")
	if !ok {
		t.Fatal("AdminRunChecksRes.Meta is missing")
	}
	if got := meta.Tag.Get("status"); got != "202" {
		t.Fatalf("AdminRunChecksRes status = %q, want 202", got)
	}
	if _, ok := reflect.TypeOf(AdminRunChecksRes{}).FieldByName("Job"); !ok {
		t.Fatal("AdminRunChecksRes.Job is missing")
	}
}

func TestAdminGetCheckJobUsesDedicatedStatusRoute(t *testing.T) {
	t.Parallel()

	meta, ok := reflect.TypeOf(AdminGetCheckJobReq{}).FieldByName("Meta")
	if !ok {
		t.Fatal("AdminGetCheckJobReq.Meta is missing")
	}
	if got := meta.Tag.Get("path"); got != "/api/v1/admin/nav/checks/jobs/{jobId}" {
		t.Fatalf("AdminGetCheckJobReq path = %q", got)
	}
}

func TestAdminCheckFilterSeparatesListableAndRunnableStates(t *testing.T) {
	t.Parallel()

	if err := g.Validator().Data(&AdminListChecksReq{Health: "exempt", Page: 1, Size: 20}).Run(context.Background()); err != nil {
		t.Fatalf("exempt should satisfy the list filter contract: %v", err)
	}
	if err := g.Validator().Data(&AdminRunChecksReq{Scope: "filtered", Health: "exempt"}).Run(context.Background()); err == nil {
		t.Fatal("exempt must not satisfy the runnable health filter contract")
	}
}

func TestAdminSetCheckExemptionUsesDedicatedIdempotentRoute(t *testing.T) {
	t.Parallel()

	meta, ok := reflect.TypeOf(AdminSetCheckExemptionReq{}).FieldByName("Meta")
	if !ok {
		t.Fatal("AdminSetCheckExemptionReq.Meta is missing")
	}
	if got := meta.Tag.Get("path"); got != "/api/v1/admin/nav/checks/{id}/exemption" {
		t.Fatalf("AdminSetCheckExemptionReq path = %q", got)
	}
	if got := meta.Tag.Get("method"); got != "PUT" {
		t.Fatalf("AdminSetCheckExemptionReq method = %q, want PUT", got)
	}
	exempt, ok := reflect.TypeOf(AdminSetCheckExemptionReq{}).FieldByName("Exempt")
	if !ok {
		t.Fatal("AdminSetCheckExemptionReq.Exempt is missing")
	}
	if exempt.Type != reflect.TypeOf((*bool)(nil)) || exempt.Tag.Get("v") != "required" {
		t.Fatalf("AdminSetCheckExemptionReq.Exempt must be an explicitly required boolean pointer")
	}
}
