package models

// MessageCmd 消息类型枚举
type MessageCmd string

// 指令类型
const (
	MessageCmdMsg       MessageCmd = "msg"
	MessageCmdLogout    MessageCmd = "logout"
	MessageCmdLogin     MessageCmd = "login"
	MessageCmdBindUser  MessageCmd = "bindUser"
	MessageCmdHeartbeat MessageCmd = "heartbeat"
)

// MessageType 消息类型枚举
type MessageType string

// 消息类型枚举
const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
)

// Message 消息的定义
type Message struct {
	To      string      `json:"to,omitempty"`   // 目标
	From    string      `json:"from,omitempty"` // 发送者
	MsgType MessageType `json:"msgType"`
	MsgBody interface{} `json:"msgBody,omitempty"` // 消息内容 文本，图像，视频，音频，文件
}
