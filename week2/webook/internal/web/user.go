package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"webook/internal/domain"
	"webook/internal/service"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	service        service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		service:        service,
	}
}

func (u *UserHandler) Register(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login")
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	var user domain.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := u.service.Signup(ctx, user); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "success"})
}
