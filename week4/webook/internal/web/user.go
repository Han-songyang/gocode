package web

import (
	"errors"
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
	codeSvc        service.CodeService
}

func NewUserHandler(service service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		service:        service,
		codeSvc:        codeSvc,
	}
}

func (u *UserHandler) Register(server *gin.Engine) {
	ug := server.Group("/users")

	ug.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)

	ug.POST("/edit", u.edit)
	ug.GET("/profile", u.ProfileSess)

	ug.POST("/login_sms", u.LoginSMS)
	ug.POST("/login_sms/code/send", u.SendSMSLoginCode)
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
		Birthday: birthday,
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
		Phone    string
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"user": User{
			Nickname: user.Nickname,
			Email:    user.Email,
			AboutMe:  user.AboutMe,
			Birthday: user.Birthday.Format(time.DateOnly),
			Phone:    user.Phone,
		}})
}

func (u *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你这边可以校验 Req
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	err := u.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不对，请重新输入",
		})
		return
	}
	user, err := u.service.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	u.setJWTToken(ctx, user.Id)
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 1 分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.Header("x-jwt-token", tokenStr)
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserAgent string
	Uid       int64
}
