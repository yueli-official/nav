package v1

import "github.com/gogf/gf/v2/frame/g"

type RecordPageViewReq struct {
	Referrer   string `json:"referrer,omitempty"`
	g.Meta     `path:"/api/v1/traffic/views" method:"post" tags:"Traffic" summary:"Record a public page view"`
	EventID    string `json:"eventId"`
	Path       string `json:"path"`
	VisitorID  string `json:"visitorId,omitempty"`
	OccurredAt string `json:"occurredAt"`
}

type RecordPageViewRes struct {
	Counted bool `json:"counted"`
	Replay  bool `json:"replay"`
}

type GetTrafficOverviewReq struct {
	g.Meta `path:"/api/v1/admin/nav/traffic" method:"get" tags:"Admin" summary:"Read public page view statistics"`
}

type TrafficTotalsView struct {
	Views             int64 `json:"views"`
	UniqueVisitorDays int64 `json:"uniqueVisitorDays"`
}

type TrafficDayView struct {
	Day string `json:"day"`
	TrafficTotalsView
}

type GetTrafficOverviewRes struct {
	Total  TrafficTotalsView `json:"total"`
	Today  TrafficTotalsView `json:"today"`
	Series []TrafficDayView  `json:"series"`
}
