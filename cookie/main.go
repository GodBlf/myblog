package cookie

import (
	"myblog/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const COOKIE_NAME = "dqq"

var loggedIn = make(map[string]string, 1000)

func Login(ctx *gin.Context) {
	uid := "8"
	key := util.RandStringRunes(20)
	loggedIn[key] = uid
	zap.L().Debug("cookie value", zap.String("key", key), zap.String("uid", uid))
	ctx.SetCookie(
		COOKIE_NAME,
		key,
		86400*7,
		"/",
		"localhost",
		false,
		true,
	)
	ctx.String(http.StatusOK, "login success")
}

func getUidFromCookie1(ctx *gin.Context) string {
	// http协议里没有cookie这个概念，cookie本质上是header里的一对KV
	for _, cookie := range strings.Split(ctx.Request.Header.Get("cookie"), ";") {
		arr := strings.Split(cookie, "=")
		key := strings.TrimSpace(arr[0])
		value := strings.TrimSpace(arr[1])
		if key == COOKIE_NAME {
			if uid, exists := loggedIn[value]; exists {
				return uid
			}
		}
	}
	return ""
}

func getUidFromCookie2(ctx *gin.Context) string {
	for _, cookie := range ctx.Request.Cookies() {
		if cookie.Name == COOKIE_NAME {
			if uid, exists := loggedIn[cookie.Value]; exists {
				return uid
			}
		}
	}
	return ""
}

func main() {

}
