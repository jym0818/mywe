package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Ctime   time.Time
	Utime   time.Time
	Author  Author
}
type Author struct {
	Id   int64
	Name string
}
