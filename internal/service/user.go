package service

import (
	"context"
	"errors"

	"github.com/jym0818/mywe/internal/domain"
	"github.com/jym0818/mywe/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号或者密码错误")

type UserService interface {
	Signup(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewuserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
func (u *userService) Signup(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return u.repo.Create(ctx, user)
}
func (u *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}
