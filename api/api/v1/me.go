package v1

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type MeView struct {
	Sub             string            `json:"sub"`
	UserKey         string            `json:"userKey,omitempty"`
	Authenticated   bool              `json:"authenticated"`
	IsAdministrator bool              `json:"isAdministrator"`
	Capabilities    []string          `json:"capabilities"`
	Membership      *MeMembershipView `json:"membership,omitempty"`
}

type MeMembershipView struct {
	Status     string    `json:"status"`
	JoinedAt   time.Time `json:"joinedAt"`
	LastSeenAt time.Time `json:"lastSeenAt"`
}

type MeReq struct {
	g.Meta `path:"/api/v1/me" method:"GET" tags:"Nav" summary:"Current caller access"`
}

type MeRes struct {
	Me *MeView `json:"me"`
}
