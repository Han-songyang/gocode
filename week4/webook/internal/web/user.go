package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

var JWTKey = []byte("xtpTacFeR4oDNWa7")

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

	ug.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)

	ug.POST("/edit", u.edit)
	ug.POST("/profile", u.ProfileSess)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	var user domain.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"bizerror": err.Error()})
		return
	}
	if err := u.service.Signup(ctx, user); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"bizerror": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var user domain.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	us, err := u.service.Login(ctx, user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": err.Error()})
		return
	}
	sses := sessions.Default(ctx)
	fmt.Println(us.Id)
	sses.Set("userId", us.Id)
	sses.Options(sessions.Options{
		MaxAge: 5,
	})
	sses.Save()
	ctx.JSON(http.StatusOK, gin.H{"message": "登陆成功"})
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	var user domain.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	us, err := u.service.Login(ctx, user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": err.Error()})
		return
	}
	uc := UserClaims{
		Uid: us.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": err.Error()})
		return
	}
	ctx.Header("x-jwt-token", tokenStr)

	ctx.JSON(http.StatusOK, gin.H{"message": "登陆成功"})
}

func (u *UserHandler) edit(ctx *gin.Context) {
	type Req struct {
		Id       int64  `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "生日格式错误", "code": err.Error()})
	}

	err = u.service.UpdateUserInfo(ctx, domain.User{
		Id:       req.Id,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Birthday: birthday.String(),
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "信息修改失败", "code": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "信息修改成功"})
}

func (u *UserHandler) ProfileSess(ctx *gin.Context) {
	uc, _ := ctx.MustGet("user").(UserClaims)
	user, err := u.service.FindById(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	type User struct {
		Nickname string
		Email    string
		AboutMe  string
		Birthday string
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"user": User{
			Nickname: user.Nickname,
			Email:    user.Email,
			AboutMe:  user.AboutMe,
			Birthday: user.Birthday,
		}})
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}
