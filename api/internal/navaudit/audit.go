package navaudit

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/yueli-official/foundation/go/audit"
	foundationauth "github.com/yueli-official/foundation/go/auth"
	"github.com/yueli-official/foundation/go/authorization"
	"go.opentelemetry.io/otel/trace"
)

type Action string

const (
	ActionNavigationPublished  Action = "nav.navigation.published"
	ActionNavigationArchived   Action = "nav.navigation.archived"
	ActionNavigationDeleted    Action = "nav.navigation.deleted"
	ActionTaxonomyChanged      Action = "nav.taxonomy.changed"
	ActionSiteProfilePublished Action = "nav.site_profile.published"
	ActionDataExported         Action = "nav.data.exported"
)

type Evidence struct {
	Revision uint64
	Digest   string
	Count    uint64
}

type Journal struct {
	core      *audit.Postgres
	contracts map[Action]audit.Contract[Evidence]
}

func New(ctx context.Context, db *sql.DB, instance string) (*Journal, error) {
	catalog, err := audit.Compile(Definition())
	if err != nil {
		return nil, err
	}
	core, err := audit.NewPostgres(ctx, catalog, audit.PostgresOptions{
		DB: db, InstanceKey: "nav:" + instance,
		Source:             audit.Source{Service: "nav-api", Module: "navigation", Instance: instance},
		EnableMirrorOutbox: true,
	})
	if err != nil {
		return nil, err
	}
	journal := &Journal{core: core, contracts: make(map[Action]audit.Contract[Evidence])}
	for _, action := range actions() {
		contract, err := audit.BindAction(catalog, audit.Action{Name: audit.ActionName(action), Version: 1}, encodeEvidence)
		if err != nil {
			return nil, err
		}
		journal.contracts[action] = contract
	}
	return journal, nil
}

func Definition() audit.Definition {
	definitions := make([]audit.ActionDefinition, 0, len(actions()))
	for _, action := range actions() {
		category := audit.CategoryAdministration
		retention := audit.RetentionClass("retention.nav_management")
		targets := []string{"nav.link", "nav.link_batch", "nav.taxonomy", "nav.site_profile"}
		if action == ActionDataExported {
			category = audit.CategoryExport
			retention = "retention.nav_export"
			targets = []string{"nav.data_export"}
		}
		commit := audit.CommitAtomicRequired
		if action == ActionDataExported {
			commit = audit.CommitIndependentAllow
		}
		definitions = append(definitions, audit.ActionDefinition{
			Action:   audit.Action{Name: audit.ActionName(action), Version: 1},
			Category: category, TargetTypes: targets,
			Commit: commit, Retention: retention,
			Evidence: []audit.FieldDefinition{
				{Key: "nav.revision", Kind: audit.EvidenceCount},
				{Key: "nav.digest", Kind: audit.EvidenceDigest},
				{Key: "nav.count", Kind: audit.EvidenceCount},
			},
		})
	}
	return audit.Definition{
		Version: 1, Consumer: "nav.audit", MaxBatch: 100, MaxEvidence: 8,
		Retention: []audit.RetentionDefinition{
			{Class: "retention.nav_management", MinimumAge: 365 * 24 * time.Hour, ArchiveBefore: true},
			{Class: "retention.nav_export", MinimumAge: 365 * 24 * time.Hour, ArchiveBefore: true},
		},
		Actions: definitions,
	}
}

func actions() []Action {
	return []Action{
		ActionNavigationPublished, ActionNavigationArchived, ActionNavigationDeleted,
		ActionTaxonomyChanged, ActionSiteProfilePublished, ActionDataExported,
	}
}

func encodeEvidence(value Evidence) []audit.EvidenceField {
	var fields []audit.EvidenceField
	if value.Revision != 0 {
		fields = append(fields, audit.Count("nav.revision", value.Revision))
	}
	if value.Digest != "" {
		fields = append(fields, audit.EvidenceDigestValue("nav.digest", value.Digest))
	}
	if value.Count != 0 {
		fields = append(fields, audit.Count("nav.count", value.Count))
	}
	return fields
}

func (journal *Journal) Hook(
	ctx context.Context,
	action Action,
	eventID string,
	target audit.Target,
	evidence Evidence,
) func(context.Context, *sql.Tx) error {
	if journal == nil {
		return nil
	}
	actor := actorFromContext(ctx)
	correlation := correlationFromContext(ctx)
	occurredAt := time.Now().UTC()
	return func(txCtx context.Context, tx *sql.Tx) error {
		appender, err := journal.core.Bind(tx)
		if err != nil {
			return err
		}
		_, err = audit.Record(txCtx, appender, journal.contracts[action], audit.Attempt[Evidence]{
			ID: audit.EventID(eventID), Actor: actor, Target: target,
			Outcome:     audit.Outcome{Kind: audit.OutcomeSucceeded},
			Correlation: correlation,
			OccurredAt:  occurredAt, Evidence: evidence,
		})
		return err
	}
}

func (journal *Journal) Reader() audit.Reader {
	if journal == nil {
		return nil
	}
	return journal.core
}

func (journal *Journal) Export(
	ctx context.Context,
	request audit.ExportRequest,
	writer io.Writer,
) (audit.ExportManifest, error) {
	manifest, err := journal.core.Export(ctx, request, writer)
	if err != nil {
		return audit.ExportManifest{}, err
	}
	command, err := audit.Prepare(journal.contracts[ActionDataExported], audit.Attempt[Evidence]{
		ID: audit.EventID(uuid.NewString()), Actor: actorFromContext(ctx),
		Target:      audit.Target{Type: "nav.data_export", ID: string(manifest.ContentDigest)},
		Outcome:     audit.Outcome{Kind: audit.OutcomeSucceeded},
		Correlation: correlationFromContext(ctx),
		Evidence:    Evidence{Count: manifest.Count, Digest: string(manifest.ContentDigest)},
	})
	if err != nil {
		return audit.ExportManifest{}, err
	}
	if _, err := journal.core.AppendIndependent(ctx, command); err != nil {
		return audit.ExportManifest{}, err
	}
	return manifest, nil
}

func actorFromContext(ctx context.Context) audit.Actor {
	if principal, ok := foundationauth.FromContext(ctx); ok && principal != nil {
		kind, _ := principal.Claim("subject_kind")
		switch kind {
		case "user":
			return audit.Actor{Kind: audit.ActorUser, ID: principal.Subject}
		case "guest":
			return audit.Actor{Kind: audit.ActorGuest, ID: principal.Subject}
		case "client":
			return audit.Actor{Kind: audit.ActorService, ID: principal.ClientID}
		}
	}
	return audit.Actor{Kind: audit.ActorSystem, ID: "nav-api"}
}

func correlationFromContext(ctx context.Context) audit.Correlation {
	value := audit.Correlation{RequestID: authorization.RequestMetadataFromContext(ctx).CorrelationID}
	span := trace.SpanContextFromContext(ctx)
	if span.IsValid() {
		value.TraceID = span.TraceID().String()
		value.SpanID = span.SpanID().String()
	}
	return value
}
