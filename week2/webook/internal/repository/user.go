package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
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

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}
}
