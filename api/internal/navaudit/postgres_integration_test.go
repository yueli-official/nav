package navaudit_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/yueli-official/foundation/go/audit"

	"github.com/yueli-official/nav/api/internal/navaudit"
)

func TestPostgresHookRecordsCommittedNavigationAction(t *testing.T) {
	database := openPostgres(t)
	journal, err := navaudit.New(context.Background(), database, "integration")
	if err != nil {
		t.Fatal(err)
	}
	tx, err := database.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	hook := journal.Hook(
		context.Background(), navaudit.ActionNavigationPublished, "event-commit",
		audit.Target{Type: "nav.link", ID: "link-1"},
		navaudit.Evidence{},
	)
	if err := hook(context.Background(), tx); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	page, err := journal.Reader().Query(context.Background(), audit.Query{})
	if err != nil || len(page.Events) != 1 || page.Events[0].Action.Name != "nav.navigation.published" {
		t.Fatalf("committed audit = %#v, %v", page, err)
	}
}

func openPostgres(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("AUDIT_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("AUDIT_POSTGRES_DSN is not configured")
	}
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	database.SetMaxOpenConns(1)
	schema := fmt.Sprintf("nav_audit_test_%d", time.Now().UnixNano())
	if _, err := database.Exec("CREATE SCHEMA " + pq.QuoteIdentifier(schema)); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec("SET search_path TO " + pq.QuoteIdentifier(schema)); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(audit.PostgresSchemaUp()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = database.Exec(audit.PostgresSchemaDown())
		_, _ = database.Exec("DROP SCHEMA " + pq.QuoteIdentifier(schema))
		_ = database.Close()
	})
	return database
}
