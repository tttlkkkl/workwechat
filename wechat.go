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
	CorpID     string
	CorpSecret string
	Cache      Cache
	HTTPClient *http.Client
}

// NewWorkWechat 客户端初始化
func NewWorkWechat(corpID, corpSecret string) *WorkWechat {
	return &WorkWechat{
		CorpID:     corpID,
		CorpSecret: corpSecret,
		Cache:      NewMemory(),
		HTTPClient: http.DefaultClient,
	}
}

// SendMessage 发送一条自定义结构的消息
func (w *WorkWechat) SendMessage(message interface{}) (*MessageResponse, error) {
	accessToken, err := w.GetAccessToken()
	if err != nil {
		loger.Error("accessToken 获取失败", err)
		return nil, err
	}
	URL, err := url.Parse(fmt.Sprintf(messageSendURL, accessToken))
	if err != nil {
		return nil, err
	}
	s, err := json.Marshal(message)
	if err != nil {
		loger.Error("消息序列化失败", err)
	}
	reqBody := strings.NewReader(string(s))
	body, err := w.httpPost(URL, reqBody)
	var repBosy MessageResponse
	err = json.Unmarshal(body, &repBosy)
	if err != nil {
		return nil, err
	}
	return &repBosy, nil
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
