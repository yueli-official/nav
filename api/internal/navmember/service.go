// Package navmember owns the product-local relationship between an Identity
// user and Nav. It deliberately excludes credentials, global profile truth,
// roles, grants, and applications.
package navmember

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/yueli-official/foundation/go/identifier"
	"github.com/yueli-official/nav/api/internal/identityclient"
	"github.com/yueli-official/nav/api/internal/navidentity"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
)

var (
	ErrInvalidUserKey = errors.New("navmember: invalid public user key")
	ErrInvalidStatus  = errors.New("navmember: invalid membership status")
	ErrReasonRequired = errors.New("navmember: suspension reason is required")
	ErrSelfSuspend    = errors.New("navmember: administrator cannot suspend own membership")
	ErrNotFound       = errors.New("navmember: membership not found")
)

type Member struct {
	ID               string
	UserKey          string
	Status           Status
	DisplayName      string
	Handle           string
	AvatarMediaKey   string
	ProfileSyncedAt  time.Time
	JoinReconciledAt time.Time
	JoinedAt         time.Time
	LastSeenAt       time.Time
	SuspendedAt      time.Time
	SuspendedBy      string
	SuspensionReason string
	UpdatedAt        time.Time
	SubmissionCount  int
}

type EnsureResult struct {
	Member             Member
	Created            bool
	NeedsJoinReconcile bool
}

type Query struct {
	Search            string
	Status            Status
	Page              int
	Size              int
	ConstrainUserKeys bool
	UserKeys          []string
}

type Counts struct {
	All       int
	Active    int
	Suspended int
}

type Page struct {
	Members []Member
	Counts  Counts
	Total   int
	Page    int
	Size    int
}

type SetStatusCommand struct {
	UserKey      string
	Status       Status
	ActorUserKey string
	Reason       string
}

// Directory is the complete membership seam used by HTTP composition. It is
// intentionally narrower than the underlying SQL model.
type Directory interface {
	Ensure(context.Context, string) (EnsureResult, error)
	MarkJoinReconciled(context.Context, string) error
	Get(context.Context, string) (Member, error)
	List(context.Context, Query) (Page, error)
	SetStatus(context.Context, SetStatusCommand) (Member, error)
}

type Service struct {
	db                     *sql.DB
	profiles               identityclient.Client
	now                    func() time.Time
	touchInterval          time.Duration
	profileRefreshInterval time.Duration
}

func New(db *sql.DB, profiles identityclient.Client) *Service {
	return &Service{
		db: db, profiles: profiles, now: func() time.Time { return time.Now().UTC() },
		touchInterval: 15 * time.Minute, profileRefreshInterval: 15 * time.Minute,
	}
}

