// Package deploy owns Nav's production binding checks. Business packages only
// depend on the OIDC behavior they consume, never on Identity internals.
package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	jose "github.com/go-jose/go-jose/v4"
)

const maxResponseBytes = int64(1 << 20)

type IdentityBinding struct {
	Issuer             string
	DiscoveryURL       string
	JWKSURL            string
	PublicUsersBaseURL string
	AllowInsecureHTTP  bool
}

func WaitForIdentity(
	ctx context.Context,
	binding IdentityBinding,
	interval time.Duration,
) error {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	var lastErr error
	for {
		if err := checkIdentity(ctx, binding); err == nil {
			return nil
		} else {
			lastErr = err
		}
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf(
				"identity binding did not converge: %w",
				errors.Join(lastErr, ctx.Err()),
			)
		case <-timer.C:
		}
	}
}

func checkIdentity(ctx context.Context, binding IdentityBinding) error {
	binding.Issuer = strings.TrimRight(strings.TrimSpace(binding.Issuer), "/")
	binding.DiscoveryURL = strings.TrimSpace(binding.DiscoveryURL)
	binding.JWKSURL = strings.TrimSpace(binding.JWKSURL)
	binding.PublicUsersBaseURL = strings.TrimRight(strings.TrimSpace(binding.PublicUsersBaseURL), "/")
	if binding.DiscoveryURL == "" && binding.Issuer != "" {
		binding.DiscoveryURL = binding.Issuer + "/.well-known/openid-configuration"
	}
	if binding.JWKSURL == "" && binding.Issuer != "" {
		binding.JWKSURL = binding.Issuer + "/oauth2/jwks.json"
	}
	for name, value := range map[string]string{
		"issuer":           binding.Issuer,
		"discovery URL":    binding.DiscoveryURL,
		"JWKS URL":         binding.JWKSURL,
		"public users URL": binding.PublicUsersBaseURL,
	} {
		if value == "" {
			return fmt.Errorf("identity %s is required", name)
		}
		if err := validateURL(value, binding.AllowInsecureHTTP); err != nil {
			return fmt.Errorf("identity %s: %w", name, err)
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var discovery struct {
		Issuer                string `json:"issuer"`
		AuthorizationEndpoint string `json:"authorization_endpoint"`
		TokenEndpoint         string `json:"token_endpoint"`
		JWKSURI               string `json:"jwks_uri"`
	}
	if err := getJSON(ctx, client, binding.DiscoveryURL, &discovery); err != nil {
		return fmt.Errorf("identity discovery: %w", err)
	}
	if strings.TrimRight(discovery.Issuer, "/") != binding.Issuer {
		return fmt.Errorf(
			"identity discovery issuer %q does not match expected %q",
			discovery.Issuer,
			binding.Issuer,
		)
	}
	if discovery.AuthorizationEndpoint == "" || discovery.TokenEndpoint == "" ||
		discovery.JWKSURI == "" {
		return errors.New("identity discovery is missing OIDC endpoints")
	}
	var set jose.JSONWebKeySet
	if err := getJSON(ctx, client, binding.JWKSURL, &set); err != nil {
		return fmt.Errorf("identity JWKS: %w", err)
	}
	usableKey := false
	for index := range set.Keys {
		key := &set.Keys[index]
		if key.Valid() && key.IsPublic() && strings.TrimSpace(key.KeyID) != "" &&
			(key.Use == "" || key.Use == "sig") {
			usableKey = true
			break
		}
	}
	if !usableKey {
		return errors.New("identity JWKS has no usable public signing key")
	}
	var publicUsers map[string]json.RawMessage
	probeURL := binding.PublicUsersBaseURL + "/api/v1/users?ids=TestA123"
	if err := getJSON(ctx, client, probeURL, &publicUsers); err != nil {
		return fmt.Errorf("identity public users: %w", err)
	}
	users, exists := publicUsers["users"]
	if !exists || len(users) == 0 || users[0] != '[' {
		return errors.New("identity public users response is missing users array")
	}
	return nil
}

func validateURL(value string, allowInsecureHTTP bool) error {
	parsed, err := url.Parse(value)
	if err != nil || !parsed.IsAbs() || parsed.Host == "" {
		return errors.New("must be an absolute URL")
	}
	if parsed.User != nil || parsed.Fragment != "" {
		return errors.New("must not contain credentials or a fragment")
	}
	switch strings.ToLower(parsed.Scheme) {
	case "https":
		return nil
	case "http":
		if allowInsecureHTTP || isLoopbackHost(parsed.Hostname()) {
			return nil
		}
		return errors.New("plain HTTP requires an explicit insecure-HTTP binding opt-in")
	default:
		return errors.New("scheme must be HTTP or HTTPS")
	}
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	address := net.ParseIP(host)
	return address != nil && address.IsLoopback()
}

func getJSON(ctx context.Context, client *http.Client, endpoint string, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return err
	}
	if int64(len(body)) > maxResponseBytes {
		return fmt.Errorf("response exceeds %d bytes", maxResponseBytes)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("GET %s returned %s", endpoint, response.Status)
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode GET %s: %w", endpoint, err)
	}
	return nil
}
