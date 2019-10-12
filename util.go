package workwechat

import (
	"math/rand"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// GetRandomString 获取随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Values 类型处理
type Values map[string]interface{}

// Get 获取一个值
func (v *Values) Get(key string) (value interface{}, exists bool) {
	val := *v
	value, exists = val[key]
	return
}

// GetString 返回字符串
func (v *Values) GetString(key string) (s string) {
	if val, ok := v.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool 返回布尔值
func (v *Values) GetBool(key string) (b bool) {
	if val, ok := v.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt 返回int
func (v *Values) GetInt(key string) (i int) {
	if val, ok := v.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 返回int64
func (v *Values) GetInt64(key string) (i64 int64) {
	if val, ok := v.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 float64
func (v *Values) GetFloat64(key string) (f64 float64) {
	if val, ok := v.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime 返回 time
func (v *Values) GetTime(key string) (t time.Time) {
	if val, ok := v.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration 返回 duration
func (v *Values) GetDuration(key string) (d time.Duration) {
	if val, ok := v.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice 返回 []string
func (v *Values) GetStringSlice(key string) (ss []string) {
	if val, ok := v.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap 返回map[string]interface{}
func (v *Values) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := v.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString 返回 map[string]string
func (v *Values) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := v.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice 返回 map[string][]string
func (v *Values) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := v.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
