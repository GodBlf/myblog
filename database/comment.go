package database

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const maxCommentContentLength = 1000

var (
	ErrInvalidCommentContent = errors.New("invalid comment content")
	ErrPublicBlogNotExist    = errors.New("public blog not exist")
	ErrCommentNotExist       = errors.New("comment not exist")
	ErrCommentNoPermission   = errors.New("no permission to delete comment")
	ErrInvalidDeleteComment  = errors.New("invalid delete comment parameter")
)

type BlogComment struct {
	Id         int       `gorm:"column:id;primaryKey"`
	BlogId     int       `gorm:"column:blog_id;not null;index"`
	UserId     int       `gorm:"column:user_id;not null;index"`
	Content    string    `gorm:"column:content;not null;type:text"`
	CreateTime time.Time `gorm:"column:create_time"`
}

type PublicBlogCommentItem struct {
	Id         int       `gorm:"column:id"`
	BlogId     int       `gorm:"column:blog_id"`
	UserId     int       `gorm:"column:user_id"`
	UserName   string    `gorm:"column:user_name"`
	Content    string    `gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (BlogComment) TableName() string {
	return "blog_comment"
}

var blogCommentMigrator sync.Once

func ensureBlogCommentTable() {
	db := GetBlogDBConnection()
	blogCommentMigrator.Do(func() {
		if err := db.AutoMigrate(&BlogComment{}); err != nil {
			zap.L().Error("migrate blog_comment failed", zap.Error(err))
		}
	})
}

func CreatePublicBlogComment(bid int, uid int, content string) error {
	if bid <= 0 || uid <= 0 {
		return fmt.Errorf("invalid blog id or user id")
	}
	trimmedContent := strings.TrimSpace(content)
	if len(trimmedContent) == 0 || len([]rune(trimmedContent)) > maxCommentContentLength {
		return ErrInvalidCommentContent
	}
	if !IsBlogPublic(bid) {
		return ErrPublicBlogNotExist
	}

	ensureBlogCommentTable()
	comment := &BlogComment{
		BlogId:     bid,
		UserId:     uid,
		Content:    trimmedContent,
		CreateTime: time.Now(),
	}

	db := GetBlogDBConnection()
	return db.Create(comment).Error
}

func GetPublicBlogComments(bid int) []*PublicBlogCommentItem {
	if bid <= 0 {
		return nil
	}

	ensureBlogCommentTable()
	db := GetBlogDBConnection()

	var comments []*PublicBlogCommentItem
	err := db.Table("blog_comment bc").
		Select("bc.id, bc.blog_id, bc.user_id, u.name AS user_name, bc.content, bc.create_time").
		Joins("LEFT JOIN `user` u ON u.id = bc.user_id").
		Where("bc.blog_id = ?", bid).
		Order("bc.create_time DESC").
		Find(&comments).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			zap.L().Error("get public blog comments failed", zap.Int("bid", bid), zap.Error(err))
		}
		return nil
	}

	return comments
}

func DeletePublicBlogComment(bid int, cid int, uid int) error {
	if bid <= 0 || cid <= 0 || uid <= 0 {
		return ErrInvalidDeleteComment
	}
	if !IsBlogPublic(bid) {
		return ErrPublicBlogNotExist
	}

	ensureBlogCommentTable()
	db := GetBlogDBConnection()

	comment := &BlogComment{}
	err := db.Where("id = ? AND blog_id = ?", cid, bid).First(comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotExist
		}
		return err
	}

	if comment.UserId != uid {
		return ErrCommentNoPermission
	}

	result := db.Where("id = ? AND blog_id = ? AND user_id = ?", cid, bid, uid).Delete(&BlogComment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCommentNotExist
	}

	return nil
}
