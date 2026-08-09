// Command bindingcheck validates Nav's concrete Identity binding before the
// runtime process starts.
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yueli-official/nav/api/internal/deploy"
)

func main() {
	timeout := durationFromEnvironment("NAV_BINDING_TIMEOUT", 3*time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := deploy.WaitForIdentity(ctx, deploy.IdentityBinding{
		Issuer:             os.Getenv("NAV_IDENTITY_ISSUER"),
		DiscoveryURL:       os.Getenv("NAV_IDENTITY_DISCOVERY_URL"),
		JWKSURL:            os.Getenv("NAV_IDENTITY_JWKS_URL"),
		PublicUsersBaseURL: os.Getenv("NAV_IDENTITY_INTERNAL_URL"),
		AllowInsecureHTTP:  boolFromEnvironment("NAV_BINDING_ALLOW_INSECURE_IDENTITY_HTTP"),
	}, 2*time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, "nav binding check:", err)
		os.Exit(1)
	}
	fmt.Println("Nav Identity binding is compatible")
}

func durationFromEnvironment(name string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		fmt.Fprintf(os.Stderr, "%s must be a positive duration\n", name)
		os.Exit(2)
	}
	return value
}

func boolFromEnvironment(name string) bool {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return false
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s must be a boolean\n", name)
		os.Exit(2)
	}
	return value
}
