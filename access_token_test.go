package workwechat

import (
	"testing"
)

const (
	corpID  = "xx"
	secret  = "xx"
	agentID = 1000034
)

var client *WorkWechat

func TestWorkWechat_GetAccessToken(t *testing.T) {
	var err error
	client, err = NewWorkWechat(corpID, secret, agentID)
	if err != nil {
		t.Errorf("init error:%v", err)
	}
	if token, err := client.GetAccessToken(); err != nil || token == "" {
		t.Errorf("accessToken获取失败!---%s---%v", token, err)
	}
}