func (service *Service) Ensure(ctx context.Context, userKey string) (EnsureResult, error) {
	if service == nil || service.db == nil {
		return EnsureResult{}, errors.New("navmember: service is not configured")
	}
	if !navidentity.IsPublicUserKey(userKey) {
		return EnsureResult{}, ErrInvalidUserKey
	}
	now := service.now()
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return EnsureResult{}, fmt.Errorf("begin ensure membership: %w", err)
	}
	defer tx.Rollback()
	member, err := scanMember(tx.QueryRowContext(ctx, `
INSERT INTO nav_memberships (id, user_key, joined_at, last_seen_at, updated_at)
VALUES ($1::uuid, $2, $3, $3, $3)
ON CONFLICT (user_key) DO NOTHING
RETURNING id::text, user_key, status, display_name_snapshot, handle_snapshot,
          avatar_media_key_snapshot, profile_synced_at, join_reconciled_at, joined_at, last_seen_at,
          suspended_at, suspended_by_user_key, suspension_reason, updated_at, 0`,
		identifier.MustNew().String(), userKey, now))
	created := err == nil
	if errors.Is(err, sql.ErrNoRows) {
		member, err = scanMember(tx.QueryRowContext(ctx, `
UPDATE nav_memberships
SET last_seen_at = $2, updated_at = $2
WHERE user_key = $1 AND last_seen_at <= $3
RETURNING id::text, user_key, status, display_name_snapshot, handle_snapshot,
          avatar_media_key_snapshot, profile_synced_at, join_reconciled_at, joined_at, last_seen_at,
          suspended_at, suspended_by_user_key, suspension_reason, updated_at, 0`,
			userKey, now, now.Add(-service.touchInterval)))
		if errors.Is(err, sql.ErrNoRows) {
			member, err = scanMember(tx.QueryRowContext(ctx, memberSelect+` WHERE membership.user_key = $1`, userKey))
		}
	}
	if err != nil {
		return EnsureResult{}, fmt.Errorf("ensure membership: %w", err)
	}
	if created {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO nav_membership_events (id, membership_id, event_type, next_status, occurred_at)
VALUES ($1::uuid, $2::uuid, 'joined', 'active', $3)`,
			identifier.MustNew().String(), member.ID, now); err != nil {
			return EnsureResult{}, fmt.Errorf("record membership join: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return EnsureResult{}, fmt.Errorf("commit ensure membership: %w", err)
	}
	if created {
		member = service.refreshMemberProfile(ctx, member)
	}
	return EnsureResult{
		Member: member, Created: created, NeedsJoinReconcile: member.JoinReconciledAt.IsZero(),
	}, nil
}

func (service *Service) MarkJoinReconciled(ctx context.Context, userKey string) error {
	if !navidentity.IsPublicUserKey(userKey) {
		return ErrInvalidUserKey
	}
	result, err := service.db.ExecContext(ctx, `
UPDATE nav_memberships
SET join_reconciled_at = COALESCE(join_reconciled_at, $2), updated_at = $2
WHERE user_key = $1`, userKey, service.now())
	if err != nil {
		return fmt.Errorf("mark membership join reconciled: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read membership reconciliation result: %w", err)
	}
	if changed == 0 {
		return ErrNotFound
	}
	return nil
}

func (service *Service) Get(ctx context.Context, userKey string) (Member, error) {
	if !navidentity.IsPublicUserKey(userKey) {
		return Member{}, ErrInvalidUserKey
	}
	member, err := scanMember(service.db.QueryRowContext(ctx, memberSelect+` WHERE membership.user_key = $1`, userKey))
	if errors.Is(err, sql.ErrNoRows) {
		return Member{}, ErrNotFound
	}
	return member, err
}

func (service *Service) List(ctx context.Context, query Query) (Page, error) {
	if service == nil || service.db == nil {
		return Page{}, errors.New("navmember: service is not configured")
	}
	if query.Status != "" && query.Status != StatusActive && query.Status != StatusSuspended {
		return Page{}, ErrInvalidStatus
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Size < 1 || query.Size > 100 {
		query.Size = 20
	}
	where, args := memberWhere(query, false)
	var counts Counts
	if err := service.db.QueryRowContext(ctx, `
SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'active'),
       COUNT(*) FILTER (WHERE status = 'suspended')
