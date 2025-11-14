package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Utime   time.Time
	Ctime   time.Time
	Status  ArticleStatus
}

func (a Article) Abstract() string {
	str := []rune(a.Content)
	if len(str) > 128 {
		str = str[:128]
	}
	return string(str)
}

type Author struct {
	Id   int64
	Name string
}
type ArticleStatus uint8

const (
	// ArticleStatusUnknown 这是一个未知状态
	ArticleStatusUnknown ArticleStatus = iota
	// ArticleStatusUnpublished 未发表
	ArticleStatusUnpublished
	// ArticleStatusPublished 已发表
	ArticleStatusPublished
	// ArticleStatusPrivate 仅自己可见
	ArticleStatusPrivate
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}
