//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	ijwt "webook/internal/web/jwt"

	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitDB, InitRedis, ioc.InitLogger(),

		// dao
		dao.NewUserDAO,

		// cache
		cache.NewRedisCodeCache, cache.NewUserCache,

		// repository
		repository.NewCachedUserRepository, repository.NewCodeRepository,

		// service
		ioc.InitSms, service.NewUserService, service.NewCodeService, ioc.InitWechatService,

		// handler
		web.NewUserHandler, web.NewOAuth2WechatHandler, ijwt.NewRedisJWTHandler,

		// web
		ioc.InitGinMiddlewares, ioc.InitWebServer,
	)
	return gin.Default()
}
