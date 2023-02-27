package websocket

import (
	"bychat/internal/models"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func ProcessData(c *Client, message []byte) {
	var req = &models.Request{}
	err := json.Unmarshal(message, req)
	if err != nil {
		logrus.Error(err)
		return
	}
	switch req.Cmd {
	case models.MessageCmdLogin:
		logrus.Infof("ProcessData Login:%s,%s", string(message), req.Cmd)
		Login(c, req.Data, req.Cmd)
	case models.MessageCmdHeartbeat:
		logrus.Infof("ProcessData Heartbeat:%s,%s", string(message), req.Cmd)
	case models.MessageCmdMsg:
	}
}
