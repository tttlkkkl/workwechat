package workwechat
import (
	log "github.com/sirupsen/logrus"
)

// Log teamkit 公共日志输出
var Log *log.Entry

func init() {
	// 打印文件等详细信息，这个对性能有影响
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	Log = log.WithField("pkg", "ko-socket")
}
