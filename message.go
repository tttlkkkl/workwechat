package workwechat

import (
	"errors"
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
	InvalidUser  StringArray `json:"invaliduser"`
	InvalidParty StringArray `json:"invalidparty"`
	InvalidTag   StringArray `json:"invalidtag"`
}

// TextMessage 文本消息体
type TextMessage struct {
	MessageHead
	MessageType string `json:"msgtype"`
	Safe        int8   `json:"safe"`
	Text        struct {
		Context string `json:"content"`
	} `json:"text"`
	EnableIDTrans int `json:"enable_id_trans"`
}

// PushTextMessage 发送文本消息
func (w *WorkWechat) PushTextMessage(head MessageHead, context string, isSafe, isEnableIDTrans bool) error {
	message := TextMessage{
		MessageHead: head,
		MessageType: "text",
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
	rep, err := w.SendMessage(&message)
	if err != nil {
		return err
	}
	if rep.ErrCode != 0 {
		return errors.New(rep.ErrMessage)
	}
	return nil
}

// PushTaskCardMessage 推送任务卡片消息
func (w *WorkWechat) PushTaskCardMessage(head MessageHead, messageBody *TaskCardMessage) error {
	message := struct {
		MessageHead
		MessageType string           `json:"msgtype"`
		TaskCard    *TaskCardMessage `json:"taskcard"`
	}{
		MessageHead: head,
		MessageType: "taskcard",
		TaskCard:    messageBody,
	}
	rep, err := w.SendMessage(&message)
	if err != nil {
		return err
	}
	if rep.ErrCode != 0 {
		return errors.New(rep.ErrMessage)
	}
	return nil
}

// TaskCardMessage 任务卡片消息内容
type TaskCardMessage struct {
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
