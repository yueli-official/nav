package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/yueli-official/nav/api/internal/naverr"
)

func CauseMappingMiddleware(request *ghttp.Request) {
	request.Middleware.Next()
	if err := request.GetError(); err != nil {
		request.SetError(naverr.MapCause(err))
	}
}
