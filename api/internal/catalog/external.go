package catalog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"regexp"
	"strings"
	"time"

	"platform/products/nav/api/internal/model"
	"platform/products/nav/api/internal/naverr"
)

const (
	linkCheckTimeout = 6 * time.Second
	maxFaviconBytes  = 1 << 20
	maxHTMLBytes     = 512 << 10
)

type LinkChecker interface {
	Check(context.Context, string) model.LinkHealth
}

type HTTPLinkChecker struct {
	client *http.Client
}

func NewHTTPLinkChecker() *HTTPLinkChecker {
	return &HTTPLinkChecker{client: newSafeHTTPClient(linkCheckTimeout, false)}
}

func (c *HTTPLinkChecker) Check(ctx context.Context, rawURL string) model.LinkHealth {
	started := time.Now()
	response, err := c.request(ctx, http.MethodHead, rawURL)
	if err == nil && (response.StatusCode == http.StatusMethodNotAllowed || response.StatusCode == http.StatusNotImplemented) {
		response.Body.Close()
		response, err = c.request(ctx, http.MethodGet, rawURL)
	}
	latency := int(time.Since(started).Milliseconds())
	if err != nil {
		status := "error"
		var netErr net.Error
		if (errors.As(err, &netErr) && netErr.Timeout()) || errors.Is(err, context.DeadlineExceeded) {
			status = "timeout"
		}
		return model.LinkHealth{Status: status, LatencyMS: latency, Error: truncateExternalError(err.Error())}
	}
	defer response.Body.Close()
	status := "broken"
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		status = "healthy"
	} else if response.StatusCode >= 300 && response.StatusCode < 400 {
		status = "redirected"
	}
	return model.LinkHealth{Status: status, HTTPStatus: response.StatusCode, LatencyMS: latency}
}

func (c *HTTPLinkChecker) request(ctx context.Context, method, rawURL string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "YueliNavHealth/1.0 (+https://nav.localhost)")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,*/*;q=0.8")
	return c.client.Do(request)
}

func (s *Service) Favicon(ctx context.Context, id string) ([]byte, string, error) {
	link, err := s.store.LinkByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, "", err
	}
	if link == nil || link.Status != StatusPublished {
		return nil, "", naverr.NotFound(id)
	}
	data, mime, err := fetchFavicon(ctx, s.faviconClient, link.URL)
	if err != nil {
		return nil, "", naverr.FaviconNotFound(id)
	}
	return data, mime, nil
}

func fetchFavicon(ctx context.Context, client *http.Client, rawURL string) ([]byte, string, error) {
	pageURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", err
	}
	originIcon := pageURL.ResolveReference(&url.URL{Path: "/favicon.ico"})
	if data, mime, iconErr := fetchImage(ctx, client, originIcon.String()); iconErr == nil {
		return data, mime, nil
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL.String(), nil)
	if err != nil {
		return nil, "", err
	}
	request.Header.Set("User-Agent", "YueliNavFavicon/1.0 (+https://nav.localhost)")
	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, "", fmt.Errorf("homepage returned HTTP %d", response.StatusCode)
	}
	html, err := io.ReadAll(io.LimitReader(response.Body, maxHTMLBytes))
	if err != nil {
		return nil, "", err
	}
	href := discoverIconHref(string(html))
	if href == "" {
		return nil, "", errors.New("favicon not declared")
	}
	iconURL, err := response.Request.URL.Parse(href)
	if err != nil {
		return nil, "", err
	}
	return fetchImage(ctx, client, iconURL.String())
}

var (
	linkTagPattern = regexp.MustCompile(`(?is)<link\b[^>]*>`)
	relPattern     = regexp.MustCompile(`(?is)\brel\s*=\s*["']([^"']+)["']`)
	hrefPattern    = regexp.MustCompile(`(?is)\bhref\s*=\s*["']([^"']+)["']`)
)

func discoverIconHref(document string) string {
	for _, tag := range linkTagPattern.FindAllString(document, -1) {
		relMatch, hrefMatch := relPattern.FindStringSubmatch(tag), hrefPattern.FindStringSubmatch(tag)
		if len(relMatch) == 2 && len(hrefMatch) == 2 && strings.Contains(strings.ToLower(relMatch[1]), "icon") {
			return strings.TrimSpace(hrefMatch[1])
		}
	}
	return ""
}

func fetchImage(ctx context.Context, client *http.Client, rawURL string) ([]byte, string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", err
	}
	request.Header.Set("User-Agent", "YueliNavFavicon/1.0 (+https://nav.localhost)")
	request.Header.Set("Accept", "image/avif,image/webp,image/png,image/jpeg,image/gif,image/x-icon,*/*;q=0.1")
	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, "", fmt.Errorf("favicon returned HTTP %d", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, maxFaviconBytes+1))
	if err != nil {
		return nil, "", err
	}
	if len(data) == 0 || len(data) > maxFaviconBytes {
		return nil, "", errors.New("favicon is empty or too large")
	}
	mime := strings.ToLower(strings.TrimSpace(strings.Split(response.Header.Get("Content-Type"), ";")[0]))
	if mime == "" || mime == "application/octet-stream" {
		mime = http.DetectContentType(data)
	}
	if !safeFaviconMIME(mime) {
		return nil, "", fmt.Errorf("favicon content type %q is not a safe image", mime)
	}
	return data, mime, nil
}

