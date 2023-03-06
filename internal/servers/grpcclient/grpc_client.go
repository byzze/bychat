package grpcclient

import (
	"bychat/internal/common"
	"context"
	"fmt"
	"time"

	"bychat/internal/models"
	"bychat/internal/protobuf"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// rpc client
// 给全体用户发送消息
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func SendMsgAll(server *models.ServerNode, appID, roomID, userID uint32, seq, cmd string, message string) (sendMsgID string, err error) {
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
		Seq:    seq,
		AppID:  appID,
		UserID: userID,
		RoomID: roomID,
		Cms:    cmd,
		Msg:    message,
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

	sendMsgID = rsp.GetSendMsgID()
	logrus.Error("给全体用户发送消息 成功:", sendMsgID)

	return
}

// 获取用户列表
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func GetRoomUserList(server *models.ServerNode, appID, roomID uint32) (userIDs []string, err error) {
	userIDs = make([]string, 0)

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

	userIDs = rsp.GetUserID()
	logrus.Info("获取用户列表 成功:", userIDs)
	return
}
