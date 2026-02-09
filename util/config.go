package util

import (
	"errors"
	"myblog/global"

	"github.com/spf13/viper"
)

func CreateConfig(file string) *viper.Viper {
	config := viper.New()
	configPath := global.ProjectRootPath + "config/"
	config.AddConfigPath(configPath)
	config.SetConfigName(file)
	config.SetConfigType("yaml")
	configFile := configPath + file + ".yaml"
	err := config.ReadInConfig()
	if err != nil {
		cerror := &viper.ConfigFileNotFoundError{}
		errors.As(err, &cerror)
		if cerror != nil {
			panic("配置文件未找到: " + configFile)
		} else {
			panic("读取配置文件失败: " + err.Error())
		}
	}
	return config
}
