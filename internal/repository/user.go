package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jym0818/mywe/internal/domain"
	"github.com/jym0818/mywe/internal/repository/cache"
	"github.com/jym0818/mywe/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func (u *userRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	//先查缓存
	user, err := u.cache.Get(ctx, uid)
	if err == nil {
		return user, nil
	}
	//查找数据库
	ue, err := u.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	user = u.toDomain(ue)
	//回写缓存 可以开一个goroutine
	go func() {
		er := u.cache.Set(ctx, user)
		if er != nil {
			//记录日志
		}
	}()
	return user, nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return u.toDomain(user), nil
}

func (u *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := u.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return u.toDomain(user), nil
}

func NewuserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}
func (u *userRepository) Create(ctx context.Context, user domain.User) error {
	return u.dao.Insert(ctx, u.toEntity(user))
}

func (u *userRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Email:    sql.NullString{String: user.Email, Valid: user.Email != ""},
		Phone:    sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Password: user.Password,
		Ctime:    user.Ctime.UnixMilli(),
		Utime:    user.Utime.UnixMilli(),
		Nickname: user.Nickname,
	}
}
func (u *userRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		Ctime:    time.UnixMilli(user.Ctime),
		Utime:    time.UnixMilli(user.Utime),
		Nickname: user.Nickname,
	}
}
