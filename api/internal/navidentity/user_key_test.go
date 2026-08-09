package navidentity_test

import (
	"context"
	"testing"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	"github.com/yueli-official/nav/api/internal/navidentity"
	"github.com/yueli-official/nav/api/internal/testidentity"
)

func TestPublicUserKeyPrefersVerifiedClaimForPairwiseSubject(t *testing.T) {
	const key = "TestA123"
	principal := testidentity.UserWithKey(t, "pairwise-subject", key, nil, nil)
	if got, ok := navidentity.PublicUserKey(principal); !ok || got != key {
		t.Fatalf("PublicUserKey() = %q, %v; want %q, true", got, ok, key)
	}
}

func TestPublicUserKeyFallsBackToFirstPartyPublicSubject(t *testing.T) {
	const key = "TestA123"
	principal := testidentity.User(t, key, nil, nil)
	if got, ok := navidentity.PublicUserKey(principal); !ok || got != key {
		t.Fatalf("PublicUserKey() = %q, %v; want fallback", got, ok)
	}
}

func TestPublicUserKeyRejectsUnstableSubjectWithoutClaim(t *testing.T) {
	ctx := foundationauth.NewContext(context.Background(), testidentity.User(t, "pairwise-subject", nil, nil))
	if got, ok := navidentity.FromContext(ctx); ok || got != "" {
		t.Fatalf("FromContext() = %q, %v; want rejected", got, ok)
	}
}
