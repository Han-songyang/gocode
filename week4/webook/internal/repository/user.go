package repository

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"webook/internal/bizerror"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUserInfo(ctx *gin.Context, u domain.User) error
	FindById(ctx *gin.Context, uid int64) (domain.User, error)
}
type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(d dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao: d,
	}
}
func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if errors.Is(err, dao.ErrRecordNotFound) {
		return domain.User{}, bizerror.New(bizerror.IncorrectUserNameOrPassword, "用户名或密码错误", err.Error())
	}
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) UpdateUserInfo(ctx *gin.Context, u domain.User) error {
	err := repo.dao.UpdateById(ctx, repo.toEntity(u))
	if err != nil {
		return err
	}
	return nil
}

func (repo *CachedUserRepository) FindById(ctx *gin.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return du, nil
	}
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)
	err = repo.cache.Set(ctx, du)
	return du, err
}

func (repo *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		ID:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
	}
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}
