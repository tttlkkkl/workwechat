package workwechat

import (
	"fmt"
	"testing"
)

func TestWorkWechat_PushTextMessage(t *testing.T) {
	type args struct {
		head            MessageHead
		context         string
		isSafe          bool
		isEnableIDTrans bool
	}
	tests := []struct {
		name    string
		args    args
		success bool
	}{
		{
			name: "send text message",
			args: args{
				head: MessageHead{
					ToUser:  StringArray{"Hua", "lihua"},
					AgentID: 1000034,
				},
				context: "你好，世界！",
			},
			success: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if m := client.PushTextMessage(&tt.args.head, tt.args.context, tt.args.isSafe, tt.args.isEnableIDTrans); m.IsSuccess() != tt.success {
				t.Errorf("WorkWechat.PushTextMessage() error = %v, success %v", m.GetError(), tt.success)
			}
		})
	}
}

func TestWorkWechat_PushTaskCardMessage(t *testing.T) {
	type args struct {
		head        MessageHead
		messageBody *TaskCardMessageBody
	}
	taskID := GetRandomString(10)
	fmt.Println("taskID:", taskID)
	tests := []struct {
		name    string
		args    args
		success bool
	}{
		{
			name: "test send task message",
			args: args{
				head: MessageHead{
					ToUser:  StringArray{"Hua", "lihua"},
					AgentID: 1000034,
				},
				messageBody: &TaskCardMessageBody{
					Title:       "===赵明的礼物申请===",
					Description: "礼品：A31茶具套装<br>用途：赠与小黑科技张总经理",
					URL:         "https://www.baidu.com",
					TaskID:      taskID,
					Btn: []TaskCardMessageBtn{
						{
							Key:         "key11",
							Name:        "批准",
							ReplaceName: "已批准",
							Color:       "blue",
							IsBold:      true,
						},
						{
							Key:         "key22",
							Name:        "驳回",
							ReplaceName: "已驳回",
							Color:       "blue",
							IsBold:      true,
						},
					},
				},
			},
			success: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if m := client.PushTaskCardMessage(&tt.args.head, tt.args.messageBody); m.IsSuccess() != tt.success {
				t.Errorf("WorkWechat.PushTaskCardMessage() error = %v, success %v", m.GetError(), tt.success)
			}
		})
	}
}

func TestWorkWechat_UpdateTaskCard(t *testing.T) {
	type args struct {
		t *TaskCardMessageUpdate
	}
	tests := []struct {
		name    string
		args    args
		success bool
	}{
		{
			name: "test btn1",
			args: args{
				t: &TaskCardMessageUpdate{
					UserIDs:    StringArray{"Hua"},
					TaskID:     "tq7t12nJIn",
					ClickedKey: "key11",
				},
			},
			success: true,
		},
		{
			name: "test btn2",
			args: args{
				t: &TaskCardMessageUpdate{
					UserIDs:    StringArray{"Hua"},
					TaskID:     "tq7t12nJIn",
					ClickedKey: "key22",
				},
			},
			success: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.UpdateTaskCard(tt.args.t); got.IsSuccess() != tt.success {
				t.Errorf("WorkWechat.UpdateTaskCard() = %v, want %v", got.GetError(), tt.success)
			}
		})
	}
}