FROM nav_memberships membership`+where, args...).Scan(&counts.All, &counts.Active, &counts.Suspended); err != nil {
		return Page{}, fmt.Errorf("count memberships: %w", err)
	}
	selectedTotal := counts.All
	if query.Status == StatusActive {
		selectedTotal = counts.Active
	} else if query.Status == StatusSuspended {
		selectedTotal = counts.Suspended
	}
	where, args = memberWhere(query, true)
	args = append(args, query.Size, (query.Page-1)*query.Size)
	rows, err := service.db.QueryContext(ctx, memberSelect+where+fmt.Sprintf(
		` ORDER BY membership.last_seen_at DESC, membership.id DESC LIMIT $%d OFFSET $%d`, len(args)-1, len(args)), args...)
	if err != nil {
		return Page{}, fmt.Errorf("list memberships: %w", err)
	}
	defer rows.Close()
	members := make([]Member, 0, query.Size)
	for rows.Next() {
		member, scanErr := scanMember(rows)
		if scanErr != nil {
			return Page{}, fmt.Errorf("scan membership: %w", scanErr)
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return Page{}, fmt.Errorf("iterate memberships: %w", err)
	}
	members = service.refreshStaleProfiles(ctx, members)
	return Page{Members: members, Counts: counts, Total: selectedTotal, Page: query.Page, Size: query.Size}, nil
}

func (service *Service) SetStatus(ctx context.Context, command SetStatusCommand) (Member, error) {
	if !navidentity.IsPublicUserKey(command.UserKey) || !navidentity.IsPublicUserKey(command.ActorUserKey) {
		return Member{}, ErrInvalidUserKey
	}
	if command.Status != StatusActive && command.Status != StatusSuspended {
		return Member{}, ErrInvalidStatus
	}
	command.Reason = strings.TrimSpace(command.Reason)
	if len([]rune(command.Reason)) > 500 {
		return Member{}, fmt.Errorf("navmember: reason exceeds 500 characters")
	}
	if command.Status == StatusSuspended && command.UserKey == command.ActorUserKey {
		return Member{}, ErrSelfSuspend
	}
	if command.Status == StatusSuspended && command.Reason == "" {
		return Member{}, ErrReasonRequired
	}
	now := service.now()
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return Member{}, fmt.Errorf("begin membership status change: %w", err)
	}
	defer tx.Rollback()
	current, err := scanMember(tx.QueryRowContext(ctx, memberSelect+` WHERE membership.user_key = $1 FOR UPDATE`, command.UserKey))
	if errors.Is(err, sql.ErrNoRows) {
		return Member{}, ErrNotFound
	}
	if err != nil {
		return Member{}, fmt.Errorf("load membership for status change: %w", err)
	}
	if current.Status == command.Status {
		if err := tx.Commit(); err != nil {
			return Member{}, err
		}
		return current, nil
	}
	suspendedAt := any(nil)
	suspendedBy := ""
	reason := ""
	if command.Status == StatusSuspended {
		suspendedAt = now
		suspendedBy = command.ActorUserKey
		reason = command.Reason
	}
	member, err := scanMember(tx.QueryRowContext(ctx, `
UPDATE nav_memberships
SET status = $2, suspended_at = $3, suspended_by_user_key = $4,
    suspension_reason = $5, updated_at = $6
WHERE user_key = $1
RETURNING id::text, user_key, status, display_name_snapshot, handle_snapshot,
          avatar_media_key_snapshot, profile_synced_at, join_reconciled_at, joined_at, last_seen_at,
          suspended_at, suspended_by_user_key, suspension_reason, updated_at,
          (SELECT COUNT(*) FROM nav_links WHERE submitter_sub = $1)`,
		command.UserKey, command.Status, suspendedAt, suspendedBy, reason, now))
	if err != nil {
		return Member{}, fmt.Errorf("update membership status: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO nav_membership_events (
    id, membership_id, actor_user_key, event_type, previous_status, next_status, reason, occurred_at
) VALUES ($1::uuid, $2::uuid, $3, 'status_changed', $4, $5, $6, $7)`,
		identifier.MustNew().String(), member.ID, command.ActorUserKey,
		current.Status, command.Status, command.Reason, now); err != nil {
		return Member{}, fmt.Errorf("record membership status change: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return Member{}, fmt.Errorf("commit membership status change: %w", err)
	}
	return member, nil
}

const memberSelect = `
SELECT membership.id::text, membership.user_key, membership.status,
       membership.display_name_snapshot, membership.handle_snapshot,
       membership.avatar_media_key_snapshot, membership.profile_synced_at,
       membership.join_reconciled_at,
       membership.joined_at, membership.last_seen_at, membership.suspended_at,
       membership.suspended_by_user_key, membership.suspension_reason, membership.updated_at,
       (SELECT COUNT(*) FROM nav_links link WHERE link.submitter_sub = membership.user_key)
FROM nav_memberships membership`

