package test

import (
	"fmt"
	"myblog/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlogById(t *testing.T) {
	blogId := 1
	u := database.GetBlogById(blogId)
	assert.NotNilf(t, u, "Expected blog with ID %d to exist but got nil", blogId)
	fmt.Println(u.Article)
}

func TestGetBlogByUserId(t *testing.T) {
	uid := 1
	blogs := database.GetBlogByUserId(uid)
	assert.NotNil(t, blogs, "Expected blogs for user ID %d but got nil", uid)
	assert.True(t, len(blogs) > 0, "Expected at least one blog for user ID %d but got %d", uid, len(blogs))
	for _, blog := range blogs {
		fmt.Println(blog.Id, blog.Title)
	}

}

func TestUpdateBlog(t *testing.T) {
	blog := database.Blog{Id: 1, Title: "双十一", Article: "双十一来临喜洋洋，购物狂欢乐无边。电商盛宴满眼芳，心愿成真喜笑颜。"}
	assert.NoError(t, database.UpdateBlog(&blog))
}
