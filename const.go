package workwechat

const (
	// accessTokenURL accessToken 接口地址
	accessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	messageSendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
)
