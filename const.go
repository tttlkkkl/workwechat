package workwechat

const (
	// accessTokenURL accessToken 接口地址
	accessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	messageSendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	updateTaskCard = "https://qyapi.weixin.qq.com/cgi-bin/message/update_taskcard?access_token=%s"
)

// MessageType 消息类型
type MessageType string

const (
	// MessageTypeText 文本消息
	MessageTypeText MessageType = "text"
	// MessageTypeTaskCard 任务卡片消息
	MessageTypeTaskCard MessageType = "taskcard"
	// MessageTypeEvent 事件消息
	MessageTypeEvent MessageType = "event"
)

// EventType 事件类型
type EventType string

const (
	// EventTypeTaskCard 任务卡片事件
	EventTypeTaskCard EventType = "taskcard_click"
)
