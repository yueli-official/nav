package controller

import (
	"context"
	"errors"
	"time"

	"github.com/yueli-official/foundation/go/authorization"
	"github.com/yueli-official/foundation/go/traffic"
	v1 "github.com/yueli-official/nav/api/api/v1"
	"github.com/yueli-official/nav/api/internal/navauthz"
	"github.com/yueli-official/nav/api/internal/naverr"
	"github.com/yueli-official/nav/api/internal/navtraffic"
)

type TrafficPublic struct {
	Service *navtraffic.Service
}

func (c *TrafficPublic) RecordPageView(ctx context.Context, req *v1.RecordPageViewReq) (*v1.RecordPageViewRes, error) {
	result, err := c.Service.Record(ctx, req.EventID, req.Path, req.VisitorID, req.OccurredAt)
	if err != nil {
		var cause *traffic.Error
		if errors.As(err, &cause) {
			switch cause.Kind {
			case traffic.ErrorInvalidInput:
				return nil, naverr.Validation(cause.Field, "invalid", nil)
			case traffic.ErrorConflict:
				return nil, naverr.Conflict(req.EventID)
			}
		}
		return nil, err
	}
	return &v1.RecordPageViewRes{Counted: result.Counted, Replay: result.Replay}, nil
}

type TrafficAdmin struct {
	Service *navtraffic.Service
}

func (c *TrafficAdmin) GetTrafficOverview(ctx context.Context, _ *v1.GetTrafficOverviewReq) (*v1.GetTrafficOverviewRes, error) {
	if err := requireCapability(ctx, navauthz.CapabilitySettingsManage, navauthz.RootScopeID, authorization.ResourceFacts{}); err != nil {
		return nil, err
	}
	report, err := c.Service.Overview(ctx, time.Now())
	if err != nil {
		return nil, err
	}
	days := make([]v1.TrafficDayView, 0, len(report.Series))
	for _, day := range report.Series {
		days = append(days, v1.TrafficDayView{Day: day.Day.String(), TrafficTotalsView: trafficTotals(day.Totals)})
	}
	return &v1.GetTrafficOverviewRes{
		Total: trafficTotals(report.Total), Today: trafficTotals(report.Today), Series: days,
	}, nil
}

func trafficTotals(totals traffic.Totals) v1.TrafficTotalsView {
	return v1.TrafficTotalsView{Views: totals.Views, UniqueVisitorDays: totals.UniqueVisitorDays}
}
