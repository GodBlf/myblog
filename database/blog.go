package database

import (
	"errors"
	"fmt"
	"myblog/util"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Blog struct {
	Id         int       `gorm:"column:id;primaryKey"`
	UserId     int       `gorm:"column:user_id"`
	Title      string    `gorm:"column:title"`
	Article    string    `gorm:"column:article"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (Blog) TableName() string {
	return "blog"
}

var (
	_all_blog_field = util.GetGormFields(Blog{}) //反射要缓存
)

func GetBlogById(id int) *Blog {
	db := GetBlogDBConnection()
	blog := &Blog{}
	err := db.Select(_all_blog_field).Where("id=?", id).First(blog).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("get blog by id failed", zap.Int("id", id), zap.Error(err))
		}
		return nil
	}
	return blog
}

// 根据作者id获取博客列表(仅包含博客id和标题)
func GetBlogByUserId(uid int) []*Blog {
	db := GetBlogDBConnection()
	var blogs []*Blog

	if err := db.Select("id, title").Where("user_id = ?", uid).Find(&blogs).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			zap.L().Error("get blogs of user failed", zap.Int("user_id", uid), zap.Error(err))
		}
		return nil
	}

	return blogs
}

// 根据博客id更新标题和正文
func UpdateBlog(blog *Blog) error {
	if blog.Id <= 0 {
		return fmt.Errorf("could not update blog of id %d", blog.Id) // errors.New("")
	}

	if len(blog.Article) == 0 || len(blog.Title) == 0 { //判空
		return fmt.Errorf("could not set blog title or article to empty")
	}

	db := GetBlogDBConnection()
	return db.Model(&Blog{}).Where("id = ?", blog.Id).Updates(map[string]any{
		"title":   blog.Title,
		"article": blog.Article,
	}).Error
}
