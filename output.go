package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Output 用做Ajax返回
type Output struct {
	Status  string // "succes" or "error"
	Message string // "error" 时为错误信息
	Data    interface{}
}

// JSON 获取Output的Json字符串
func (o Output) JSON() string {
	outputjson := `{"Status":"error","Message":"Convert Error.","Data":null}`
	outputBytes, err := json.Marshal(o)
	if err == nil {
		outputjson = string(outputBytes)
	}
	return outputjson
}

// Writer2Response 将Output写入到ResponseWriter中
func (o Output) Writer2Response(w http.ResponseWriter) (n int, err error) {
	return fmt.Fprintln(w, o.JSON())
}
