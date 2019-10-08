package workwechat

import (
	"testing"
)

const (
	corpID = ""
	secret = ""
)

func TestWorkWechat_GetAccessToken(t *testing.T) {
	client := NewWorkWechat(corpID, secret)
	if token, err := client.GetAccessToken(); err != nil || token == "" {
		t.Errorf("accessToken获取失败!---%s---%v", token, err)
	}
}
