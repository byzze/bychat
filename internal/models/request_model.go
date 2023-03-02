package models

/************************  请求数据  **************************/

// Request 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            // 消息的唯一Id
	Cmd  string      `json:"cmd"`            // 请求命令字
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// OpenRequest 登录请求数据
type OpenRequest struct {
	ServiceToken string `json:"serviceToken"` // 验证用户是否登录
	AppID        uint32 `json:"appID,omitempty"`
	RoomID       uint32 `json:"roomID,omitempty"`
	UserID       string `json:"userID,omitempty"`
}

// HeartBeatRequest 心跳请求数据
type HeartBeatRequest struct {
	UserID string `json:"userID,omitempty"`
}

// MsgRequest 消息请求
type MsgRequest struct {
	AppID   uint32 `json:"appID"`
	UserID  string `jsong:"userID"`
	MsgID   string `json:"msgID"`
	Message string `json:"message"`
	Cmd     string `json:"cmd"`
}

// EnterRoomRequest 进入房间请求数据
type EnterRoomRequest struct {
	AppID  uint32 `json:"appID,omitempty"`
	RoomID uint32 `json:"roomID,omitempty"`
	UserID string `json:"userID,omitempty"`
}
