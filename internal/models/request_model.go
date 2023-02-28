package models

/************************  请求数据  **************************/

// Request 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            // 消息的唯一Id
	Cmd  string      `json:"cmd"`            // 请求命令字
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// LoginRequest 登录请求数据
type LoginRequest struct {
	ServiceToken string `json:"serviceToken"` // 验证用户是否登录
	AppID        uint32 `json:"appId,omitempty"`
	UserID       string `json:"userId,omitempty"`
}

// HeartBeatRequest 心跳请求数据
type HeartBeatRequest struct {
	UserID string `json:"userId,omitempty"`
}

// MsgRequest 消息请求
type MsgRequest struct {
	AppID   uint32 `json:"appId"`
	UserID  string `jsong:"userId"`
	MsgID   string `json:"msgId"`
	Message string `json:"message"`
	Cmd     string `json:"cmd"`
}
