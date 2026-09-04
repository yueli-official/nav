package controller

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/yueli-official/nav/api/internal/catalog"
	"github.com/yueli-official/nav/api/internal/model"
	"github.com/yueli-official/nav/api/internal/runtime"
	"io"
	"net/http"
	"testing"
	"time"
)

type faviconStore struct{ catalog.Store }

func (faviconStore) LinkByID(context.Context, string) (*model.Link, error) {
	return &model.Link{ID: "cached", URL: "https://example.com", Status: "published"}, nil
}
func (faviconStore) FaviconByLinkID(context.Context, string) (*model.FaviconCache, error) {
	return &model.FaviconCache{SourceURL: "https://example.com", Content: []byte("png-fixture"), ContentType: "image/png", ContentHash: "fixture", RefreshAfter: gtime.New(time.Now().Add(time.Hour))}, nil
}
func (faviconStore) UpsertFavicon(context.Context, *model.FaviconCache) error { return nil }

func TestFaviconConditionalResponseHasNoJSONBody(t *testing.T) {
	server := g.Server(t.Name())
	server.SetAddr("127.0.0.1:0")
	server.SetDumpRouterMap(false)
	middleware := runtime.MustAPIMiddleware(nil)
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Handle, CauseMappingMiddleware)
		group.Bind(NewPublic(catalog.New(faviconStore{}, catalog.Site{})))
	})
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer server.Shutdown()
	url := fmt.Sprintf("http://127.0.0.1:%d/api/v1/nav/links/cached/favicon", server.GetListenedPort())
	response, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(response.Body)
	response.Body.Close()
	if response.StatusCode != 200 || string(body) != "png-fixture" {
		t.Fatalf("status=%d body=%q", response.StatusCode, body)
	}
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("If-None-Match", response.Header.Get("ETag"))
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, _ = io.ReadAll(response.Body)
	if response.StatusCode != 304 || len(body) != 0 {
		t.Fatalf("status=%d body=%q", response.StatusCode, body)
	}
}
