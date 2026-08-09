package deploy

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jose "github.com/go-jose/go-jose/v4"
)

func TestIdentityBindingRequiresOIDCAndPublicUsersContract(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"issuer": server.URL, "authorization_endpoint": server.URL + "/oauth2/authorize",
				"token_endpoint": server.URL + "/oauth2/token", "jwks_uri": server.URL + "/oauth2/jwks.json",
			})
		case "/oauth2/jwks.json":
			_ = json.NewEncoder(writer).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
				Key: &privateKey.PublicKey, KeyID: "identity-key", Use: "sig",
			}}})
		case "/api/v1/users":
			_ = json.NewEncoder(writer).Encode(map[string]any{"users": []any{}})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	err = checkIdentity(context.Background(), IdentityBinding{
		Issuer: server.URL, DiscoveryURL: server.URL + "/.well-known/openid-configuration",
		JWKSURL: server.URL + "/oauth2/jwks.json", PublicUsersBaseURL: server.URL,
		AllowInsecureHTTP: true,
	})
	if err != nil {
		t.Fatalf("checkIdentity() error = %v", err)
	}
}

func TestIdentityBindingRejectsMissingPublicUsersEndpoint(t *testing.T) {
	err := checkIdentity(context.Background(), IdentityBinding{
		Issuer: "https://identity.example", DiscoveryURL: "https://identity.example/discovery",
		JWKSURL: "https://identity.example/jwks", PublicUsersBaseURL: "",
	})
	if err == nil || !strings.Contains(err.Error(), "public users URL is required") {
		t.Fatalf("checkIdentity() error = %v", err)
	}
}
