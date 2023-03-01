package common

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// 设置日志格式为json格式
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
}

// SetOutPutFile 设置输出文件
func SetOutPutFile(level logrus.Level) {
	dir := viper.GetString("log.dir")
	name := viper.GetString("log.name")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic(fmt.Errorf("create log dir '%s' error: %s", dir, err))
		}
	}

	fileName := path.Join(dir, name)

	var err error
	os.Stderr, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("open log file err", err)
	}

	mw := io.MultiWriter(os.Stdout, os.Stderr)
	logrus.SetOutput(mw)
	logrus.SetLevel(level)
	return
}

// 实现日志滚动。
// Refer to https://www.cnblogs.com/jssyjam/p/11845475.html.
// logger := &lumberjack.Logger{
// 	Filename:   fmt.Sprintf("%v/%v", dir, name), // 日志输出文件路径。
// 	MaxSize:    LogConf.MaxSize,                                 // 日志文件最大 size(MB)，缺省 100MB。
// 	MaxBackups: 10,                                              // 最大过期日志保留的个数。
// 	MaxAge:     30,                                              // 保留过期文件的最大时间间隔，单位是天。
// 	LocalTime:  true,                                            // 是否使用本地时间来命名备份的日志。
// }
// logrus.SetOutput(logger)
