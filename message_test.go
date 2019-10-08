package workwechat

import (
	"testing"
)

func TestWorkWechat_PushTextMessage(t *testing.T) {
	client := NewWorkWechat(corpID, secret)
	type args struct {
		head            MessageHead
		context         string
		isSafe          bool
		isEnableIDTrans bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.PushTextMessage(tt.args.head, tt.args.context, tt.args.isSafe, tt.args.isEnableIDTrans); (err != nil) != tt.wantErr {
				t.Errorf("WorkWechat.PushTextMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkWechat_PushTaskCardMessage(t *testing.T) {
	client := NewWorkWechat(corpID, secret)
	type args struct {
		head        MessageHead
		messageBody *TaskCardMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test send task message",
			args: args{
				head: MessageHead{
					ToUser:  StringArray{"Hua", "lihua"},
					AgentID: 1000034,
				},
				messageBody: &TaskCardMessage{
					Title:       "赵明的礼物申请",
					Description: "礼品：A31茶具套装<br>用途：赠与小黑科技张总经理",
					URL:         "https://www.baidu.com",
					TaskID:      "111123",
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.PushTaskCardMessage(tt.args.head, tt.args.messageBody); (err != nil) != tt.wantErr {
				t.Errorf("WorkWechat.PushTaskCardMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
