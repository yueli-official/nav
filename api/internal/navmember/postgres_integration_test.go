package navmember

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/yueli-official/nav/api/internal/identityclient"
)

type profileFixture map[string]identityclient.PublicUser

func (fixture profileFixture) GetMany(_ context.Context, keys []string) (map[string]identityclient.PublicUser, error) {
	result := map[string]identityclient.PublicUser{}
	for _, key := range keys {
		if profile, ok := fixture[key]; ok {
			result[key] = profile
		}
	}
	return result, nil
}

func TestMembershipLifecycleOnPostgres(t *testing.T) {
	database := openMembershipPostgres(t)
	const memberKey = "TestB234"
	const adminKey = "TestA123"
	clock := time.Date(2026, 8, 9, 4, 0, 0, 0, time.UTC)
	service := New(database, profileFixture{memberKey: {
		UserKey: memberKey, Handle: "alice", DisplayName: "Alice",
		Avatar: &identityclient.MediaRef{MediaKey: "31Pj0mXv7cfR5fdZIUvra"},
	}})
	service.now = func() time.Time { return clock }

	first, err := service.Ensure(context.Background(), memberKey)
	if err != nil || !first.Created || !first.NeedsJoinReconcile || first.Member.DisplayName != "Alice" {
		t.Fatalf("first Ensure() = %#v, %v", first, err)
	}
	if err := service.MarkJoinReconciled(context.Background(), memberKey); err != nil {
		t.Fatalf("MarkJoinReconciled() error = %v", err)
	}
	clock = clock.Add(5 * time.Minute)
	second, err := service.Ensure(context.Background(), memberKey)
	if err != nil || second.Created || second.NeedsJoinReconcile || !second.Member.LastSeenAt.Equal(first.Member.LastSeenAt) {
		t.Fatalf("second Ensure() = %#v, %v", second, err)
	}
	clock = clock.Add(16 * time.Minute)
	third, err := service.Ensure(context.Background(), memberKey)
	if err != nil || !third.Member.LastSeenAt.Equal(clock) {
		t.Fatalf("throttled Ensure() = %#v, %v", third, err)
	}
	if _, err := service.SetStatus(context.Background(), SetStatusCommand{
		UserKey: memberKey, ActorUserKey: memberKey, Status: StatusSuspended,
	}); !errors.Is(err, ErrSelfSuspend) {
		t.Fatalf("self suspend error = %v", err)
	}
	if _, err := service.SetStatus(context.Background(), SetStatusCommand{
		UserKey: memberKey, ActorUserKey: adminKey, Status: StatusSuspended,
	}); !errors.Is(err, ErrReasonRequired) {
		t.Fatalf("missing reason error = %v", err)
	}
	suspended, err := service.SetStatus(context.Background(), SetStatusCommand{
		UserKey: memberKey, ActorUserKey: adminKey, Status: StatusSuspended, Reason: "policy review",
	})
	if err != nil || suspended.Status != StatusSuspended || suspended.SuspendedBy != adminKey {
		t.Fatalf("SetStatus() = %#v, %v", suspended, err)
	}
	page, err := service.List(context.Background(), Query{Search: "alice", Status: StatusSuspended, Page: 1, Size: 20})
	if err != nil || page.Total != 1 || page.Counts.Suspended != 1 || len(page.Members) != 1 {
		t.Fatalf("List() = %#v, %v", page, err)
	}
	var events int
	if err := database.QueryRow(`SELECT COUNT(*) FROM nav_membership_events`).Scan(&events); err != nil || events != 2 {
		t.Fatalf("membership events = %d, %v", events, err)
	}
}

func openMembershipPostgres(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("NAV_MEMBERSHIP_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("NAV_MEMBERSHIP_POSTGRES_DSN is not configured")
	}
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	database.SetMaxOpenConns(1)
	schema := fmt.Sprintf("nav_membership_test_%d", time.Now().UnixNano())
	if _, err := database.Exec("CREATE SCHEMA " + pq.QuoteIdentifier(schema)); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec("SET search_path TO " + pq.QuoteIdentifier(schema)); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`CREATE TABLE nav_links (submitter_sub TEXT NOT NULL DEFAULT '')`); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../manifest/sql/migrations/0012_membership_v1.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = database.Exec("DROP SCHEMA " + pq.QuoteIdentifier(schema) + " CASCADE")
		_ = database.Close()
	})
	return database
}
