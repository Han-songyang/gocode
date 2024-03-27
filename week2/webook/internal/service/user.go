package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"webook/internal/domain"
	"webook/internal/repository"
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
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
