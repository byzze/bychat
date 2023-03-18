package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig(cname string) {
	viper.SetConfigName("configs/" + cname)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Panic("Fatal error config file:", err)
	}
	// 获取app属性，redis属性
	logrus.Info("config app:", viper.Get("app"))
	logrus.Info("config redis:", viper.Get("redis"))
}
