package catalog

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/yueli-official/foundation/go/problem"
	"github.com/yueli-official/nav/api/internal/model"
	"github.com/yueli-official/nav/api/internal/naverr"
)

func TestSafeFaviconMIMERejectsActiveContent(t *testing.T) {
	for _, mime := range []string{"image/svg+xml", "text/html", "application/octet-stream"} {
		if safeFaviconMIME(mime) {
			t.Fatalf("safeFaviconMIME(%q) = true, want false", mime)
		}
	}
	if !safeFaviconMIME("image/png") {
		t.Fatal("safeFaviconMIME(image/png) = false, want true")
	}
}

func TestFetchImageRejectsUnsafeResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/svg":
			response.Header().Set("Content-Type", "image/svg+xml")
			_, _ = response.Write([]byte(`<svg onload="alert(1)"/>`))
		case "/large":
			response.Header().Set("Content-Type", "image/png")
			_, _ = response.Write(bytes.Repeat([]byte{0}, maxFaviconBytes+1))
		case "/missing":
			response.WriteHeader(http.StatusNotFound)
		default:
			response.Header().Set("Content-Type", "text/html")
			_, _ = response.Write([]byte("not an image"))
		}
	}))
	defer server.Close()

	for _, path := range []string{"/svg", "/large", "/missing", "/html"} {
		if _, _, err := fetchImage(context.Background(), server.Client(), server.URL+path); err == nil {
			t.Fatalf("fetchImage(%s) succeeded, want rejection", path)
		}
	}
}

func TestPublicTargetIPRejectsPrivateNetworks(t *testing.T) {
	for _, value := range []string{
		"127.0.0.1", "10.0.0.1", "172.16.0.1", "192.168.1.1", "169.254.1.1", "100.64.0.1",
		"192.0.2.1", "198.18.0.1", "198.51.100.1", "203.0.113.1", "::1", "64:ff9b:1::1",
		"2001:db8::1", "3fff::1", "5f00::1", "fec0::1",
	} {
		if isPublicTargetIP(net.ParseIP(value)) {
			t.Fatalf("isPublicTargetIP(%s) = true, want false", value)
		}
	}
	if !isPublicTargetIP(net.ParseIP("8.8.8.8")) {
		t.Fatal("isPublicTargetIP(8.8.8.8) = false, want true")
	}
}

func TestSafeClientRejectsRedirectToReservedAddress(t *testing.T) {
	client := newSafeHTTPClient(time.Second, true)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://100.64.0.1/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.CheckRedirect(request, []*http.Request{{}}); err == nil {
		t.Fatal("redirect to carrier-grade NAT address succeeded, want rejection")
	}
}

func TestFaviconMissQueuesFetchAndSubsequentReadsUsePersistentCache(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		requests++
		response.Header().Set("Content-Type", "image/png")
		_, _ = response.Write([]byte("cached-png"))
	}))
	defer server.Close()
	store := &fakeStore{links: []*model.Link{{ID: "cached-icon", URL: server.URL, Status: StatusPublished}}}
	service := New(store, Site{})
	service.faviconClient = server.Client()
	var queued func()
	service.faviconRunner = func(task func()) { queued = task }

	_, _, err := service.Favicon(context.Background(), "cached-icon")
	mapped, ok, resolveErr := problem.FromError(err, "test-trace")
	if resolveErr != nil || !ok || mapped.Code != naverr.CodeNotFound {
		t.Fatalf("error = %#v, want nav.not_found", err)
	}
	if requests != 0 || queued == nil {
		t.Fatalf("public miss requests=%d queued=%v, want zero blocking upstream requests and one background job", requests, queued != nil)
	}

	queued()
	data, mime, err := service.Favicon(context.Background(), "cached-icon")
	if err != nil || string(data) != "cached-png" || mime != "image/png" || requests != 1 {
		t.Fatalf("cached favicon data=%q mime=%q requests=%d err=%v", data, mime, requests, err)
	}
	_, _, err = service.Favicon(context.Background(), "cached-icon")
	if err != nil || requests != 1 {
		t.Fatalf("warm cache performed another upstream request: requests=%d err=%v", requests, err)
	}
}

