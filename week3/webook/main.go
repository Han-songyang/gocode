package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	db := initDB()
	server := initWebServer()
	initUser(server, db)
	_ = server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/webook"), &gorm.Config{})
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
	loginSession(server)
	return server
}

func loginSession(server *gin.Engine) {
	login := middleware.LoginMiddlewareBuild{}
	store, err := redis.NewStore(16, "tcp", "localhost:6379",
		"", []byte("xtpTacFeR4oDNWap"),
		[]byte("cfR5BdYotg7n8QOM"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}

func initUser(server *gin.Engine, db *gorm.DB) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewCachedUserRepository(ud)
	us := service.NewUserService(ur)
	c := web.NewUserHandler(us)
	c.Register(server)
}
