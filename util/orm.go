package util

import (
	"reflect"
	"strings"
)

func GetGormFields(stc any) []string {
	typ := reflect.TypeOf(stc)
	// 如果传的是指针类型，先解析指针获取基础类型
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Struct {
		columns := make([]string, 0, typ.NumField())
		for i := 0; i < typ.NumField(); i++ {
			fieldType := typ.Field(i)
			// 只关注可导出成员 (Public fields)
			if fieldType.IsExported() {
				// 如果 tag 标记为 "-", 则不做 ORM 映射的字段跳过
				if fieldType.Tag.Get("gorm") == "-" {
					continue
				}

				// 默认逻辑：如果没有 gorm Tag，则把驼峰命名转为蛇形命名
				name := Camel2Snake(fieldType.Name)

				// 如果存在 gorm Tag，尝试解析出其中的 column 定义
				if len(fieldType.Tag.Get("gorm")) > 0 {
					content := fieldType.Tag.Get("gorm")
					if strings.HasPrefix(content, "column:") {
						content = content[7:]
						pos := strings.Index(content, ";")
						if pos > 0 {
							name = content[0:pos]
						} else if pos < 0 {
							name = content
						}
					}
				}
				columns = append(columns, name)
			}
		}
		return columns
	} else {
		// 如果 stc 不是结构体则返回空切片
		return nil
	}
}
