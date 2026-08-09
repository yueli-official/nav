package navidentity

import (
	"context"
	"regexp"
	"strings"

	foundationauth "github.com/yueli-official/foundation/go/auth"
)

var publicUserKeyPattern = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{8}$`)

func IsPublicUserKey(value string) bool {
	return publicUserKeyPattern.MatchString(strings.TrimSpace(value))
}

// PublicUserKey returns the platform-stable user reference from a verified
// principal. The subject fallback keeps existing first-party public-subject
// clients working while they migrate to Identity's additive user_key claim.
func PublicUserKey(principal *foundationauth.Principal) (string, bool) {
	if principal == nil {
		return "", false
	}
	kind, _ := principal.Claim("subject_kind")
	if kind != "user" || strings.TrimSpace(principal.Subject) == "" {
		return "", false
	}
	if claim, ok := principal.Claim("user_key"); ok {
		if value, valueOK := claim.(string); valueOK && IsPublicUserKey(value) {
			return strings.TrimSpace(value), true
		}
	}
	if IsPublicUserKey(principal.Subject) {
		return strings.TrimSpace(principal.Subject), true
	}
	return "", false
}

func FromContext(ctx context.Context) (string, bool) {
	principal, ok := foundationauth.FromContext(ctx)
	if !ok {
		return "", false
	}
	return PublicUserKey(principal)
}
