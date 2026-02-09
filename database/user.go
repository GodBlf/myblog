package database

import (
	"errors"
	"myblog/util"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	Id     int    `gorm:"column:id;primaryKey"`
	Name   string `gorm:"column:name"`     //name
	PassWd string `gorm:"column:password"` //pass_wd
}

func (User) TableName() string { //gorm能够识别到这个函数，就会用它返回的字符串当作表名，而不是默认的结构体名转蛇形
	return "user"
}

// 反射
var (
	_all_user_field = util.GetGormFields(User{}) //反射要缓存
)

// 根据用户名检索用户
func GetUserByName(name string) *User {
	db := GetBlogDBConnection()
	var user User
	if err := db.Select(_all_user_field).Where("name=?", name).First(&user).
		Error; err != nil { //name不存在重复，所以用First即可
		if !errors.Is(err, gorm.ErrRecordNotFound) { // 如果是用户名不存在，不需要打错误日志
			zap.L().Error("get user by name failed", zap.String("name", name), zap.Error(err))
		}
		return nil
	}
	return &user
}

func CreateUser(name, password string) error {
	db := GetBlogDBConnection()
	user := &User{
		Name:   name,
		PassWd: password,
	}
	err := db.Create(user).Error
	if err != nil {
		zap.L().Error("create user failed", zap.String("name", name), zap.Error(err))
		return err
	}
	zap.L().Info("create user success", zap.String("name", name))
	return nil
}

func DeleteUser(name string) error {
	db := GetBlogDBConnection()
	err := db.Where("name=?", name).Delete(&User{}).Error
	if err != nil {
		zap.L().Error("delete user failed", zap.String("name", name), zap.Error(err))
		return err
	}
	zap.L().Info("delete user success", zap.String("name", name))
	return nil
}
