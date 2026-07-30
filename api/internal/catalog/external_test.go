package catalog

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestFaviconMapsUpstreamFailureToNotFound(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	store := &fakeStore{links: []*model.Link{{ID: "missing-icon", URL: server.URL, Status: StatusPublished}}}
	service := New(store, Site{})
	service.faviconClient = server.Client()
	_, _, err := service.Favicon(context.Background(), "missing-icon")
	mapped, ok, resolveErr := problem.FromError(err, "test-trace")
	if resolveErr != nil || !ok || mapped.Code != naverr.CodeNotFound {
		t.Fatalf("error = %#v, want nav.not_found", err)
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
