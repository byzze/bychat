package models

import "encoding/json"

/************************  响应数据  **************************/

// Response websocket 响应结构
type Response struct {
	MsgSeq     string      `json:"msgSeq"`  // 消息的Id
	MsgType    MessageCmd  `json:"msgType"` // 消息的cmd 动作
	Code       uint32      `json:"code"`
	CodeMsg    string      `json:"codeMsg"`
	MsgContent interface{} `json:"data"` // 消息体
}

// NewResponse 设置返回消息
func NewResponse(msgSeq string, code uint32, codeMsg string, msgContent interface{}, msgType MessageCmd) *Response {
	return &Response{MsgSeq: msgSeq, MsgType: msgType, Code: code, CodeMsg: codeMsg, MsgContent: msgContent}
}

func (h *Response) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)
	return
}
