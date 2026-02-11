package test

import (
	"strings"
	"testing"

	"myblog/database"

	"github.com/stretchr/testify/assert"
)

func TestCreatePublicBlogCommentInvalidInput(t *testing.T) {
	t.Run("invalid blog id", func(t *testing.T) {
		err := database.CreatePublicBlogComment(0, 1, "hello")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid blog id")
	})

	t.Run("invalid user id", func(t *testing.T) {
		err := database.CreatePublicBlogComment(1, 0, "hello")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid blog id")
	})

	t.Run("empty content", func(t *testing.T) {
		err := database.CreatePublicBlogComment(1, 1, "   ")
		assert.ErrorIs(t, err, database.ErrInvalidCommentContent)
	})

	t.Run("content too long", func(t *testing.T) {
		tooLong := strings.Repeat("a", 1001)
		err := database.CreatePublicBlogComment(1, 1, tooLong)
		assert.ErrorIs(t, err, database.ErrInvalidCommentContent)
	})
}

func TestDeletePublicBlogCommentInvalidInput(t *testing.T) {
	t.Run("invalid blog id", func(t *testing.T) {
		err := database.DeletePublicBlogComment(0, 1, 1)
		assert.ErrorIs(t, err, database.ErrInvalidDeleteComment)
	})

	t.Run("invalid comment id", func(t *testing.T) {
		err := database.DeletePublicBlogComment(1, 0, 1)
		assert.ErrorIs(t, err, database.ErrInvalidDeleteComment)
	})

	t.Run("invalid user id", func(t *testing.T) {
		err := database.DeletePublicBlogComment(1, 1, 0)
		assert.ErrorIs(t, err, database.ErrInvalidDeleteComment)
	})
}
