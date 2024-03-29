package cache

import (
	"bychat/internal/redislib"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	submitAgainPrefix = "acc:submit:again:" // 数据不重复提交

	seqDuplicatesDefaultTime = 30
)

/*********************  查询数据是否处理过  ************************/

// 获取数据提交去除key
func getSubmitAgainKey(from string, value string) (key string) {
	key = fmt.Sprintf("%s%s:%s", submitAgainPrefix, from, value)
	return
}

// 重复提交
// return true:重复提交 false:第一次提交
func submitAgain(from string, second int, value string) (isSubmitAgain bool) {
	// 默认重复提交
	isSubmitAgain = true
	key := getSubmitAgainKey(from, value)

	redisClient := redislib.GetClient()
	number, err := redisClient.Do("setNx", key, "1").Int()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key":    key,
			"number": number,
			"err":    err,
		}).Error("submitAgain")
		return
	}

	if number != 1 {
		return
	}
	// 第一次提交
	isSubmitAgain = false

	redisClient.Do("Expire", key, second)
	return

}

// SeqDuplicates 重复提交
func SeqDuplicates(seq string) (result bool) {
	result = submitAgain("seq", seqDuplicatesDefaultTime, seq)
	return
}
