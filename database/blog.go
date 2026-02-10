package database

import (
	"errors"
	"fmt"
	"myblog/util"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Blog struct {
	Id         int       `gorm:"column:id;primaryKey"`
	UserId     int       `gorm:"column:user_id"`
	Title      string    `gorm:"column:title"`
	Article    string    `gorm:"column:article"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

type PublicBlog struct {
	BlogId      int       `gorm:"column:blog_id;primaryKey"`
	UserId      int       `gorm:"column:user_id;not null"`
	PublishTime time.Time `gorm:"column:publish_time;not null"`
}

type PublicBlogPreview struct {
	Id         int       `gorm:"column:id"`
	UserId     int       `gorm:"column:user_id"`
	UserName   string    `gorm:"column:user_name"`
	Title      string    `gorm:"column:title"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

type PublicBlogDetail struct {
	Id         int       `gorm:"column:id"`
	UserId     int       `gorm:"column:user_id"`
	UserName   string    `gorm:"column:user_name"`
	Title      string    `gorm:"column:title"`
	Article    string    `gorm:"column:article"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (Blog) TableName() string {
	return "blog"
}

func (PublicBlog) TableName() string {
	return "public_blog"
}

var (
	allBlogField       = util.GetGormFields(Blog{})
	publicBlogMigrator sync.Once
)

func ensurePublicBlogTable() {
	db := GetBlogDBConnection()
	publicBlogMigrator.Do(func() {
		if err := db.AutoMigrate(&PublicBlog{}); err != nil {
			zap.L().Error("migrate public_blog failed", zap.Error(err))
		}
	})
}

func GetBlogById(id int) *Blog {
	db := GetBlogDBConnection()
	blog := &Blog{}
	err := db.Select(allBlogField).Where("id=?", id).First(blog).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("get blog by id failed", zap.Int("id", id), zap.Error(err))
		}
		return nil
	}
	return blog
}

func GetBlogByUserId(uid int) []*Blog {
	db := GetBlogDBConnection()
	var blogs []*Blog

	if err := db.Select("id, title").Where("user_id = ?", uid).Find(&blogs).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("get blogs of user failed", zap.Int("user_id", uid), zap.Error(err))
		}
		return nil
	}

	return blogs
}

func GetPublicBlogList() []*PublicBlogPreview {
	ensurePublicBlogTable()
	db := GetBlogDBConnection()

	var blogs []*PublicBlogPreview
	err := db.Table("blog b").
		Select("b.id, b.user_id, u.name AS user_name, b.title, b.update_time").
		Joins("INNER JOIN public_blog pb ON pb.blog_id = b.id").
		Joins("LEFT JOIN `user` u ON u.id = b.user_id").
		Order("pb.publish_time DESC").
		Find(&blogs).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("get public blog list failed", zap.Error(err))
		}
		return nil
	}
	return blogs
}

func GetPublicBlogById(bid int) *PublicBlogDetail {
	ensurePublicBlogTable()
	db := GetBlogDBConnection()

	blog := &PublicBlogDetail{}
	err := db.Table("blog b").
		Select("b.id, b.user_id, u.name AS user_name, b.title, b.article, b.update_time").
		Joins("INNER JOIN public_blog pb ON pb.blog_id = b.id").
		Joins("LEFT JOIN `user` u ON u.id = b.user_id").
		Where("b.id = ?", bid).
		First(blog).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("get public blog by id failed", zap.Int("bid", bid), zap.Error(err))
		}
		return nil
	}
	return blog
}

func IsBlogPublic(bid int) bool {
	ensurePublicBlogTable()
	db := GetBlogDBConnection()

	var count int64
	err := db.Model(&PublicBlog{}).Where("blog_id = ?", bid).Count(&count).Error
	if err != nil {
		zap.L().Error("check blog public status failed", zap.Int("bid", bid), zap.Error(err))
		return false
	}
	return count > 0
}

func PublishBlog(bid, uid int) error {
	if bid <= 0 || uid <= 0 {
		return fmt.Errorf("invalid blog id or user id")
	}
	ensurePublicBlogTable()

	publishTime := time.Now()
	publicBlog := &PublicBlog{
		BlogId:      bid,
		UserId:      uid,
		PublishTime: publishTime,
	}

	db := GetBlogDBConnection()
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "blog_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"user_id":      uid,
			"publish_time": publishTime,
		}),
	}).Create(publicBlog).Error
}

func UnpublishBlog(bid, uid int) error {
	if bid <= 0 || uid <= 0 {
		return fmt.Errorf("invalid blog id or user id")
	}
	ensurePublicBlogTable()

	db := GetBlogDBConnection()
	return db.Where("blog_id = ? AND user_id = ?", bid, uid).Delete(&PublicBlog{}).Error
}

func UpdateBlog(blog *Blog) error {
	if blog.Id <= 0 {
		return fmt.Errorf("could not update blog of id %d", blog.Id)
	}

	if len(blog.Article) == 0 || len(blog.Title) == 0 {
		return fmt.Errorf("could not set blog title or article to empty")
	}

	db := GetBlogDBConnection()
	return db.Model(&Blog{}).Where("id = ?", blog.Id).Updates(map[string]any{
		"title":   blog.Title,
		"article": blog.Article,
	}).Error
}

func CreateBlog(blog *Blog) error {
	if blog.UserId <= 0 {
		return fmt.Errorf("invalid user id")
	}
	if len(blog.Title) == 0 || len(blog.Article) == 0 {
		return fmt.Errorf("could not set blog title or article to empty")
	}
	if blog.UpdateTime.IsZero() {
		blog.UpdateTime = time.Now()
	}

	db := GetBlogDBConnection()
	return db.Create(blog).Error
}
