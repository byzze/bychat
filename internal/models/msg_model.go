package models

import "bychat/internal/common"

// 消息类型
const (
	MessageTypeText = "text"
	MessageTypeImg  = "img"

	MessageCmdMsg       = "msg"
	MessageCmdEnter     = "enter"
	MessageCmdExit      = "exit"
	MessageCmdLogin     = "login"
	MessageCmdHeartbeat = "heartbeat"
)

// Message 消息的定义
type Message struct {
	To   string      `json:"to"`   // 目标
	Type string      `json:"type"` // 消息类型 text/img/
	Msg  interface{} `json:"msg"`  // 消息内容 文本，图像，视频，音频，文件
	From string      `json:"from"` // 发送者
}

// NewTextMsg 文本消息构造
func NewTextMsg(from string, Msg string) (message *Message) {
	message = &Message{
		Type: MessageTypeText,
		From: from,
		Msg:  Msg,
	}

	return
}

// getTextMsgData 获取文本消息
func getTextMsgData(uuid, cmd, msgID, message string) string {
	textMsg := NewTextMsg(uuid, message)
	head := NewResponseHead(msgID, cmd, common.OK, "Ok", textMsg)

	return head.String()
}

// GetMsgData 文本消息
func GetMsgData(uuid, msgID, cmd, message string) string {
	return getTextMsgData(uuid, cmd, msgID, message)
}

// GetTextMsgData 文本消息
func GetTextMsgData(uuid, msgID, message string) string {
	return getTextMsgData(uuid, "msg", msgID, message)
}

// GetTextMsgDataEnter 用户进入消息
func GetTextMsgDataEnter(uuid, msgID, message string) string {
	return getTextMsgData(uuid, "enter", msgID, message)
}

// GetTextMsgDataExit 用户退出消息
func GetTextMsgDataExit(uuid, msgID, message string) string {
	return getTextMsgData(uuid, "exit", msgID, message)
}
