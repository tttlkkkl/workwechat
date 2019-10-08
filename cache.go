package workwechat

import (
	"sync"
	"time"
)

// Cache AccessToken cache
type Cache interface {
	Get(key string) interface{}
	Set(key string, val interface{}, timeout time.Duration) error
	IsExist(key string) bool
	Del(key string) error
}

// Memory 内存缓存实现
type Memory struct {
	sync.Mutex

	data map[string]*data
}

type data struct {
	Data    interface{}
	Expired time.Time
}

// NewMemory create new memcache
func NewMemory() Cache {
	return &Memory{
		data: map[string]*data{},
	}
}

// Get return cached value
func (mem *Memory) Get(key string) interface{} {
	if ret, ok := mem.data[key]; ok {
		if ret.Expired.Before(time.Now()) {
			mem.deleteKey(key)
			return nil
		}
		return ret.Data
	}
	return nil
}

// Set cached value with key and expire time.
func (mem *Memory) Set(key string, val interface{}, timeout time.Duration) (err error) {
	mem.Lock()
	defer mem.Unlock()

	mem.data[key] = &data{
		Data:    val,
		Expired: time.Now().Add(timeout),
	}
	return nil
}

// IsExist check value exists in memcache.
func (mem *Memory) IsExist(key string) bool {
	if ret, ok := mem.data[key]; ok {
		if ret.Expired.Before(time.Now()) {
			mem.deleteKey(key)
			return false
		}
		return true
	}
	return false
}

// Del 删除
func (mem *Memory) Del(key string) error {
	return mem.deleteKey(key)
}

func (mem *Memory) deleteKey(key string) error {
	mem.Lock()
	defer mem.Unlock()
	delete(mem.data, key)
	return nil
}
