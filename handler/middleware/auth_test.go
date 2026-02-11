package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"myblog/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestToken(t *testing.T, userDefined map[string]any) string {
	t.Helper()
	header := util.DefaultHeader
	payload := &util.JwtPayload{
		Issue:       "test",
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(time.Hour).Unix(),
		UserDefined: userDefined,
	}
	token, err := util.GenJwt(&header, payload, util.JWT_SECRET)
	require.NoError(t, err)
	return token
}

func TestGetUidFromJwt(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid token", func(t *testing.T) {
		token := newTestToken(t, map[string]any{"uid": 123})
		uid := GetUidFromJwt(token)
		assert.Equal(t, 123, uid)
	})

	t.Run("invalid token returns zero", func(t *testing.T) {
		uid := GetUidFromJwt("not-a-jwt")
		assert.Equal(t, 0, uid)
	})

	t.Run("missing uid returns zero", func(t *testing.T) {
		token := newTestToken(t, map[string]any{"name": "tester"})
		uid := GetUidFromJwt(token)
		assert.Equal(t, 0, uid)
	})

	t.Run("uid type mismatch panics", func(t *testing.T) {
		token := newTestToken(t, map[string]any{"uid": "123"})
		assert.Panics(t, func() {
			_ = GetUidFromJwt(token)
		})
	})
}

func TestGetLoginUid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("read uid from auth_token header", func(t *testing.T) {
		token := newTestToken(t, map[string]any{"uid": 321})

		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("auth_token", token)
		ctx.Request = request

		uid := GetLoginUid(ctx)
		assert.Equal(t, 321, uid)
	})

	t.Run("missing auth_token returns zero", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		uid := GetLoginUid(ctx)
		assert.Equal(t, 0, uid)
	})
}

func TestAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid token aborts with 403", func(t *testing.T) {
		router := gin.New()
		handlerExecuted := false
		router.GET("/protected", Auth(), func(ctx *gin.Context) {
			handlerExecuted = true
			ctx.Status(http.StatusOK)
		})

		writer := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/protected", nil)
		request.Header.Set("auth_token", "invalid-token")

		router.ServeHTTP(writer, request)

		assert.Equal(t, http.StatusForbidden, writer.Code)
		assert.Contains(t, writer.Body.String(), "auth failed")
		assert.False(t, handlerExecuted)
	})

	t.Run("valid token sets uid and continues", func(t *testing.T) {
		router := gin.New()
		handlerExecuted := false
		capturedUid := 0
		router.GET("/protected", Auth(), func(ctx *gin.Context) {
			handlerExecuted = true
			uidValue, exists := ctx.Get("uid")
			if !exists {
				ctx.Status(http.StatusInternalServerError)
				return
			}
			uid, ok := uidValue.(int)
			if !ok {
				ctx.Status(http.StatusInternalServerError)
				return
			}
			capturedUid = uid
			ctx.Status(http.StatusNoContent)
		})

		token := newTestToken(t, map[string]any{"uid": 77})
		writer := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/protected", nil)
		request.Header.Set("auth_token", token)

		router.ServeHTTP(writer, request)

		assert.Equal(t, http.StatusNoContent, writer.Code)
		assert.True(t, handlerExecuted)
		assert.Equal(t, 77, capturedUid)
	})
}
