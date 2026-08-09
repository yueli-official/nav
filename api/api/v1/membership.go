package v1

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type MembershipStatusCountsView struct {
	All       int `json:"all"`
	Active    int `json:"active"`
	Suspended int `json:"suspended"`
}

type MemberRoleView struct {
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
	Source      string `json:"source"`
}

type MemberView struct {
	UserKey             string           `json:"userKey"`
	Status              string           `json:"status"`
	DisplayName         string           `json:"displayName"`
	Handle              string           `json:"handle"`
	AvatarMediaKey      string           `json:"avatarMediaKey,omitempty"`
	JoinedAt            time.Time        `json:"joinedAt"`
	LastSeenAt          time.Time        `json:"lastSeenAt"`
	SuspendedAt         time.Time        `json:"suspendedAt,omitempty"`
	SuspendedBy         string           `json:"suspendedBy,omitempty"`
	SuspensionReason    string           `json:"suspensionReason,omitempty"`
	SubmissionCount     int              `json:"submissionCount"`
	PendingApplications int              `json:"pendingApplications"`
	Roles               []MemberRoleView `json:"roles"`
}

type AdminListMembersReq struct {
	g.Meta `path:"/api/v1/admin/nav/members" method:"GET" tags:"Admin Nav" summary:"List navigation members"`
	Q      string `p:"q"`
	Status string `p:"status" v:"in:active,suspended"`
	Role   string `p:"role"`
	Page   int    `p:"page" d:"1" v:"min:1"`
	Size   int    `p:"size" d:"20" v:"min:1|max:100"`
}

type AdminListMembersRes struct {
	Members []MemberView               `json:"members"`
	Counts  MembershipStatusCountsView `json:"counts"`
	Roles   []AuthorizationRoleView    `json:"roles"`
	Total   int                        `json:"total"`
	Page    int                        `json:"page"`
	Size    int                        `json:"size"`
}

type AdminSetMemberStatusReq struct {
	g.Meta  `path:"/api/v1/admin/nav/members/{userKey}/status" method:"PATCH" tags:"Admin Nav" summary:"Set navigation membership status"`
	UserKey string `p:"userKey" v:"required|regex:^[1-9A-HJ-NP-Za-km-z]{8}$"`
	Status  string `json:"status" v:"required|in:active,suspended"`
	Reason  string `json:"reason" v:"length:0,500"`
}

type AdminSetMemberStatusRes struct {
	Member MemberView `json:"member"`
}
