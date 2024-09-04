package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"webook/internal/domain"
	"webook/internal/repository"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	UpdateUserInfo(ctx *gin.Context, u domain.User) error
	FindById(ctx *gin.Context, uid int64) (domain.User, error)
	FindOrCreate(ctx *gin.Context, phone string) (domain.User, error)
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
		//logger: zap.L(),
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) UpdateUserInfo(ctx *gin.Context, u domain.User) error {
	err := svc.repo.UpdateUserInfo(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (svc *userService) FindById(ctx *gin.Context, uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid)
}

func (svc *userService) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	// 是否已经注册
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return u, nil
	}

	// 没有注册，创建用户
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && !errors.Is(err, repository.ErrDuplicateUser) {
		return domain.User{}, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}
