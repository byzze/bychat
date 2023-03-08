package models

/************************  请求数据  **************************/

// Request websocket 通用请求数据格式
type Request struct {
	MsgSeq string      `json:"msgSeq"` // 消息的唯一Id
	Cmd    string      `json:"cmd"`    // 请求命令字
	Data   interface{} `json:"data"`   // 数据 json
}

// OpenRequest 登录请求数据
type OpenRequest struct {
	Token  string `json:"token"` // 验证用户是否登录
	AppID  uint32 `json:"appID"`
	UserID uint32 `json:"userID"`
}

// HeartBeatRequest 心跳请求数据
type HeartBeatRequest struct {
	UserID uint32 `json:"userID"`
}

// MsgRequest 消息请求
type MsgRequest struct {
	AppID          uint32 `json:"appID"`
	UserID         uint32 `json:"userID"`
	MsgSeq         string `json:"msgSeq"`
	MessageContent string `json:"messageContent"`
	Cmd            string `json:"cmd"`
}
