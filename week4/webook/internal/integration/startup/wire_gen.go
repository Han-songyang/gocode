// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/jwt"
	"webook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	logger := InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, logger)
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSms()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(handler, userService, codeService)
	wechatService := InitWechatService(logger)
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, handler, userService)
	articleDAO := dao.NewArticleGORMDAO(db)
	articleCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(articleDAO, userRepository, articleCache)
	articleService := service.NewArticleService(articleRepository)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, logger)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	articleHandler := web.NewArticleHandler(logger, articleService, interactiveService)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	return engine
}

func InitArticleHandler(dao2 dao.ArticleDAO) *web.ArticleHandler {
	logger := InitLogger()
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	cmdable := InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	articleCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(dao2, userRepository, articleCache)
	articleService := service.NewArticleService(articleRepository)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, logger)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	articleHandler := web.NewArticleHandler(logger, articleService, interactiveService)
	return articleHandler
}

func InitInteractiveService() service.InteractiveService {
	db := InitDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := InitRedis()
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	logger := InitLogger()
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, logger)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	return interactiveService
}

// wire.go:

var thirdPartySet = wire.NewSet(
	InitRedis, InitDB,
	InitLogger)

var userSvcProvider = wire.NewSet(dao.NewUserDAO, cache.NewUserCache, repository.NewCachedUserRepository, service.NewUserService)

var articlSvcProvider = wire.NewSet(repository.NewCachedArticleRepository, cache.NewArticleRedisCache, dao.NewArticleGORMDAO, service.NewArticleService)

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewInteractiveRedisCache, repository.NewCachedInteractiveRepository, service.NewInteractiveService)
