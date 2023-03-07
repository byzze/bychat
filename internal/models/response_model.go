package models

import "encoding/json"

/************************  响应数据  **************************/

// Head 消息头
type Head struct {
	ChatID   string    `json:"chatID"`   // 会话id
	Seq      string    `json:"seq"`      // 消息的Id
	Cmd      string    `json:"cmd"`      // 消息的cmd 动作
	Response *Response `json:"response"` // 消息体
}

// Response 响应结构
type Response struct {
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Data    interface{} `json:"data"` // 数据 json
}

// NewResponseHead 设置返回消息
func NewResponseHead(seq string, cmd string, code uint32, codeMsg string, data interface{}) *Head {
	response := NewResponse(code, codeMsg, data)

	return &Head{Seq: seq, Cmd: cmd, Response: response}
}

func (h *Head) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)
	return
}

// NewResponse 构造返返回体
func NewResponse(code uint32, codeMsg string, data interface{}) *Response {
	return &Response{Code: code, CodeMsg: codeMsg, Data: data}
}
