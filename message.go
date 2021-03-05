package workwechat

import (
	"fmt"
	"net/url"
	"strings"
)

// MessageHead 通用消息体结构
type MessageHead struct {
	ToUser  StringArray `json:"touser"`
	ToParty StringArray `json:"toparty"`
	ToTag   StringArray `json:"totag"`
	AgentID int         `json:"agentid"`
}

// MessageResponse 消息发送返回
type MessageResponse struct {
	Response
	InvalidUser  StringArray `json:"invaliduser,omitempty"`
	InvalidParty StringArray `json:"invalidparty,omitempty"`
	InvalidTag   StringArray `json:"invalidtag,omitempty"`
	err          error
}

// IsSuccess 是否测送成功
func (m *MessageResponse) IsSuccess() bool {
	return m.ErrCode == 0 && m.err == nil
}

// GetError 获取错误信息
func (m *MessageResponse) GetError() error {
	if m.IsSuccess() {
		return nil
	}
	return fmt.Errorf(m.ErrMessage+"%w", m.err)
}

// TextMessage 文本消息体
type TextMessage struct {
	MessageHead
	MessageType MessageType `json:"msgtype"`
	Safe        int8        `json:"safe"`
	Text        struct {
		Context string `json:"content"`
	} `json:"text"`
	EnableIDTrans int `json:"enable_id_trans"`
}

// PushTextMessage 发送文本消息
func (w *WorkWechat) PushTextMessage(head *MessageHead, context string, isSafe, isEnableIDTrans bool) *MessageResponse {
	message := TextMessage{
		MessageHead: *head,
		MessageType: MessageTypeText,
		Text: struct {
			Context string `json:"content"`
		}{Context: context},
	}
	if isSafe {
		message.Safe = 1
	}
	if isEnableIDTrans {
		message.EnableIDTrans = 1
	}
	return w.SendMessage(&message)
}

// TaskCardMessageUpdate 更新任务卡片消息状态
type TaskCardMessageUpdate struct {
	UserIDs    StringArray `json:"userids"`
	TaskID     string      `json:"task_id"`
	ClickedKey string      `json:"clicked_key"`
}

// UpdateTaskCard 更新卡片状态
func (w *WorkWechat) UpdateTaskCard(t *TaskCardMessageUpdate) *MessageResponse {
	message := struct {
		TaskCardMessageUpdate
		AgentID int64 `json:"agentid"`
	}{
		TaskCardMessageUpdate: *t,
		AgentID:               w.AgentID,
	}
	repBody := MessageResponse{}
	accessToken, err := w.GetAccessToken()
	if err != nil {
		Log.Error("accessToken 获取失败", err)
		repBody.err = err
		return &repBody
	}
	URL, err := url.Parse(fmt.Sprintf(updateTaskCard, accessToken))
	if err != nil {
		repBody.err = err
		return &repBody
	}
	s, err := json.Marshal(&message)
	if err != nil {
		Log.Error("消息序列化失败", err)
	}
	reqBody := strings.NewReader(string(s))
	body, err := w.httpPost(URL, reqBody)
	err = json.Unmarshal(body, &repBody)
	if err != nil {
		repBody.err = err
		return &repBody
	}
	return &repBody
}

// PushTaskCardMessage 推送任务卡片消息
func (w *WorkWechat) PushTaskCardMessage(head *MessageHead, messageBody *TaskCardMessageBody) *MessageResponse {
	message := TaskCardMessage{
		MessageHead: *head,
		MessageType: MessageTypeTaskCard,
		TaskCard:    messageBody,
	}
	return w.SendMessage(&message)
}

// TaskCardMessage 任务卡片消息
type TaskCardMessage struct {
	MessageHead
	MessageType MessageType          `json:"msgtype"`
	TaskCard    *TaskCardMessageBody `json:"taskcard"`
}

// TaskCardMessageBody 任务卡片消息内容
type TaskCardMessageBody struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	URL         string               `json:"url"`
	TaskID      string               `json:"task_id"`
	Btn         []TaskCardMessageBtn `json:"btn"`
}

// TaskCardMessageBtn 任务卡片消息按钮
type TaskCardMessageBtn struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	ReplaceName string `json:"replace_name"`
	Color       string `json:"color,omitempty"`
	IsBold      bool   `json:"is_bold,omitempty"`
}
