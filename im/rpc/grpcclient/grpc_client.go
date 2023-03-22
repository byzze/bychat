package grpcclient

import (
	"bychat/im/models"
	"bychat/im/rpc/protobuf"
	"bychat/pkg/common"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// SendMsgAll 给全体用户发送消息 link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func SendMsgAll(server *models.ServerNode, appID, roomID, userID uint32, message string) (err error) {
	logrus.WithFields(logrus.Fields{
		"appID":   appID,
		"roomID":  roomID,
		"userID":  userID,
		"message": message,
	}).Info("grpc client SendMsgAll")
	// Set up a connection to the server.
	conn, err := grpc.Dial(server.String(), grpc.WithInsecure())
	if err != nil {
		logrus.Error("连接失败", server.String())
		return
	}
	defer conn.Close()

	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := protobuf.SendMsgAllReq{
		AppID:  appID,
		UserID: userID,
		RoomID: roomID,
		Data:   message,
	}
	rsp, err := c.SendMsgAll(ctx, &req)
	if err != nil {
		logrus.Error("给全体用户发送消息", err)
		err = fmt.Errorf("发送消息失败 err:%d", err)
		return
	}

	if rsp.GetRetCode() != common.OK {
		logrus.Error("给全体用户发送消息", rsp.String())
		err = fmt.Errorf("发送消息失败 code:%d", rsp.GetRetCode())
		return
	}

	logrus.Info("给全体用户发送消息 成功")
	return
}

// GetRoomUserList 获取用户列表 link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func GetRoomUserList(server *models.ServerNode, appID, roomID uint32) (userList []uint32, err error) {
	userList = make([]uint32, 0)

	conn, err := grpc.Dial(server.String(), grpc.WithInsecure())
	if err != nil {
		logrus.Error("连接失败", server.String())
		return
	}
	defer conn.Close()

	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := protobuf.GetRoomUserListReq{
		AppID:  appID,
		RoomID: roomID,
	}
	rsp, err := c.GetRoomUserList(ctx, &req)
	if err != nil {
		logrus.Error("获取用户列表 发送请求错误:", err)
		return
	}

	if rsp.GetRetCode() != common.OK {
		logrus.Error("获取用户列表 返回码错误:", rsp.String())
		err = fmt.Errorf("发送消息失败 code:%d", rsp.GetRetCode())
		return
	}

	for _, v := range rsp.GetResUserOnline() {
		userList = append(userList, v.Id)
	}
	logrus.Info("rpc 获取用户列表 成功:", userList)
	return
}
