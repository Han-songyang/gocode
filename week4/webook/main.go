package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/service/sms"
	"webook/internal/service/sms/localsms"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	server := InitWebServer()

	_ = server.Run(":8081")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&dao.User{})
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	//redisClind := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//server.Use(ratelimit.NewBuilder(redisClind, time.Second, 1).Build())
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		// 这个是允许前端访问你的后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))
	//useSession(server)
	useJWT(server)
	return server

}

func useJWT(server *gin.Engine) {
	login := middleware.LoginJWTMiddlewareBuild{}
	server.Use(login.CheckLogin())
}

//func useSession(server *gin.Engine) {
//	login := middleware.LoginMiddlewareBuild{}
//	store, err := redis.NewStore(16, "tcp", "localhost:6379",
//		"", []byte("xtpTacFeR4oDNWap"),
//		[]byte("cfR5BdYotg7n8QOM"))
//	if err != nil {
//		panic(err)
//	}
//	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
//}

func initUser(server *gin.Engine,
	redisClint redis.Cmdable,
	codeSvc service.CodeService,
	db *gorm.DB) {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(redisClint)
	ur := repository.NewCachedUserRepository(ud, uc)
	us := service.NewUserService(ur)
	c := web.NewUserHandler(us, codeSvc)
	c.Register(server)
}

func initCodeSvc(redisClint redis.Cmdable) service.CodeService {
	cc := cache.NewRedisCodeCache(redisClint)
	cr := repository.NewCodeRepository(cc)
	return service.NewCodeService(cr, initMemorySms())
}

func initMemorySms() sms.Service {
	return localsms.NewService()
}
