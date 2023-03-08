package websocket

import (
	"bychat/internal/common"
	"bychat/internal/models"
	"encoding/json"
	"sync"

	"github.com/sirupsen/logrus"
)

// DisposeFunc 处理函数
type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

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

// ProcessData 处理数据
func ProcessData(c *Client, message []byte) {
	logrus.WithFields(logrus.Fields{
		"addr": c.Addr,
		"data": string(message),
	}).Info("ProcessData Request")

	var req = &models.Request{}
	err := json.Unmarshal(message, req)
	if err != nil {
		logrus.Error(err)
		return
	}
	requestData, err := json.Marshal(req.Data)
	if err != nil {
		logrus.Error("处理数据 json Marshal", err)
		c.SendMsg([]byte("处理数据失败"))
		return
	}

	seq := req.MsgSeq
	cmd := models.MessageCmd(req.Cmd)

	var (
		code uint32
		msg  string
		data interface{}
	)

	if v, ok := getHandlers(cmd); ok {
		code, msg, data = v(c, seq, requestData)
	} else {
		code = common.RoutingNotExist
		logrus.WithFields(logrus.Fields{
			"client.Addr": c.Addr,
			"cmd":         cmd,
		}).Error("处理数据 路由不存在")
	}

	msg = common.GetErrorMessage(code, msg)

	responseHead := models.NewResponse(seq, code, msg, data, cmd)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		logrus.Error("处理数据 json Marshal", err)
		return
	}

	c.SendMsg(headByte)

	logrus.WithFields(logrus.Fields{
		"cmd":      cmd,
		"code":     code,
		"headByte": string(headByte),
	}).Info("acc_response send")
}
