package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
}
type CachedUserRepository struct {
	dao dao.UserDAO
}

func NewCachedUserRepository(dao dao.UserDAO) UserRepository {
	return &CachedUserRepository{
		dao: dao,
	}
}
func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Email:    u.Email,
		Password: u.Password,
	}
}
