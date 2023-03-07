package grpcserver

import (
	"bychat/internal/common"
	"bychat/internal/protobuf"
	"bychat/internal/websocket"
	"bychat/lib/cache"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type server struct {
	protobuf.UnimplementedAccServerServer
}

func setErr(rsp proto.Message, code uint32, message string) {
	message = common.GetErrorMessage(code, message)
	switch v := rsp.(type) {
	case *protobuf.SendMsgRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.SendMsgAllRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.GetRoomUserListRsp:
		v.RetCode = code
		v.ErrMsg = message
	default:

	}

}

func (server *server) SendMsgAll(c context.Context, req *protobuf.SendMsgAllReq) (rsp *protobuf.SendMsgAllRsp, err error) {
	rsp = &protobuf.SendMsgAllRsp{}

	websocket.AllSendMessages(req.GetAppID(), req.GetRoomID(), req.GetUserID(), req.GetMsg())

	setErr(rsp, common.OK, "")

	logrus.Info("grpc_response 给本机全体用户发消息:", rsp.String())
	return
}

// GetRoomUserList 获取本机用户列表
func (server *server) GetRoomUserList(c context.Context, req *protobuf.GetRoomUserListReq) (rsp *protobuf.GetRoomUserListRsp,
	err error) {

	fmt.Println("grpc_request 获取本机用户列表", req.String())

	// appID := req.GetAppID()
	rsp = &protobuf.GetRoomUserListRsp{}

	// 本机
	userResList := cache.GetChatRoomUser(req.GetRoomID())

	setErr(rsp, common.OK, "")
	var userList []*protobuf.ResponUserOnline
	for _, v := range userResList {
		tmp := &protobuf.ResponUserOnline{
			Id:       v.ID,
			NickName: v.NickName,
			Avatar:   v.Avatar,
		}
		userList = append(userList, tmp)
	}
	rsp.ResUserOnline = userList

	logrus.Info("grpc_response 获取用户列表:", rsp.String())

	return
}

// Init 初始化grpc
func Init() {
	rpcPort := viper.GetString("app.rpcPort")
	fmt.Println("rpc server 启动", rpcPort)

	lis, err := net.Listen("tcp", ":"+rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protobuf.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