type rowScanner interface {
	Scan(...any) error
}

func scanMember(row rowScanner) (Member, error) {
	var member Member
	var profileSyncedAt, joinReconciledAt, suspendedAt sql.NullTime
	err := row.Scan(
		&member.ID, &member.UserKey, &member.Status, &member.DisplayName, &member.Handle,
		&member.AvatarMediaKey, &profileSyncedAt, &joinReconciledAt, &member.JoinedAt, &member.LastSeenAt,
		&suspendedAt, &member.SuspendedBy, &member.SuspensionReason, &member.UpdatedAt,
		&member.SubmissionCount,
	)
	if profileSyncedAt.Valid {
		member.ProfileSyncedAt = profileSyncedAt.Time
	}
	if joinReconciledAt.Valid {
		member.JoinReconciledAt = joinReconciledAt.Time
	}
	if suspendedAt.Valid {
		member.SuspendedAt = suspendedAt.Time
	}
	return member, err
}

func memberWhere(query Query, includeStatus bool) (string, []any) {
	conditions := make([]string, 0, 3)
	args := make([]any, 0, 3)
	if search := strings.TrimSpace(query.Search); search != "" {
		args = append(args, "%"+strings.ToLower(search)+"%")
		index := len(args)
		conditions = append(conditions, fmt.Sprintf(`(
lower(membership.user_key) LIKE $%d OR lower(membership.display_name_snapshot) LIKE $%d
OR lower(membership.handle_snapshot) LIKE $%d)`, index, index, index))
	}
	if query.ConstrainUserKeys {
		args = append(args, pq.Array(query.UserKeys))
		conditions = append(conditions, fmt.Sprintf("membership.user_key = ANY($%d)", len(args)))
	}
	if includeStatus && query.Status != "" {
		args = append(args, query.Status)
		conditions = append(conditions, fmt.Sprintf("membership.status = $%d", len(args)))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (service *Service) refreshStaleProfiles(ctx context.Context, members []Member) []Member {
	if service.profiles == nil || len(members) == 0 {
		return members
	}
	cutoff := service.now().Add(-service.profileRefreshInterval)
	keys := make([]string, 0, len(members))
	for _, member := range members {
		if member.ProfileSyncedAt.IsZero() || member.ProfileSyncedAt.Before(cutoff) {
			keys = append(keys, member.UserKey)
		}
	}
	profiles, err := service.profiles.GetMany(ctx, keys)
	if err != nil {
		return members
	}
	for index, member := range members {
		profile, ok := profiles[member.UserKey]
		if !ok {
			continue
		}
		members[index] = service.applyProfile(ctx, member, profile)
	}
	return members
}

func (service *Service) refreshMemberProfile(ctx context.Context, member Member) Member {
	if service.profiles == nil {
		return member
	}
	profiles, err := service.profiles.GetMany(ctx, []string{member.UserKey})
	if err != nil {
		return member
	}
	profile, ok := profiles[member.UserKey]
	if !ok {
		return member
	}
	return service.applyProfile(ctx, member, profile)
}

func (service *Service) applyProfile(ctx context.Context, member Member, profile identityclient.PublicUser) Member {
	syncedAt := service.now()
	avatar := ""
	if profile.Avatar != nil {
		avatar = profile.Avatar.MediaKey
	}
	if _, err := service.db.ExecContext(ctx, `
UPDATE nav_memberships
SET display_name_snapshot = $2, handle_snapshot = $3,
    avatar_media_key_snapshot = $4, profile_synced_at = $5, updated_at = $5
WHERE user_key = $1`, member.UserKey, profile.DisplayName, profile.Handle, avatar, syncedAt); err != nil {
		return member
	}
	member.DisplayName = profile.DisplayName
	member.Handle = profile.Handle
	member.AvatarMediaKey = avatar
	member.ProfileSyncedAt = syncedAt
	member.UpdatedAt = syncedAt
	return member
}
