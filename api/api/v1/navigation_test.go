package v1

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
)

func TestAdminRunChecksRequestAcceptsFifteenIDs(t *testing.T) {
	t.Parallel()

	req := &AdminRunChecksReq{IDs: []string{
		"adobe-color",
		"ambientcg",
		"artstation",
		"bangumi",
		"bbc-learning-english",
		"behance",
		"bigjpg",
		"blender",
		"blueprintue",
		"cambridge-dictionary",
		"chatgpt",
		"codeforces",
		"codepen",
		"convertio",
		"coursera",
	}}

	if err := g.Validator().Data(req).Run(context.Background()); err != nil {
		t.Fatalf("15 ids should satisfy the request contract: %v", err)
	}
}
