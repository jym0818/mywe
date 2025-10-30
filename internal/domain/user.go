package domain

import "time"

type User struct {
	Id       int64
	Nickname string
	Email    string
	Phone    string
	Password string
	Ctime    time.Time
	Utime    time.Time
}
