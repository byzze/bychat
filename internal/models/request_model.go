package models

/************************  请求数据  **************************/
// 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            // 消息的唯一Id
	Cmd  string      `json:"cmd"`            // 请求命令字
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// 登录请求数据
type Login struct {
	ServiceToken string `json:"serviceToken"` // 验证用户是否登录
	AppID        uint32 `json:"appId,omitempty"`
	UserID       string `json:"userId,omitempty"`
}

// 心跳请求数据
type HeartBeat struct {
	UserID string `json:"userId,omitempty"`
}

type Msg struct {
	AppID   uint32 `json:"appId"`
	UserID  string `jsong:"userId"`
	MsgID   string `json:"msgId"`
	Message string `json:"message"`
	Cmd     string `json:"cmd"`
}
