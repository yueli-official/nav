package controller

import (
	"context"
	"slices"

	"github.com/gogf/gf/v2/frame/g"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	"platform/products/nav/api/internal/naverr"
)

func requireAdmin(ctx context.Context) error {
	principal, ok := foundationauth.FromContext(ctx)
	if !ok || !slices.Contains(g.Cfg().MustGet(ctx, "nav.operatorSubs").Strings(), principal.Subject) {
		return naverr.Forbidden()
	}
	return nil
}
