package test

import (
	"myblog/database"
	"myblog/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserByName(t *testing.T) {
	util.InitLogger("log")
	testDatas := []struct {
		name   string
		expect string
	}{
		{"tech_guru", "482c811da5d5b4bc6d497ffa98491e38"},
		{"traveler_01", "ebee4f7952a0354d6e79dc79b8fb39e5"},
	}
	for _, data := range testDatas {
		u := database.GetUserByName(data.name)
		t.Run("TestGetUserByName", func(t *testing.T) {
			assert.NotNilf(t, u, "Expected user %s to exist but got nil", data.name)
			assert.Equal(t, data.expect, u.PassWd, "Expected password hash %s but got %s", data.expect, u.PassWd)
		})
	}

}

func TestCreateUser(t *testing.T) {
	util.InitLogger("log")
	err := database.CreateUser("dqq", "123")
	assert.NoError(t, err, "Expected CreateUser to succeed but got error: %v", err)

}
func TestCreateUser1(t *testing.T) {
	util.InitLogger("log")
	tmp := "123"
	md5 := util.Md5(tmp)
	err := database.CreateUser("dqq", md5)
	assert.NoError(t, err, "Expected CreateUser to succeed but got error: %v", err)

}

func TestDeleteUser(t *testing.T) {
	util.InitLogger("log")
	err := database.DeleteUser("dqq")
	assert.NoError(t, err, "Expected DeleteUser to succeed but got error: %v", err)
}
func TestDeleteUser1(t *testing.T) {
	util.InitLogger("log")
	db := database.GetBlogDBConnection()
	u := &database.User{
		PassWd: util.Md5("123"),
	}
	err := db.Model(u).Where("1 = 1", "*").Updates(u).Error
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

}
