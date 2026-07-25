package navauthz

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yueli-official/foundation/go/authorization"
)

func SyncResourceScopes(ctx context.Context, db *sql.DB, runtime Runtime) error {
	rows, err := db.QueryContext(ctx, `SELECT id::text FROM nav_categories ORDER BY id`)
	if err != nil {
		return fmt.Errorf("list categories for authorization scope sync: %w", err)
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return err
		}
		if err := ensureScope(ctx, runtime, authorization.RegisterScopeCommand{
			ID: CategoryScopeID(id), Type: ScopeCategory, ParentID: RootScopeID,
		}); err != nil {
			_ = rows.Close()
			return err
		}
	}
	if err := rows.Close(); err != nil {
		return err
	}

	rows, err = db.QueryContext(ctx, `SELECT id::text, category_id::text FROM nav_groups ORDER BY id`)
	if err != nil {
		return fmt.Errorf("list groups for authorization scope sync: %w", err)
	}
	for rows.Next() {
		var id, categoryID string
		if err := rows.Scan(&id, &categoryID); err != nil {
			_ = rows.Close()
			return err
		}
		if err := ensureScope(ctx, runtime, authorization.RegisterScopeCommand{
			ID: GroupScopeID(id), Type: ScopeGroup, ParentID: CategoryScopeID(categoryID),
		}); err != nil {
			_ = rows.Close()
			return err
		}
	}
	if err := rows.Close(); err != nil {
		return err
	}

	rows, err = db.QueryContext(ctx, `SELECT id::text, group_id::text FROM nav_links ORDER BY id`)
	if err != nil {
		return fmt.Errorf("list links for authorization scope sync: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, groupID string
		if err := rows.Scan(&id, &groupID); err != nil {
			return err
		}
		if err := ensureScope(ctx, runtime, authorization.RegisterScopeCommand{
			ID: LinkScopeID(id), Type: ScopeLink, ParentID: GroupScopeID(groupID),
		}); err != nil {
			return err
		}
	}
	return rows.Err()
}

func ensureScope(ctx context.Context, runtime Runtime, command authorization.RegisterScopeCommand) error {
	if runtime == nil {
		return &authorization.Error{Kind: authorization.ErrorUnavailable, Field: "runtime", Message: "is not configured"}
	}
	_, err := runtime.RegisterScope(ctx, command)
	if err == nil {
		return nil
	}
	if authorization.Is(err, authorization.ErrorConflict) {
		_, moveErr := runtime.ReparentScope(ctx, authorization.ReparentScopeCommand{
			ID: command.ID, ParentID: command.ParentID,
		})
		if moveErr == nil {
			return nil
		}
		return fmt.Errorf("repair authorization scope %q: %w", command.ID, moveErr)
	}
	return fmt.Errorf("register authorization scope %q: %w", command.ID, err)
}
