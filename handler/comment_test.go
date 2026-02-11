package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"myblog/handler/middleware"
	"myblog/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCommentTestToken(t *testing.T, uid int) string {
	t.Helper()
	header := util.DefaultHeader
	payload := &util.JwtPayload{
		Issue:      "comment-test",
		IssueAt:    time.Now().Unix(),
		Expiration: time.Now().Add(time.Hour).Unix(),
		UserDefined: map[string]any{
			"uid": uid,
		},
	}
	token, err := util.GenJwt(&header, payload, util.JWT_SECRET)
	require.NoError(t, err)
	return token
}

func TestNewPublicBlogCommentsInvalidBid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/blog/public/:bid/comments", NewPublicBlogComments())

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/blog/public/abc/comments", nil)
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Contains(t, writer.Body.String(), "invalid blog id")
}

func TestNewPublicBlogCommentCreateInvalidBid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/blog/public/:bid/comments", middleware.Auth(), NewPublicBlogCommentCreate())

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/blog/public/not-number/comments", strings.NewReader("content=hello"))
	request.Header.Set("auth_token", newCommentTestToken(t, 1))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Contains(t, writer.Body.String(), "invalid blog id")
}

func TestNewPublicBlogCommentCreateAuthFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/blog/public/:bid/comments", middleware.Auth(), NewPublicBlogCommentCreate())

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/blog/public/1/comments", strings.NewReader("content=hello"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusForbidden, writer.Code)
	assert.Contains(t, writer.Body.String(), "auth failed")
}

func TestNewPublicBlogCommentCreateInvalidParameter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/blog/public/:bid/comments", middleware.Auth(), NewPublicBlogCommentCreate())

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/blog/public/1/comments", strings.NewReader(""))
	request.Header.Set("auth_token", newCommentTestToken(t, 1))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Contains(t, writer.Body.String(), "invalid parameter")
}

func TestNewPublicBlogCommentDeleteInvalidBid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/blog/public/:bid/comments/:cid", middleware.Auth(), NewPublicBlogCommentDelete())

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/blog/public/not-number/comments/1", nil)
	request.Header.Set("auth_token", newCommentTestToken(t, 1))
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Contains(t, writer.Body.String(), "invalid blog id")
}

func TestNewPublicBlogCommentDeleteInvalidCommentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/blog/public/:bid/comments/:cid", middleware.Auth(), NewPublicBlogCommentDelete())

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/blog/public/1/comments/not-number", nil)
	request.Header.Set("auth_token", newCommentTestToken(t, 1))
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Contains(t, writer.Body.String(), "invalid comment id")
}

func TestNewPublicBlogCommentDeleteAuthFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/blog/public/:bid/comments/:cid", middleware.Auth(), NewPublicBlogCommentDelete())

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/blog/public/1/comments/1", nil)
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusForbidden, writer.Code)
	assert.Contains(t, writer.Body.String(), "auth failed")
}
