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
		"xx",
		"xxx",
		1000034,
		workwechat.SetReceiveMessageAPI("xxx", "xxx"),
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
