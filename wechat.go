package workwechat // import "github.com/tttlkkkl/workwechat"

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// WorkWechat 企业微信接口操作对象
type WorkWechat struct {
	Options
	CorpID         string
	CorpSecret     string
	AgentID        int64
	ReceiveMessage *ReceiveMessage
	Cache          Cache
	HTTPClient     *http.Client
}

// Options 配置选项
type Options struct {
	// ReceiveMessageAPI 接收消息服务器配置
	ReceiveMessageAPI struct {
		Token          string
		EncodingAESKey string
	}
}

// Option 设置配置选项
type Option func(*Options)

// SetReceiveMessageAPI 设置消息接收服务器
func SetReceiveMessageAPI(token, encodingAESKey string) Option {
	return func(o *Options) {
		o.ReceiveMessageAPI.Token = token
		o.ReceiveMessageAPI.EncodingAESKey = encodingAESKey
	}
}

// NewWorkWechat 客户端初始化
func NewWorkWechat(corpID, corpSecret string, agentID int64, opt ...Option) (*WorkWechat, error) {
	opts := new(Options)
	for _, o := range opt {
		o(opts)
	}
	var err error
	rm := &ReceiveMessage{}
	if opts.ReceiveMessageAPI.Token != "" && opts.ReceiveMessageAPI.EncodingAESKey != "" {
		rm, err = newReceiveMessage(opts.ReceiveMessageAPI.Token, opts.ReceiveMessageAPI.EncodingAESKey)
		if err != nil {
			return nil, err
		}
	}
	return &WorkWechat{
		Options:        *opts,
		CorpID:         corpID,
		CorpSecret:     corpSecret,
		AgentID:        agentID,
		ReceiveMessage: rm,
		Cache:          NewMemory(),
		HTTPClient:     http.DefaultClient,
	}, nil
}

// SendMessage 发送一条自定义结构的消息
func (w *WorkWechat) SendMessage(message interface{}) *MessageResponse {
	repBody := MessageResponse{}
	accessToken, err := w.GetAccessToken()
	if err != nil {
		loger.Error("accessToken 获取失败", err)
		repBody.err = err
		return &repBody
	}
	URL, err := url.Parse(fmt.Sprintf(messageSendURL, accessToken))
	if err != nil {
		repBody.err = err
		return &repBody
	}
	s, err := json.Marshal(message)
	if err != nil {
		loger.Error("消息序列化失败", err)
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
func (w *WorkWechat) httpGet(URL *url.URL) ([]byte, error) {
	rep, err := w.HTTPClient.Get(URL.String())
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(rep.Body)
}

func (w *WorkWechat) httpPost(URL *url.URL, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", URL.String(), body)
	if err != nil {
		return nil, err
	}
	rep, err := w.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rep.Body.Close()
	return ioutil.ReadAll(rep.Body)
}
