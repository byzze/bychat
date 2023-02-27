package websocket

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

func Login(c *Client, msg interface{}, cmd string) {
	curtime := time.Now().UnixNano()

	cl := NewClient(c.Socket.LocalAddr().Network(), c.Socket, uint64(curtime))
	clientManager.AddClients(cl)

	b, err := json.Marshal(msg)
	if err != nil {
		logrus.Error(err)
		return
	}
	var cls = &login{}
	err = json.Unmarshal(b, cls)
	if err != nil {
		logrus.WithField("err", err.Error()).Error("Login")
		return
	}
	cls.Client = cl
	cl.UserID = cls.UserID
	cl.AppID = cls.AppID
	clientManager.Login <- cls
}
