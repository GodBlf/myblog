package test

import (
	"encoding/json"
	"fmt"
	"myblog/util"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
)

func TestCamel2Snake(t *testing.T) {
	datas := []struct {
		test   string
		expect string
	}{
		{
			test:   "CamelCase",
			expect: "camel_case",
		},
		{
			test:   "HTTPServer",
			expect: "http_server",
		},
		{
			test:   "HelloWorld",
			expect: "hello_world",
		},
	}
	for _, data := range datas {
		t.Run("TestCamel2Snake", func(t *testing.T) {
			assert.Equal(t, data.expect, util.Camel2Snake(data.test), "Expected %s but got %s", data.expect, util.Camel2Snake(data.test))
		})
	}

}

func TestRandStringRunes(t *testing.T) {
	fmt.Println(util.RandStringRunes(10))
	fmt.Println(util.RandStringRunes(30))
}

type User struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	height   float32
	Birthday time.Time `json:"birthday"`
}

var user *User = &User{
	Name:     "Alice",
	Age:      30,
	height:   5.6,
	Birthday: time.Date(1993, time.March, 15, 0, 0, 0, 0, time.UTC),
}

// 基准测试
func BenchmarkStdJson(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marshal, _ := json.Marshal(user)
		tmp := &User{}
		_ = json.Unmarshal(marshal, tmp)
	}
}
func BenchmarkSonicJson(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marshal, _ := sonic.Marshal(user)
		ans := &User{}
		_ = sonic.Unmarshal(marshal, ans)
	}
}
