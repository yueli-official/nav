// Package identityclient adapts Identity's public-user HTTP contract. It never
// reads credentials, private account state, or Identity implementation code.
package identityclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	foundationhttpclient "github.com/yueli-official/foundation/go/httpclient"
)

const batchSize = 100

type MediaRef struct {
	MediaKey string `json:"mediaKey"`
}

type PublicUser struct {
	UserKey     string    `json:"userKey"`
	Handle      string    `json:"handle"`
	DisplayName string    `json:"displayName"`
	Avatar      *MediaRef `json:"avatar"`
}

type Client interface {
	GetMany(context.Context, []string) (map[string]PublicUser, error)
}

type HTTPClient struct {
	base   string
	client *http.Client
}

func NewHTTP(baseURL string, client *http.Client) *HTTPClient {
	if client == nil {
		client = &http.Client{Timeout: 2 * time.Second}
	}
	return &HTTPClient{base: strings.TrimRight(baseURL, "/"), client: client}
}

func (client *HTTPClient) GetMany(ctx context.Context, userKeys []string) (map[string]PublicUser, error) {
	unique := make([]string, 0, len(userKeys))
	seen := make(map[string]struct{}, len(userKeys))
	for _, key := range userKeys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, key)
	}
	result := make(map[string]PublicUser, len(unique))
	for start := 0; start < len(unique); start += batchSize {
		end := min(start+batchSize, len(unique))
		query := url.Values{"ids": {strings.Join(unique[start:end], ",")}}
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.base+"/api/v1/users?"+query.Encode(), nil)
		if err != nil {
			return nil, fmt.Errorf("build public users request: %w", err)
		}
		response, err := client.client.Do(request)
		if err != nil {
			return nil, fmt.Errorf("get public users: %w", err)
		}
		decoded, decodeErr := foundationhttpclient.DecodeJSON[struct {
			Users []PublicUser `json:"users"`
		}](response, foundationhttpclient.Limits{})
		response.Body.Close()
		if decodeErr != nil {
			return nil, fmt.Errorf("decode public users: %w", decodeErr)
		}
		for _, user := range decoded.Users {
			if user.UserKey != "" {
				result[user.UserKey] = user
			}
		}
	}
	return result, nil
}
