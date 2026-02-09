package test

import (
	"myblog/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 模拟 Camel2Snake 逻辑，实际测试时请确保你的函数逻辑一致
//func Camel2Snake(s string) string {
//	// 简单实现：仅演示用
//	var result string
//	for i, r := range s {
//		if i > 0 && r >= 'A' && r <= 'Z' {
//			result += "_"
//		}
//		result += strings.ToLower(string(r))
//	}
//	return result
//}

func TestGetGormFields(t *testing.T) {
	// 定义测试用的结构体
	type TestUser struct {
		ID         uint   `gorm:"primaryKey"`                  // 预期：id (由于没有 column 标签，走 Camel2Snake)
		UserName   string `gorm:"column:login_name;type:text"` // 预期：login_name (解析 column)
		Email      string // 预期：email (走 Camel2Snake)
		Password   string `gorm:"-"`                 // 预期：忽略
		CreatedAt  int64  `gorm:"column:created_at"` // 预期：created_at
		unexported string // 预期：忽略 (不可导出)
	}
	type testType struct {
		name     string
		input    any
		expected []string
	}
	tests := []testType{
		{
			name:     "基础结构体测试",
			input:    TestUser{},
			expected: []string{"id", "login_name", "email", "created_at"},
		},
		{
			name:     "结构体指针测试",
			input:    &TestUser{},
			expected: []string{"id", "login_name", "email", "created_at"},
		},
		{
			name: "带分号的复杂 Tag 测试",
			input: struct {
				RoleName string `gorm:"column:role_id;index;not null"`
			}{},
			expected: []string{"role_id"},
		},
		{
			name: "无 Tag 自动转换蛇形",
			input: struct {
				ProjectName string
				IsActive    bool
			}{},
			expected: []string{"project_name", "is_active"},
		},
		{
			name:     "非结构体输入",
			input:    "not a struct",
			expected: nil,
		},
		{
			name:     "空结构体",
			input:    struct{}{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := util.GetGormFields(tt.input)
			assert.Equal(t, tt.expected, actual, "用例 [%s] 的期望结果不符", tt.name)
		})
	}
}
