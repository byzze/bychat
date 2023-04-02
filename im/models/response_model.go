package models

import "encoding/json"

/************************  响应数据  **************************/

// Response websocket 响应结构
type Response struct {
	MsgSeq  string      `json:"msgSeq"`  // 消息的Id
	MsgCmd  MessageCmd  `json:"msgType"` // 消息的cmd 动作
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	MsgBody interface{} `json:"data"` // 消息体
}

// NewResponse 设置返回消息
func NewResponse(msgSeq string, code uint32, codeMsg string, message interface{}, msgType MessageCmd) *Response {
	return &Response{MsgSeq: msgSeq, MsgCmd: msgType, Code: code, CodeMsg: codeMsg, MsgBody: message}
}

func (h *Response) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)
	return
}
