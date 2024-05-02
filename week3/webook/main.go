package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middleware/ratelimit"
)

func main() {
	db := initDB()
	server := initWebServer()
	initUser(server, db)
	_ = server.Run(":8081")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
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
	redisClind := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClind, time.Second, 1).Build())
	//server.Use(cors.New(cors.Config{
	//	AllowCredentials: true,
	//	AllowHeaders:     []string{"Content-Type"},
	//	AllowOriginFunc: func(origin string) bool {
	//		if strings.HasPrefix(origin, "http://localhost") {
	//			return strings.Contains(origin, "localhost")
	//		}
	//		return false
	//	},
	//	MaxAge:           12 * time.Hour,
	//}))
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

func initUser(server *gin.Engine, db *gorm.DB) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewCachedUserRepository(ud)
	us := service.NewUserService(ur)
	c := web.NewUserHandler(us)
	c.Register(server)
}
