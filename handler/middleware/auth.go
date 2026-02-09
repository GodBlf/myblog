package middleware

import (
	"myblog/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	KeyConfig = util.CreateConfig("key")
)

func GetUidFromJwt(jwt string) int {
	_, payload, err := util.VerifyJwt(jwt, util.JWT_SECRET)
	if err != nil {
		return 0
	}
	for k, v := range payload.UserDefined {
		if k == "uid" {
			return int(v.(float64))
		}
	}
	return 0
}

func GetLoginUid(ctx *gin.Context) int {
	token := ctx.Request.Header.Get("auth_token")
	uid := GetUidFromJwt(token)
	return uid
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loginUid := GetLoginUid(ctx)
		if loginUid <= 0 {
			ctx.String(http.StatusForbidden, "auth failed")
			ctx.Abort()
			return
		}
		ctx.Set("uid", loginUid)
		ctx.Next()
	}
}
