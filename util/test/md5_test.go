package test

import (
	"myblog/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMd5(t *testing.T) {
	// 准备测试用例
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空字符串",
			input:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "标准字符串 - hello",
			input:    "hello",
			expected: "5d41402abc4b2a76b9719d911017c592", // 注意：hello world 的 MD5
		},
		{
			name:     "包含特殊字符",
			input:    "123!@#",
			expected: "1fdb7184e697ab9355a3f1438ddc6ef9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用 testify/assert 进行断言
			actual := util.Md5(tt.input)

			// Equal 会比较预期值和实际值，如果不一致会打印漂亮的 Diff
			assert.Equal(t, tt.expected, actual, "输入数据为 [%s] 时，MD5 结果不匹配", tt.input)
		})
	}
}
