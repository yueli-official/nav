// Package navauthz declares Navigation's consumer-owned authorization model.
// It intentionally differs from Docs: four scope levels, submitter relations,
// moderation, and scoped delegated administration.
package navauthz

import (
	"context"

	"github.com/yueli-official/foundation/go/authorization"
)

const (
	RootScopeID authorization.ScopeID = "nav"

	ScopeSite     authorization.ScopeType = "site"
	ScopeCategory authorization.ScopeType = "category"
	ScopeGroup    authorization.ScopeType = "nav_group"
	ScopeLink     authorization.ScopeType = "link"

	RoleAdministrator authorization.RoleKey = "administrator"
	RoleCurator       authorization.RoleKey = "curator"

	CapabilityLinkRead        authorization.CapabilityKey = "nav.link.read"
	CapabilityLinkSubmit      authorization.CapabilityKey = "nav.link.submit"
	CapabilityLinkUpdate      authorization.CapabilityKey = "nav.link.update"
	CapabilityLinkModerate    authorization.CapabilityKey = "nav.link.moderate"
	CapabilityStructureManage authorization.CapabilityKey = "nav.structure.manage"
	CapabilityHealthCheckRun  authorization.CapabilityKey = "nav.health_check.run"
	CapabilitySettingsManage  authorization.CapabilityKey = "nav.settings.manage"

	RelationSubmitter authorization.RelationKind = "submitter"

	ConstraintNormalRoleOwnsLink authorization.ConstraintKey = "nav.normal_role_owns_link"
)

func CategoryScopeID(id string) authorization.ScopeID {
	return authorization.ScopeID("category:" + id)
}

func GroupScopeID(id string) authorization.ScopeID {
	return authorization.ScopeID("group:" + id)
}

func LinkScopeID(id string) authorization.ScopeID {
	return authorization.ScopeID("link:" + id)
}

func Definition() authorization.Definition {
	normalSubjects := []authorization.SubjectKind{
		authorization.SubjectUser, authorization.SubjectService,
	}
	return authorization.Definition{
		Consumer: "nav",
		Version:  1,
		Capabilities: []authorization.CapabilityDefinition{
			{
				Key: CapabilityLinkRead, Version: 1, Binding: authorization.BindingAccessLayerEligible,
				AllowedScopes: []authorization.ScopeType{ScopeSite, ScopeCategory, ScopeGroup, ScopeLink},
			},
			{
				Key: CapabilityLinkSubmit, Version: 1, Binding: authorization.BindingNormal,
				AllowedScopes:    []authorization.ScopeType{ScopeGroup},
				EligibleSubjects: normalSubjects, Delegable: true,
			},
			{
				Key: CapabilityLinkUpdate, Version: 1, Binding: authorization.BindingNormal,
				AllowedScopes:    []authorization.ScopeType{ScopeSite, ScopeCategory, ScopeGroup, ScopeLink},
				EligibleSubjects: normalSubjects, QueryableRelation: RelationSubmitter, Delegable: true,
			},
			{
				Key: CapabilityLinkModerate, Version: 1, Binding: authorization.BindingNormal,
				AllowedScopes:    []authorization.ScopeType{ScopeCategory, ScopeGroup, ScopeLink},
				EligibleSubjects: normalSubjects, Delegable: true,
			},
			protectedCapability(CapabilityStructureManage, []authorization.ScopeType{ScopeSite, ScopeCategory}),
			protectedCapability(CapabilityHealthCheckRun, []authorization.ScopeType{ScopeSite}),
			protectedCapability(CapabilitySettingsManage, []authorization.ScopeType{ScopeSite}),
		},
		Scopes: authorization.ScopeSchema{Types: []authorization.ScopeTypeDefinition{
			{Key: ScopeSite, Root: true, Children: []authorization.ScopeType{ScopeCategory}},
			{Key: ScopeCategory, Children: []authorization.ScopeType{ScopeGroup}},
			{Key: ScopeGroup, Children: []authorization.ScopeType{ScopeLink}},
			{Key: ScopeLink},
		}},
		AccessLayers: []authorization.AccessLayerDefinition{
			{Key: authorization.AccessLayerVisitor, Capabilities: []authorization.CapabilityKey{CapabilityLinkRead}},
			{
				Key: authorization.AccessLayerAuthenticated,
				Capabilities: []authorization.CapabilityKey{
					CapabilityLinkRead,
					authorization.CapabilityApplicationCreate,
					authorization.CapabilityApplicationReadOwn,
					authorization.CapabilityApplicationWithdraw,
					authorization.CapabilityInvitationAccept,
				},
			},
		},
		Roles: []authorization.RoleDefinition{
			{
				Key: RoleAdministrator, DisplayName: "管理员", Protected: true,
				Capabilities: []authorization.CapabilityKey{
					authorization.CapabilityManage, authorization.CapabilityAuditRead,
					CapabilityLinkSubmit, CapabilityLinkUpdate,
					CapabilityLinkModerate, CapabilityStructureManage,
					CapabilityHealthCheckRun, CapabilitySettingsManage,
				},
			},
			{
				Key: RoleCurator, DisplayName: "内容维护者",
				Capabilities: []authorization.CapabilityKey{
					CapabilityLinkSubmit, CapabilityLinkUpdate,
				},
				Assignment: authorization.AssignmentPolicy{Sources: []authorization.GrantSource{
					authorization.GrantSourceApplication, authorization.GrantSourceInvitation,
					authorization.GrantSourceDirect, authorization.GrantSourceGroup,
				}},
			},
		},
		Constraints: []authorization.ConstraintDefinition{{
			Key: ConstraintNormalRoleOwnsLink, Version: 1, Mode: authorization.ConstraintSource,
			Capabilities: []authorization.CapabilityKey{CapabilityLinkUpdate}, AllNormalRoles: true,
		}},
	}
}

func ConstraintEvaluators() map[authorization.ConstraintKey]authorization.ConstraintEvaluator {
	return map[authorization.ConstraintKey]authorization.ConstraintEvaluator{
		ConstraintNormalRoleOwnsLink: authorization.ConstraintFunc(
			func(_ context.Context, input authorization.ConstraintInput) authorization.ConstraintResult {
				for _, submitter := range input.Resource.Relations[RelationSubmitter] {
					if submitter == input.Subject {
						return authorization.ConstraintResult{}
					}
				}
				return authorization.ConstraintResult{Denied: true}
			},
		),
	}
}

func protectedCapability(
	key authorization.CapabilityKey,
	scopes []authorization.ScopeType,
) authorization.CapabilityDefinition {
	return authorization.CapabilityDefinition{
		Key: key, Version: 1, Binding: authorization.BindingProtectedOnly,
		Risk: authorization.RiskHigh, Audit: authorization.AuditFull,
		AllowedScopes: scopes,
		EligibleSubjects: []authorization.SubjectKind{
			authorization.SubjectUser, authorization.SubjectService,
		},
	}
}
