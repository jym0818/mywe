package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrUserDuplicateEmail = errors.New("邮件冲突")

type UserDAO interface {
	Insert(ctx context.Context, user User) error
}
type userDAO struct {
	db *gorm.DB
}

func (u *userDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().Unix()
	user.Ctime = now
	user.Utime = now
	err := u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok {
			const uniqueIndexErrNo uint16 = 1062
			if me.Number == uniqueIndexErrNo {
				return ErrUserDuplicateEmail
			}
		}
	}
	return err
}

func NewuserDAO(db *gorm.DB) UserDAO {
	return &userDAO{db: db}
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password string
	Nickname string
	Ctime    int64
	Utime    int64
}
