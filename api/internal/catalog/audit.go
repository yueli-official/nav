package catalog

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/yueli-official/foundation/go/audit"

	"github.com/yueli-official/nav/api/internal/dao"
	"github.com/yueli-official/nav/api/internal/model"
	"github.com/yueli-official/nav/api/internal/navaudit"
)

type auditedStore interface {
	InsertLinkWithHook(context.Context, *model.Link, dao.TransactionHook) error
	UpdateLinkWithHook(context.Context, *model.Link, dao.TransactionHook) (bool, error)
	DeleteLinkWithHook(context.Context, string, dao.TransactionHook) (bool, error)
	BulkUpdateLinksWithHook(context.Context, []string, string, dao.TransactionHook) (int, error)
	BulkDeleteLinksWithHook(context.Context, []string, dao.TransactionHook) (int, error)
	InsertCategoryWithHook(context.Context, *model.Category, dao.TransactionHook) error
	UpdateCategoryWithHook(context.Context, *model.Category, dao.TransactionHook) (bool, error)
	DeleteCategoryWithHook(context.Context, string, dao.TransactionHook) (bool, error)
	InsertGroupWithHook(context.Context, *model.Group, dao.TransactionHook) error
	UpdateGroupWithHook(context.Context, *model.Group, dao.TransactionHook) (bool, error)
	DeleteGroupWithHook(context.Context, string, dao.TransactionHook) (bool, error)
	RenameTagWithHook(context.Context, string, string, dao.TransactionHook) (int, error)
	DeleteTagWithHook(context.Context, string, dao.TransactionHook) (int, error)
}

func (s *Service) SetAudit(journal *navaudit.Journal) {
	s.audit = journal
}

func (s *Service) linkAuditHook(
	ctx context.Context,
	action navaudit.Action,
	targetType string,
	targetID string,
	digest string,
	count int,
) dao.TransactionHook {
	if s.audit == nil || action == "" {
		return nil
	}
	return s.audit.Hook(
		ctx, action, uuid.NewString(),
		audit.Target{Type: targetType, ID: targetID},
		navaudit.Evidence{Digest: digest, Count: uint64(count)},
	)
}

func (s *Service) taxonomyAuditHook(ctx context.Context, kind, id string, count int) dao.TransactionHook {
	if s.audit == nil {
		return nil
	}
	sum := sha256.Sum256([]byte(kind + "\x00" + id))
	return s.audit.Hook(
		ctx, navaudit.ActionTaxonomyChanged, uuid.NewString(),
		audit.Target{Type: "nav.taxonomy", ID: kind + ":" + id},
		navaudit.Evidence{Digest: hex.EncodeToString(sum[:]), Count: uint64(count)},
	)
}

func linkDigest(link *model.Link) string {
	raw := []byte(
		link.ID + "\x00" + link.CategoryID + "\x00" + link.GroupID + "\x00" +
			link.Title + "\x00" + link.URL + "\x00" + link.Status,
	)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func navigationAction(status string) navaudit.Action {
	switch status {
	case "published":
		return navaudit.ActionNavigationPublished
	case "archived":
		return navaudit.ActionNavigationArchived
	default:
		return ""
	}
}
