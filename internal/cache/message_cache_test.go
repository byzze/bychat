package cache

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitRedis() {
	viper.SetConfigName("../../config/app")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Panic("Fatal error config file: %s \n", err)
	}
	// redislib.InitRedlisClient()
}
func TestZSetMessage(t *testing.T) {
	t.Run("one", func(t *testing.T) {
		InitRedis()
		ZSetMessage(1001, "一条消息")
		res, err := ZGetMessageAll(1001)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info(res)
	})
}
