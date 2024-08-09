package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type LocalCodeCache struct {
	cache      *cache.Cache  // 本地缓存
	lock       sync.Mutex    // 锁
	expiration time.Duration // 过期时间
}

type codeItem struct {
	code   string    // 验证码
	cnt    int       // 可验证次数
	expire time.Time // 过期时间
}

func NewLocalCodeCache(c *cache.Cache, expiration time.Duration) CodeCache {
	return &LocalCodeCache{
		cache:      c,
		expiration: expiration,
	}
}

func (l *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now()
	key := l.key(biz, phone)
	// 先查看是否有验证码
	val, ok := l.cache.Get(key)
	if !ok {
		// 没有验证码，直接写
		l.cache.Set(key, codeItem{
			code:   code,
			cnt:    3,
			expire: now.Add(l.expiration),
		}, l.expiration)
	}

	// 有验证码，查看是否是在一分钟内多次发送
	item, ok := val.(codeItem)
	if !ok {
		return errors.New("系统错误")
	}
	if item.expire.Sub(now) > time.Minute*9 {
		// 一分钟内多次发送
		return ErrCodeSendTooMany
	}

	// 超过一分钟，重新发
	l.cache.Set(key, codeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expiration),
	}, l.expiration)
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return false, nil
}

func (l *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
