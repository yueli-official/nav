package identityclient_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yueli-official/nav/api/internal/identityclient"
)

func TestHTTPClientBatchResolvesAndDeduplicatesPublicUsers(t *testing.T) {
	const first = "TestA123"
	const second = "TestB234"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if got := request.URL.Query().Get("ids"); got != first+","+second {
			t.Fatalf("ids = %q", got)
		}
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(map[string]any{"users": []map[string]any{
			{"userKey": first, "handle": "alice", "displayName": "Alice"},
			{"userKey": second, "displayName": "Bob"},
		}})
	}))
	defer server.Close()

	users, err := identityclient.NewHTTP(server.URL, server.Client()).GetMany(
		context.Background(), []string{first, first, "", second},
	)
	if err != nil || len(users) != 2 || users[first].Handle != "alice" || users[second].DisplayName != "Bob" {
		t.Fatalf("GetMany() = %#v, %v", users, err)
	}
}

func TestHTTPClientSurfacesUnavailableProfilesForResilientFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		http.Error(writer, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()
	if _, err := identityclient.NewHTTP(server.URL, server.Client()).GetMany(
		context.Background(), []string{"TestA123"},
	); err == nil {
		t.Fatal("GetMany() error = nil, want unavailable")
	}
}
