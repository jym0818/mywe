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
	Profile(ctx context.Context, uid int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func (u *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	//查找
	user, err := u.repo.FindByPhone(ctx, phone)
	if err == nil {
		return user, nil
	}
	//系统错误
	if err != repository.ErrUserNotFound {
		return domain.User{}, err
	}
	//不存在 就去注册
	err = u.repo.Create(ctx, domain.User{Phone: phone})
	//注册失败 返回错误
	if err != nil && err != repository.ErrUserDuplicateEmail {
		return domain.User{}, err
	}
	return u.repo.FindByPhone(ctx, phone)

}

func (u *userService) Profile(ctx context.Context, uid int64) (domain.User, error) {
	return u.repo.FindById(ctx, uid)
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
