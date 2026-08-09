package navauthz

import (
	"context"
	"database/sql"
	"fmt"

	foundationauth "github.com/yueli-official/foundation/go/auth"
	"github.com/yueli-official/foundation/go/authorization"
	"github.com/yueli-official/nav/api/internal/navidentity"
)

type Runtime interface {
	authorization.Authorizer
	authorization.QueryPlanner
	authorization.AccessReader
	authorization.ResourceScopeRegistry
	authorization.ResourceScopeRelocator
	authorization.RoleManager
	authorization.RoleReader
	authorization.GrantReader
	authorization.WorkflowManager
	authorization.WorkflowReader
	authorization.PolicyManager
	authorization.PolicyReader
	authorization.Reconciler
}

type Service struct {
	runtime Runtime
	db      *sql.DB
}

func New(runtime Runtime, db *sql.DB) *Service {
	return &Service{runtime: runtime, db: db}
}

func (service *Service) Runtime() Runtime {
	if service == nil {
		return nil
	}
	return service.runtime
}

func (service *Service) Subject(ctx context.Context) authorization.SubjectRef {
	principal, ok := foundationauth.FromContext(ctx)
	if !ok {
		return authorization.SubjectRef{Kind: authorization.SubjectAnonymous}
	}
	subjectKind, _ := principal.Claim("subject_kind")
	if subjectKind == "user" {
		if userKey, ok := navidentity.PublicUserKey(principal); ok {
			return authorization.SubjectRef{Kind: authorization.SubjectUser, ID: userKey}
		}
	}
	if subjectKind == "client" && principal.ClientID != "" {
		return authorization.SubjectRef{Kind: authorization.SubjectService, ID: principal.ClientID}
	}
	return authorization.SubjectRef{Kind: authorization.SubjectAnonymous}
}

func (service *Service) Decide(
	ctx context.Context,
	capability authorization.CapabilityKey,
	scopeID authorization.ScopeID,
	resource authorization.ResourceFacts,
) (authorization.Decision, error) {
	if service == nil || service.runtime == nil {
		return authorization.Decision{}, &authorization.Error{
			Kind: authorization.ErrorUnavailable, Field: "runtime", Message: "is not configured",
		}
	}
	return service.runtime.Decide(ctx, authorization.DecisionRequest{
		Subject: service.Subject(ctx), Capability: capability, ScopeID: scopeID, Resource: resource,
	})
}

func (service *Service) EffectiveAccess(ctx context.Context) (authorization.EffectiveAccess, error) {
	if service == nil || service.runtime == nil {
		return authorization.EffectiveAccess{}, &authorization.Error{
			Kind: authorization.ErrorUnavailable, Field: "runtime", Message: "is not configured",
		}
	}
	return service.runtime.EffectiveAccess(ctx, authorization.EffectiveAccessQuery{
		Subject: service.Subject(ctx), ScopeID: RootScopeID,
	})
}

// ReconcileNewMember evaluates automatic grants exactly once for the product
// membership join event. Permission reads must never replay this lifecycle event.
func (service *Service) ReconcileNewMember(ctx context.Context) error {
	if service == nil || service.runtime == nil {
		return &authorization.Error{
			Kind: authorization.ErrorUnavailable, Field: "runtime", Message: "is not configured",
		}
	}
	subject := service.Subject(ctx)
	if subject.Kind != authorization.SubjectUser || subject.ID == "" {
		return nil
	}
	preview, err := service.runtime.PreviewReconcileSubject(ctx, authorization.ReconcileSubjectCommand{
		Subject: subject,
	})
	if err != nil || preview.Created == 0 {
		return err
	}
	_, err = service.runtime.ReconcileSubject(ctx, authorization.ReconcileSubjectCommand{
		Subject: subject,
	})
	return err
}

func (service *Service) IsAdministrator(ctx context.Context) bool {
	decision, err := service.Decide(ctx, authorization.CapabilityManage, RootScopeID, authorization.ResourceFacts{})
	return err == nil && decision.Allowed
}

func (service *Service) EnsureCategoryScope(ctx context.Context, categoryID string) error {
	return ensureScope(ctx, service.runtime, authorization.RegisterScopeCommand{
		ID: CategoryScopeID(categoryID), Type: ScopeCategory, ParentID: RootScopeID,
	})
}

func (service *Service) EnsureGroupScope(ctx context.Context, groupID, categoryID string) error {
	if err := service.EnsureCategoryScope(ctx, categoryID); err != nil {
		return err
	}
	return ensureScope(ctx, service.runtime, authorization.RegisterScopeCommand{
		ID: GroupScopeID(groupID), Type: ScopeGroup, ParentID: CategoryScopeID(categoryID),
	})
}

func (service *Service) EnsureLinkScope(ctx context.Context, linkID, groupID, categoryID string) error {
	if err := service.EnsureGroupScope(ctx, groupID, categoryID); err != nil {
		return err
	}
	return ensureScope(ctx, service.runtime, authorization.RegisterScopeCommand{
		ID: LinkScopeID(linkID), Type: ScopeLink, ParentID: GroupScopeID(groupID),
	})
}

