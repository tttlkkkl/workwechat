package workwechat

import (
	"fmt"
	"net/url"
	"time"
)

// AccessTokenResponse accessToken 请求返回
type AccessTokenResponse struct {
	Response
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// GetAccessToken 获取accessToken
func (w *WorkWechat) GetAccessToken() (string, error) {
	token := w.Cache.Get(w.getAccessTokenCacheKey())
	if token != nil {
		accessToken, ok := token.(string)
		if ok {
			return accessToken, nil
		}
	}
	URL, err := url.Parse(fmt.Sprintf(accessTokenURL, w.CorpID, w.CorpSecret))
	if err != nil {
		return "", err
	}
	body, err := w.httpGet(URL)
	var rep AccessTokenResponse
	err = json.Unmarshal(body, &rep)
	if err != nil {
		return "", err
	}
	if rep.ErrCode != 0 {
		return "", err
	}
	rep.ExpiresIn = rep.ExpiresIn - 10
	w.Cache.Set(w.getAccessTokenCacheKey(), rep.AccessToken, time.Second*time.Duration(rep.ExpiresIn))
	return rep.AccessToken, nil
}

func (w *WorkWechat) getAccessTokenCacheKey() string {
	return fmt.Sprintf("at-%s", w.CorpSecret)
}
