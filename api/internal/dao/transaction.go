package dao

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

var ErrSiteSettingsRevisionConflict = errors.New("navigation site settings revision conflict")

type TransactionHook func(context.Context, *sql.Tx) error

func runTransactionHook(ctx context.Context, tx gdb.TX, hook TransactionHook) error {
	if hook == nil {
		return nil
	}
	return hook(ctx, tx.GetSqlTX())
}

func ComposeTransactionHooks(hooks ...TransactionHook) TransactionHook {
	return func(ctx context.Context, tx *sql.Tx) error {
		for _, hook := range hooks {
			if hook != nil {
				if err := hook(ctx, tx); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func (p *PG) SaveSiteSettingsWithHook(
	ctx context.Context,
	searchPlaceholder string,
	expectedRevision uint64,
	hook TransactionHook,
) error {
	return p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Ctx(ctx).Exec(`
UPDATE nav_site_settings
SET search_placeholder = ?, runtime_revision = runtime_revision + 1, updated_at = now()
WHERE id = 1 AND runtime_revision = ?
`, strings.TrimSpace(searchPlaceholder), expectedRevision)
		if err != nil {
			return err
		}
		changed, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if changed != 1 {
			return ErrSiteSettingsRevisionConflict
		}
		if hook == nil {
			return nil
		}
		return hook(ctx, tx.GetSqlTX())
	})
}

func (p *PG) CutoverSiteSettingsWithHook(
	ctx context.Context,
	searchPlaceholder string,
	hook TransactionHook,
) error {
	return p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Ctx(ctx).Exec(`
INSERT INTO nav_site_settings (
    id, name, title, description, search_placeholder, footer_tagline, runtime_revision, updated_at
) VALUES (1, NULL, NULL, NULL, ?, NULL, 1, now())
ON CONFLICT (id) DO UPDATE SET
    name = NULL,
    title = NULL,
    description = NULL,
    search_placeholder = EXCLUDED.search_placeholder,
    footer_tagline = NULL,
    updated_at = now()
`, strings.TrimSpace(searchPlaceholder)); err != nil {
			return err
		}
		if hook == nil {
			return nil
		}
		return hook(ctx, tx.GetSqlTX())
	})
}
