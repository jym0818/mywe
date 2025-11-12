package dao

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art PublishedArticle) error
	SyncStatus(ctx context.Context, id int64, author int64, status uint8) error
}

type articleDAO struct {
	db *gorm.DB
}

func NewarticleDAO(db *gorm.DB) ArticleDAO {
	return &articleDAO{
		db: db,
	}
}
func (dao *articleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}
func (dao *articleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := dao.db.WithContext(ctx).Model(&art).Where("id = ? AND author_id = ?", art.Id, art.AuthorId).Updates(map[string]interface{}{
		"utime":   now,
		"content": art.Content,
		"title":   art.Title,
		"status":  art.Status,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，可能是创作者非法 id %d, author_id %d", art.Id, art.AuthorId)

	}
	return nil
}

func (dao *articleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		txDAO := NewarticleDAO(tx)
		if id > 0 {
			err = txDAO.UpdateById(ctx, art)
		} else {
			id, err = txDAO.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		art.Id = id
		return txDAO.Upsert(ctx, PublishedArticle{art})

	})
	return id, err

}

func (dao *articleDAO) Upsert(ctx context.Context, art PublishedArticle) error {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"utime":   art.Utime,
			"content": art.Content,
			"title":   art.Title,
			"status":  art.Status,
		}),
	}).Create(&art).Error
}

func (dao *articleDAO) SyncStatus(ctx context.Context, id int64, author int64, status uint8) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? AND author_id = ?", id, author).Updates(map[string]interface{}{
			"status": status,
			"utime":  now,
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("可能有人在搞你，误操作非自己的文章, uid: %d, aid: %d", author, id)
		}
		return tx.Model(&PublishedArticle{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status": status,
			"utime":  now,
		}).Error
	})
}

type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 长度 1024
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`

	AuthorId int64 `gorm:"index"`
	Status   uint8
	Ctime    int64
	Utime    int64
}
type PublishedArticle struct {
	Article
}
