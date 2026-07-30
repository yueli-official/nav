package v1

import "github.com/gogf/gf/v2/frame/g"

type MeView struct {
	Sub             string   `json:"sub"`
	Authenticated   bool     `json:"authenticated"`
	IsAdministrator bool     `json:"isAdministrator"`
	Capabilities    []string `json:"capabilities"`
}

type MeReq struct {
	g.Meta `path:"/api/v1/me" method:"GET" tags:"Nav" summary:"Current caller access"`
}

type MeRes struct {
	Me *MeView `json:"me"`
}
