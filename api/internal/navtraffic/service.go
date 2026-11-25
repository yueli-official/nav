package navtraffic

import (
	"context"
	"strings"
	"time"

	"github.com/yueli-official/foundation/go/traffic"
	"github.com/yueli-official/nav/api/internal/naverr"
)

const pageKind traffic.ResourceKind = "page"

type store interface {
	traffic.Recorder
	traffic.VisitorTokenizer
	traffic.Reader
}

type Service struct {
	store   store
	catalog *traffic.Catalog
}

type Overview struct {
	Total  traffic.Totals
	Today  traffic.Totals
	Series []traffic.SeriesPoint
}

func Catalog() *traffic.Catalog {
	return traffic.MustCompile(traffic.Definition{
		Version:       traffic.DefinitionVersion,
		TimeZone:      "Asia/Shanghai",
		ResourceKinds: []traffic.ResourceKindDefinition{{Key: pageKind}},
	})
}

func New(storage store, catalog *traffic.Catalog) *Service {
	return &Service{store: storage, catalog: catalog}
}

func (service *Service) Record(ctx context.Context, eventID, path, visitorID, timestamp string) (traffic.RecordResult, error) {
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	slug, isCategory := strings.CutPrefix(path, "/categories/")
	category := isCategory && slug != "" && !strings.ContainsAny(slug, "/\\?#")
	if (path != "/" && !category) || len(path) > 200 {
		return traffic.RecordResult{}, naverr.Validation("path", "invalid", nil)
	}
	occurredAt, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return traffic.RecordResult{}, naverr.Validation("occurredAt", "invalid", nil)
	}
	observation := traffic.Observation{
		EventID:    traffic.EventID(eventID),
		Resource:   traffic.Resource{Kind: pageKind, ID: path},
		OccurredAt: occurredAt,
		Class:      traffic.VisitUnknown,
	}
	if visitorID != "" {
		if len(visitorID) < 16 || len(visitorID) > 128 {
			return traffic.RecordResult{}, naverr.Validation("visitorId", "invalid", nil)
		}
		token, err := service.store.TokenizeVisitor(ctx, occurredAt, []byte("browser\x00"+visitorID))
		if err != nil {
			return traffic.RecordResult{}, err
		}
		observation.HasVisitor = true
		observation.VisitorToken = token
	}
	return service.store.Record(ctx, observation)
}

func (service *Service) Overview(ctx context.Context, now time.Time) (Overview, error) {
	today := service.catalog.DayAt(now)
	tomorrow := service.catalog.DayAt(now.Add(24 * time.Hour))
	recent := service.catalog.DayAt(now.Add(-13 * 24 * time.Hour))
	scope := traffic.InstanceScope()

	total, err := service.store.Summary(ctx, traffic.SummaryQuery{Scope: scope})
	if err != nil {
		return Overview{}, err
	}
	day, err := service.store.Summary(ctx, traffic.SummaryQuery{
		Scope: scope, Range: &traffic.DateRange{From: today, To: tomorrow},
	})
	if err != nil {
		return Overview{}, err
	}
	series, err := service.store.Series(ctx, traffic.SeriesQuery{
		Scope: scope, Range: traffic.DateRange{From: recent, To: tomorrow},
	})
	if err != nil {
		return Overview{}, err
	}
	return Overview{Total: total.Totals, Today: day.Totals, Series: series}, nil
}