func TestFaviconServesLastKnownGoodWhileFailedRefreshBacksOff(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	now := time.Date(2026, 8, 9, 12, 0, 0, 0, time.UTC)
	store := &fakeStore{
		links: []*model.Link{{ID: "stale-icon", URL: server.URL, Status: StatusPublished}},
		favicons: map[string]*model.FaviconCache{
			"stale-icon": {
				LinkID: "stale-icon", SourceURL: server.URL, Content: []byte("last-known-good"),
				ContentType: "image/png", ContentHash: "hash", FetchedAt: gtime.NewFromTime(now.Add(-8 * 24 * time.Hour)),
				RefreshAfter: gtime.NewFromTime(now.Add(-time.Minute)), LastAttemptAt: gtime.NewFromTime(now.Add(-8 * 24 * time.Hour)),
			},
		},
	}
	service := New(store, Site{})
	service.faviconClient = server.Client()
	service.faviconNow = func() time.Time { return now }
	var queued func()
	service.faviconRunner = func(task func()) { queued = task }

	data, mime, err := service.Favicon(context.Background(), "stale-icon")
	if err != nil || string(data) != "last-known-good" || mime != "image/png" || queued == nil {
		t.Fatalf("stale read data=%q mime=%q queued=%v err=%v", data, mime, queued != nil, err)
	}
	queued()
	queued = nil
	data, _, err = service.Favicon(context.Background(), "stale-icon")
	cached, _ := store.FaviconByLinkID(context.Background(), "stale-icon")
	if err != nil || string(data) != "last-known-good" || queued != nil || cached.LastError == "" || !cached.RefreshAfter.Time.After(now) {
		t.Fatalf("failed refresh did not preserve/back off cache: cache=%#v queued=%v err=%v", cached, queued != nil, err)
	}
}

func TestFaviconNegativeCachePreventsRepeatedRefreshUntilBackoffExpires(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	now := time.Date(2026, 8, 9, 12, 0, 0, 0, time.UTC)
	store := &fakeStore{links: []*model.Link{{ID: "missing-icon", URL: server.URL, Status: StatusPublished}}}
	service := New(store, Site{})
	service.faviconClient = server.Client()
	service.faviconNow = func() time.Time { return now }
	var queued func()
	service.faviconRunner = func(task func()) { queued = task }

	_, _, err := service.Favicon(context.Background(), "missing-icon")
	if err == nil || queued == nil {
		t.Fatalf("initial miss err=%v queued=%v, want not found with background refresh", err, queued != nil)
	}
	queued()
	queued = nil
	_, _, err = service.Favicon(context.Background(), "missing-icon")
	cached, _ := store.FaviconByLinkID(context.Background(), "missing-icon")
	if err == nil || queued != nil || cached == nil || cached.LastError == "" || !cached.RefreshAfter.Time.After(now) {
		t.Fatalf("negative cache did not back off: cache=%#v queued=%v err=%v", cached, queued != nil, err)
	}
}

func TestFetchFaviconUsesOriginIcon(t *testing.T) {
	homepageRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/favicon.ico" {
			response.Header().Set("Content-Type", "image/png")
			_, _ = response.Write([]byte("safe-png-fixture"))
			return
		}
		homepageRequests++
		http.NotFound(response, request)
	}))
	defer server.Close()

	data, mime, err := fetchFavicon(context.Background(), server.Client(), server.URL+"/docs/start")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "safe-png-fixture" || mime != "image/png" {
		t.Fatalf("favicon = %q, %q, want origin PNG", data, mime)
	}
	if homepageRequests != 0 {
		t.Fatalf("homepage requests = %d, want origin icon short-circuit", homepageRequests)
	}
}

func TestFetchFaviconUsesDeclaredIconWhenOriginIconIsMissing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/favicon.ico":
			http.NotFound(response, request)
		case "/docs/start":
			response.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = response.Write([]byte(`<html><head><link href="/assets/site.png" rel="shortcut icon"></head></html>`))
		case "/assets/site.png":
			response.Header().Set("Content-Type", "image/png")
			_, _ = response.Write([]byte("declared-png-fixture"))
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	data, mime, err := fetchFavicon(context.Background(), server.Client(), server.URL+"/docs/start")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "declared-png-fixture" || mime != "image/png" {
		t.Fatalf("favicon = %q, %q, want declared PNG", data, mime)
	}
}

func TestDiscoverIconHref(t *testing.T) {
	document := `<html><head><link rel="stylesheet" href="/app.css"><link href="/assets/site.svg" rel="icon shortcut"></head></html>`
	if got := discoverIconHref(document); got != "/assets/site.svg" {
		t.Fatalf("discoverIconHref() = %q, want /assets/site.svg", got)
	}
}

func TestHTTPLinkCheckerClassifiesResponsesAndFallsBackToGet(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		requests++
		if request.Method == http.MethodHead {
			response.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		response.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	checker := &HTTPLinkChecker{client: server.Client()}
	result := checker.Check(context.Background(), server.URL)
	if result.Status != "healthy" || result.HTTPStatus != http.StatusNoContent {
		t.Fatalf("result = %#v, want healthy 204", result)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want HEAD then GET", requests)
	}
}

func TestHTTPLinkCheckerReportsRedirectWithoutFollowing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Location", "https://example.com/")
		response.WriteHeader(http.StatusMovedPermanently)
	}))
	defer server.Close()

	client := server.Client()
	client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
	result := (&HTTPLinkChecker{client: client}).Check(context.Background(), server.URL)
	if result.Status != "redirected" || result.HTTPStatus != http.StatusMovedPermanently {
		t.Fatalf("result = %#v, want redirected 301", result)
	}
}
