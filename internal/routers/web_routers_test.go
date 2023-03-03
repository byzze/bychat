package routers

import (
	"bychat/lib/redislib"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitRedis() {
	viper.SetConfigName("../../config/app")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Panic(err)
	}
	redislib.InitRedlisClient()
}

func TestRegisterRouter(t *testing.T) {
	InitRedis()
	g := gin.Default()
	InitWeb(g)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/historyMessageList?appID=101", nil)
	g.ServeHTTP(w, req)
	fmt.Printf("w.Body.String(): %v\n", w.Body.String())
}
