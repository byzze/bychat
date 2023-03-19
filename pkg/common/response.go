package common

// JSONResult 响应结果
type JSONResult struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response 响应结果
func Response(code uint32, message string, data interface{}) JSONResult {
	message = GetErrorMessage(code, message)
	jsonMap := grantMap(code, message, data)

	return jsonMap
}

// 按照接口格式生成原数据数组
func grantMap(code uint32, message string, data interface{}) JSONResult {
	jsonMap := JSONResult{
		Code: code,
		Msg:  message,
		Data: data,
	}
	return jsonMap
}
