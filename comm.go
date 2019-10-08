package workwechat

import (
	"strings"
)

// Response 公共返回结果
type Response struct {
	ErrCode    int64  `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

// StringArray 数组字符串
type StringArray []string

// MarshalJSON 自定义数组字符串json编码
func (s *StringArray) MarshalJSON() ([]byte, error) {
	str := strings.Join(*s, "|")
	return json.Marshal(str)
}

// UnmarshalJSON 自定义数组字符串json解码
func (s *StringArray) UnmarshalJSON(b []byte) error {
	strArr := strings.Split(string(b), "|")
	*s = strArr
	return nil
}
