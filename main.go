package main

import (
	"io"
	"os"

	"nav/conf"
	"nav/models"

	"github.com/rs/zerolog/log"

	"github.com/Yuelioi/gkit/logx/zero"
	"github.com/Yuelioi/gkit/web/gin/middleware/apikey"
	"github.com/Yuelioi/gkit/web/gin/middleware/log/gzero"
	"github.com/Yuelioi/gkit/web/gin/middleware/ratelimit"
	"github.com/Yuelioi/gkit/web/gin/server"
	"github.com/gin-gonic/gin"
)

var resource []models.CategoryData
var configPath = "./data/config.yaml"

func main() {
	// 初始化日志
	logger := zero.Default()
	log.Logger = logger

	// 加载导航配置
	var err error
	resource, err = conf.LoadNavConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("资源加载失败")
	}

	key := conf.InitAPIKey()

	// 禁用默认 Gin 输出
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	// 服务器配置
	cfg := server.ServerConfig{
		Addr:      ":9000",
		Logger:    logger,
		Mode:      os.Getenv("APP_MODE"),
		APIPrefix: "/api/v1",
		Middlewares: []gin.HandlerFunc{
			gzero.Default(logger),
			gzero.GinRecovery(logger),
			ratelimit.Default(),
		},
		EnableCORS: true,
		SPAPath:    "./frontend/dist",
	}

	logger.Info().Msg("后台访问密码为:" + key)

	// 启动服务器
	err = server.Start(cfg, func(api *gin.RouterGroup) {
		// 公开路由
		public := api.Group("/")
		{
			public.GET("/navs", getNavs)

			public.PUT("/categories/:cid/groups/:gid/collections", updateCollection)

			public.POST("/categories/:cid/groups/:gid/collections", createCollection)

		}

		// 需要认证的路由
		auth := api.Group("/", apikey.Default(conf.IsValidAPIKey))
		{
			// 分类管理
			auth.POST("/categories", createCategory)
			auth.DELETE("/categories/:cid", deleteCategory)
			auth.PUT("/categories/:cid", updateCategory)

			// 分组管理
			auth.POST("/categories/:cid/groups", createGroup)
			auth.DELETE("/categories/:cid/groups/:gid", deleteGroup)
			auth.PUT("/categories/:cid/groups/:gid", updateGroup)

			// 收藏项管理
			auth.DELETE("/categories/:cid/groups/:gid/collections", deleteCollection)
		}
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("服务启动失败")
	}
}