func safeFaviconMIME(value string) bool {
	switch value {
	case "image/png", "image/jpeg", "image/gif", "image/webp", "image/avif", "image/x-icon", "image/vnd.microsoft.icon":
		return true
	default:
		return false
	}
}

func newSafeHTTPClient(timeout time.Duration, followRedirects bool) *http.Client {
	transport := &http.Transport{
		DialContext:           dialPublicAddress,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          32,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   4 * time.Second,
		ResponseHeaderTimeout: timeout,
	}
	client := &http.Client{Transport: transport, Timeout: timeout}
	client.CheckRedirect = func(request *http.Request, via []*http.Request) error {
		if !followRedirects {
			return http.ErrUseLastResponse
		}
		if len(via) >= 3 {
			return errors.New("too many redirects")
		}
		return ensurePublicHost(request.Context(), request.URL.Hostname())
	}
	return client
}

func dialPublicAddress(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	addresses, err := publicTargetAddresses(ctx, host)
	if err != nil {
		return nil, err
	}
	var lastErr error
	for _, ip := range addresses {
		connection, dialErr := (&net.Dialer{Timeout: 4 * time.Second}).DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		if dialErr == nil {
			return connection, nil
		}
		lastErr = dialErr
	}
	return nil, lastErr
}

func isPublicTargetIP(ip net.IP) bool {
	address, ok := netip.AddrFromSlice(ip)
	if !ok {
		return false
	}
	address = address.Unmap()
	if !address.IsGlobalUnicast() || address.IsPrivate() || address.IsLoopback() || address.IsLinkLocalUnicast() || address.IsUnspecified() {
		return false
	}
	for _, prefix := range deniedTargetPrefixes {
		if prefix.Contains(address) {
			return false
		}
	}
	return true
}

var deniedTargetPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"), netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("192.0.0.0/24"), netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.88.99.0/24"), netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"), netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("240.0.0.0/4"), netip.MustParsePrefix("64:ff9b::/96"),
	netip.MustParsePrefix("64:ff9b:1::/48"), netip.MustParsePrefix("100::/64"),
	netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2001:db8::/32"), netip.MustParsePrefix("2002::/16"),
	netip.MustParsePrefix("3fff::/20"), netip.MustParsePrefix("5f00::/16"),
	netip.MustParsePrefix("fec0::/10"),
}

func ensurePublicHost(ctx context.Context, host string) error {
	_, err := publicTargetAddresses(ctx, host)
	return err
}

func publicTargetAddresses(ctx context.Context, host string) ([]net.IP, error) {
	if parsed := net.ParseIP(host); parsed != nil {
		if isPublicTargetIP(parsed) {
			return []net.IP{parsed}, nil
		}
		return nil, errors.New("target is private, reserved, or non-routable")
	}
	addresses, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	public := make([]net.IP, 0, len(addresses))
	for _, address := range addresses {
		if isPublicTargetIP(address.IP) {
			public = append(public, address.IP)
		}
	}
	if len(public) == 0 {
		return nil, errors.New("target resolves only to private, reserved, or non-routable addresses")
	}
	return public, nil
}

func truncateExternalError(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 240 {
		return value[:240]
	}
	return value
}
