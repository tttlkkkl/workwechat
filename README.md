### 企业微信接口简单封装
因运维插件需要发送业务卡片消息，github上貌似没有找到类似代码库，故而进行简单封装,尚不支持服务商模式。其他功能接口按需增加。
如果有 gopher 需要增加接口支持可以给我提issues，有空帮加。
### 消息服务器使用示例：
```golang
package main

// 服务器交互测试
import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tttlkkkl/workwechat"
)

var router *gin.Engine
var chat *workwechat.WorkWechat

func init() {
	router = gin.Default()
	var err error
	chat, err = workwechat.NewWorkWechat(
		"corpID",
		"secret",
		1000034,
		workwechat.SetReceiveMessageAPI("token", "EncodingAESKey"),
	)
	if err != nil {
		log.Fatalln("初始化失败", err)
	}
	router.GET("/", func(c *gin.Context) {
		c.Data(200, "", []byte("hello go work wechat!"))
	})
	// 事件处理
	router.Any("/event", func(c *gin.Context) {
		m, err := chat.ReceiveMessage.Handle(c.Request, c.Writer)
		if err != nil {
			log.Println(err)
			c.Data(http.StatusInternalServerError, "", []byte("fail"))
		}
		messageType := workwechat.MessageType(m.Data.GetString("MsgType"))
		switch messageType {
		case workwechat.MessageTypeText:
		// todo
		case workwechat.MessageTypeEvent:
			eventType := workwechat.EventType(m.Data.GetString("Event"))
			switch eventType {
			case workwechat.EventTypeTaskCard:
				fmt.Println("收到任务卡片点击事件")
			}
		}
	})
}
func main() {
	svr := http.Server{
		Addr:    ":80",
		Handler: router,
	}
	if err := svr.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
```
### 支持的接口:
#### 消息推送
- 发送文本消息
- 发送任务卡片消息
- 接收事件消息