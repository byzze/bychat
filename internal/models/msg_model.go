package models

import (
	"bychat/internal/common"
)

// MessageType 消息类型枚举
type MessageType string

// 消息类型
const (
	MessageTypeText  MessageType = "text"
	MessageTypeImg   MessageType = "img"
	MessageTypeVedio MessageType = "vedio"
	MessageTypeFile  MessageType = "file"
	MessageTypeSound MessageType = "sound"
)

// MessageCmd 消息命令枚举
type MessageCmd string

// 消息指令类型
const (
	MessageCmdMsg       MessageCmd = "msg"
	MessageCmdEnter     MessageCmd = "enter"
	MessageCmdExit      MessageCmd = "exit"
	MessageCmdLogin     MessageCmd = "login"
	MessageCmdHeartbeat MessageCmd = "heartbeat"
)

// Message 消息的定义
type Message struct {
	To      string   `json:"to,omitempty"`      // 目标
	From    string   `json:"from,omitempty"`    // 发送者
	MsgBody *MsgBody `json:"msgBody,omitempty"` // 消息内容 文本，图像，视频，音频，文件

}

// MsgBody 消息体的定义
type MsgBody struct {
	MsgType    MessageType `json:"msgType"`    // 消息类型 text/img/
	MsgContent interface{} `json:"msgContent"` // 消息内容 文本，图像，视频，音频，文件 TextMessage ImgMessage VideoMessage SoundMessage FileMessage
}

// TextMessage 消息的定义
type TextMessage struct {
	Text string `json:"text"`
}

// ImgMessage 消息的定义
type ImgMessage struct {
	URL    string `json:"url"`
	Size   int64  `json:"size"`
	Name   string `json:"name"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
}

// VideoMessage 消息的定义
type VideoMessage struct {
	URL    string `json:"url"`
	Name   string `json:"name"`
	Format string `json:"format"`
	Size   int64  `json:"size"`
	Second int32  `json:"second"`
}

// SoundMessage 消息的定义
type SoundMessage struct {
	URL    string `json:"url"`
	Size   int64  `json:"size"`
	Second int32  `json:"second"`
}

// FileMessage 消息的定义
type FileMessage struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
	Name string `json:"name"`
}

func newMsgResponse(from, to string, body *MsgBody, msgSeq string, cmd MessageCmd) string {
	msg := NewMsg(from, to, body)
	res := NewResponse(msgSeq, common.OK, "", msg, cmd)
	return res.String()
}

// NewTextMsgBody 文本消息构造
func NewTextMsgBody(text string) (msgBody *MsgBody) {
	msgBody = &MsgBody{
		MsgType:    MessageTypeText,
		MsgContent: TextMessage{Text: text},
	}
	return
}

// NewImgMsgNody 图片消息
func NewImgMsgNody(url, name string, size int64, width, height int32) (msgBody *MsgBody) {
	msgBody = &MsgBody{
		MsgType: MessageTypeImg,
		MsgContent: ImgMessage{
			URL:    url,
			Name:   name,
			Size:   size,
			Width:  width,
			Height: height,
		},
	}
	return
}

// NewVedioMsgBody 视频消息
func NewVedioMsgBody(url, name, format string, size int64, second int32) (msgBody *MsgBody) {
	msgBody = &MsgBody{
		MsgType: MessageTypeVedio,
		MsgContent: VideoMessage{
			URL:    url,
			Name:   name,
			Format: format,
			Size:   size,
			Second: second,
		},
	}
	return
}

// NewSoundMsgBody 音频消息
func NewSoundMsgBody(url string, size int64, second int32) (msgBody *MsgBody) {
	msgBody = &MsgBody{
		MsgType: MessageTypeSound,
		MsgContent: SoundMessage{
			URL:    url,
			Size:   size,
			Second: second,
		},
	}
	return
}

// NewFileMsgBody 文件消息
func NewFileMsgBody(url, name string, size int64) (msgBody *MsgBody) {
	msgBody = &MsgBody{
		MsgType: MessageTypeFile,
		MsgContent: FileMessage{
			URL:  url,
			Name: name,
			Size: size,
		},
	}
	return
}

// NewMsg 消息结构
func NewMsg(from, to string, body *MsgBody) (msg *Message) {
	msg = &Message{
		From:    from,
		To:      to,
		MsgBody: body,
	}
	return
}

// newTextMsgData 封装文本消息
func newTextMsgData(from, to, msgSeq, msgContent string, cmd MessageCmd) string {
	body := NewTextMsgBody(msgContent)
	return newMsgResponse(from, to, body, msgSeq, cmd)
}

// newImgMsgData 封装图片消息
func newImgMsgData(from, to, msgSeq, url, name string, size int64, width, height int32, cmd MessageCmd) string {
	body := NewImgMsgNody(url, name, size, width, height)
	return newMsgResponse(from, to, body, msgSeq, cmd)
}

// newVedioMsgData 封装视频消息
func newVedioMsgData(from, to, msgSeq, url, name, format string, size int64, second int32, cmd MessageCmd) string {
	body := NewVedioMsgBody(url, name, format, size, second)
	return newMsgResponse(from, to, body, msgSeq, cmd)
}

// newSoundMsgData 封装音频消息
func newSoundMsgData(from, to, msgSeq, url string, size int64, second int32, cmd MessageCmd) string {
	body := NewSoundMsgBody(url, size, second)
	return newMsgResponse(from, to, body, msgSeq, cmd)
}

// newFileMsgData 封装文件消息
func newFileMsgData(from, to, msgSeq, url, name string, size int64, cmd MessageCmd) string {
	body := NewFileMsgBody(url, name, size)
	return newMsgResponse(from, to, body, msgSeq, cmd)
}

// GetTextMsgData 文本消息
func GetTextMsgData(from, to, msgSeq, msgContent string) string {
	return newTextMsgData(from, to, msgSeq, msgContent, MessageCmdMsg)
}

// GetTextMsgDataEnter 用户进入文本消息
func GetTextMsgDataEnter(from, to, msgSeq, msgContent string) string {
	return newTextMsgData(from, to, msgSeq, msgContent, MessageCmdEnter)
}

// GetTextMsgDataExit 用户退出文本消息
func GetTextMsgDataExit(from, to, msgSeq, msgContent string) string {
	return newTextMsgData(from, to, msgSeq, msgContent, MessageCmdExit)
}

// GetImgMsgData 图片消息
func GetImgMsgData(from, to, msgSeq, url, name string, size int64, width, height int32) string {
	return newImgMsgData(from, to, msgSeq, url, name, size, width, height, MessageCmdMsg)
}

// GetVedioMsgData 视频消息
func GetVedioMsgData(from, to, msgSeq, url, name, format string, size int64, second int32) string {
	return newVedioMsgData(from, to, msgSeq, url, name, format, size, second, MessageCmdMsg)
}

// GetSoundMsgData 音频消息
func GetSoundMsgData(from, to, msgSeq, url string, size int64, second int32) string {
	return newSoundMsgData(from, to, msgSeq, url, size, second, MessageCmdMsg)
}

// GetFileMsgData 文件消息
func GetFileMsgData(from, to, msgSeq, url, name string, size int64) string {
	return newFileMsgData(from, to, msgSeq, url, name, size, MessageCmdMsg)
}
