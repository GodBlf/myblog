package test

import (
	"fmt"
	"myblog/util"
	"testing"
)

func TestConfig(t *testing.T) {
	dbViper := util.CreateConfig("mysql")
	dbViper.WatchConfig()              //监听配置文件变化,如果配置更改也能动态修改变量值
	if !dbViper.IsSet("myblog.port") { //查看配置文件中是否存在某个字段
		t.Fatal("配置文件中未找到 myblog.port")
	}
	port := dbViper.GetInt("myblog.port")
	fmt.Println("port:", port)

	//unmarshal
	logViper := util.CreateConfig("log")
	logViper.WatchConfig()
	type LogConfig struct {
		Level string `mapstructure:"level"`
		File  string `mapstructure:"file"`
	}
	lc := &LogConfig{}
	err := logViper.Unmarshal(lc)
	if err != nil {
		panic("解析配置文件失败: " + err.Error())
		t.Fatal()
	}
	fmt.Printf("log config: %+v\n", lc)
}