func (service *Service) ReparentGroupScope(ctx context.Context, groupID, categoryID string) error {
	if err := service.EnsureCategoryScope(ctx, categoryID); err != nil {
		return err
	}
	_, err := service.runtime.ReparentScope(ctx, authorization.ReparentScopeCommand{
		ID: GroupScopeID(groupID), ParentID: CategoryScopeID(categoryID),
	})
	return err
}

func (service *Service) ReparentLinkScope(ctx context.Context, linkID, groupID, categoryID string) error {
	if err := service.EnsureGroupScope(ctx, groupID, categoryID); err != nil {
		return err
	}
	_, err := service.runtime.ReparentScope(ctx, authorization.ReparentScopeCommand{
		ID: LinkScopeID(linkID), ParentID: GroupScopeID(groupID),
	})
	return err
}

func LinkResource(id, submitterSub string) authorization.ResourceFacts {
	resource := authorization.ResourceFacts{
		Type: "link", ID: authorization.ResourceID(id), ScopeID: LinkScopeID(id),
	}
	if submitterSub != "" {
		resource.Relations = map[authorization.RelationKind][]authorization.SubjectRef{
			RelationSubmitter: {{Kind: authorization.SubjectUser, ID: submitterSub}},
		}
	}
	return resource
}

type ManageLinkAccess struct {
	All                bool
	OwnedAll           bool
	SubjectID          string
	UnrestrictedGroups []string
	OwnedGroups        []string
}

// ManageLinkFilter translates Nav's relation-constrained collection access.
// Descendant grants are checked per group so scoped curators cannot see links
// outside the scope delegated to them.
func (service *Service) ManageLinkFilter(ctx context.Context) (ManageLinkAccess, error) {
	if service == nil || service.runtime == nil || service.db == nil {
		return ManageLinkAccess{}, &authorization.Error{
			Kind: authorization.ErrorUnavailable, Field: "runtime", Message: "is not configured",
		}
	}
	subject := service.Subject(ctx)
	rootPlan, err := service.runtime.Plan(ctx, authorization.QueryRequest{
		Subject: subject, Capability: CapabilityLinkUpdate, ScopeID: RootScopeID,
	})
	if err != nil {
		return ManageLinkAccess{}, err
	}
	switch rootPlan.Kind {
	case authorization.QueryAll:
		return ManageLinkAccess{All: true}, nil
	case authorization.QueryRelation:
		if rootPlan.Relation != RelationSubmitter {
			return ManageLinkAccess{}, fmt.Errorf("unsupported nav query relation %q", rootPlan.Relation)
		}
		return ManageLinkAccess{OwnedAll: true, SubjectID: string(subject.ID)}, nil
	}

	rows, err := service.db.QueryContext(ctx, `SELECT id::text FROM nav_groups ORDER BY id`)
	if err != nil {
		return ManageLinkAccess{}, fmt.Errorf("list groups for authorization query: %w", err)
	}
	defer rows.Close()
	access := ManageLinkAccess{SubjectID: string(subject.ID)}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ManageLinkAccess{}, fmt.Errorf("scan group for authorization query: %w", err)
		}
		plan, planErr := service.runtime.Plan(ctx, authorization.QueryRequest{
			Subject: subject, Capability: CapabilityLinkUpdate, ScopeID: GroupScopeID(id),
		})
		if planErr != nil {
			return ManageLinkAccess{}, planErr
		}
		switch plan.Kind {
		case authorization.QueryAll:
			access.UnrestrictedGroups = append(access.UnrestrictedGroups, id)
		case authorization.QueryRelation:
			if plan.Relation == RelationSubmitter {
				access.OwnedGroups = append(access.OwnedGroups, id)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return ManageLinkAccess{}, fmt.Errorf("iterate groups for authorization query: %w", err)
	}
	return access, nil
}

func (service *Service) EffectiveManagementAccess(ctx context.Context) ([]authorization.CapabilityKey, error) {
	seen := map[authorization.CapabilityKey]struct{}{}
	add := func(access authorization.EffectiveAccess) {
		for _, capability := range access.Capabilities {
			seen[capability] = struct{}{}
		}
	}
	root, err := service.EffectiveAccess(ctx)
	if err != nil {
		return nil, err
	}
	add(root)
	if service.db != nil {
		rows, queryErr := service.db.QueryContext(ctx, `SELECT id::text FROM nav_groups ORDER BY id`)
		if queryErr != nil {
			return nil, queryErr
		}
		defer rows.Close()
		for rows.Next() {
			var id string
			if scanErr := rows.Scan(&id); scanErr != nil {
				return nil, scanErr
			}
			access, accessErr := service.runtime.EffectiveAccess(ctx, authorization.EffectiveAccessQuery{
				Subject: service.Subject(ctx), ScopeID: GroupScopeID(id),
			})
			if accessErr != nil {
				return nil, accessErr
			}
			add(access)
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			return nil, rowsErr
		}
	}
	result := make([]authorization.CapabilityKey, 0, len(seen))
	for capability := range seen {
		result = append(result, capability)
	}
	return result, nil
}
