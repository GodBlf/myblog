package test

import (
	"fmt"
	"myblog/util"
	"strings"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	secret := "123456"
	header := &util.DefaultHeader
	payload := &util.JwtPayload{
		ID:          "rj4t49tu49",
		Issue:       "微信",
		Audience:    "王者荣耀",
		Subject:     "购买道具",
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(2 * time.Hour).Unix(),
		UserDefined: map[string]any{"name": strings.Repeat("大乔乔", 1)}, //信息量很大时, jwt长度可能会超过4K
	}

	if token, err := util.GenJwt(header, payload, secret); err != nil {
		fmt.Printf("生成json web token失败: %v", err)
	} else {
		fmt.Println(token)
		if _, p, err := util.VerifyJwt(token, secret); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("JWT验证通过。欢迎 %s !\n", p.UserDefined["name"])
		}
	}
}
