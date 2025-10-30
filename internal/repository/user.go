package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jym0818/mywe/internal/domain"
	"github.com/jym0818/mywe/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (err error)
}

type userRepository struct {
	dao dao.UserDAO
}

func NewuserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}
func (u *userRepository) Create(ctx context.Context, user domain.User) (err error) {
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
