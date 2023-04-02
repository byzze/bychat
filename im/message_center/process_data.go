package messagecenter

import (
	"bychat/im/client"
	"bychat/im/models"
	"bychat/pkg/common"
	"encoding/json"
	"sync"

	"github.com/sirupsen/logrus"
)

// DisposeFunc 处理函数
type DisposeFunc func(c *client.Client, msgSeq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[models.MessageCmd]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// Register 注册
func Register(key models.MessageCmd, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandlers(key models.MessageCmd) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]
	return
}

// ProcessData websocket处理数据
func ProcessData(c *client.Client, data []byte) {
	logrus.WithFields(logrus.Fields{
		"data": string(data),
	}).Info("ProcessData Request")

	var req = &models.Request{}
	err := json.Unmarshal(data, req)
	if err != nil {
		logrus.WithError(err).Error("ProcessData Unmarshal")
		ReponseMsg(c, common.ParameterIllegal, "", "data format is invalid", "", "")
		return
	}

	msgSeq := req.MsgSeq
	msgCmd := models.MessageCmd(req.MsgCmd)

	requestData, err := json.Marshal(req.MsgContent)
	if err != nil {
		logrus.WithError(err).Error("handle json Marshal")
		ReponseMsg(c, common.ParameterIllegal, msgSeq, "data format is invalid", "", msgCmd)
		return
	}

	var (
		code    uint32
		msg     string
		message interface{}
	)

	if v, ok := getHandlers(msgCmd); ok {
		code, msg, message = v(c, msgSeq, requestData)
	} else {
		ReponseMsg(c, common.RoutingNotExist, msgSeq, "router not found", message, msgCmd)
		return
	}

	ReponseMsg(c, code, msgSeq, msg, message, msgCmd)
}
