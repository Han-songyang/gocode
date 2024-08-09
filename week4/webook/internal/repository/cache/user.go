package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/domain"
)

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *UserCache) Get(ctx *gin.Context, uid int64) (domain.User, error) {
	key := c.key(uid)
	fmt.Println("key:", key)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var du domain.User
	err = json.Unmarshal([]byte(data), &du)
	return du, err
}

func (c *UserCache) Set(ctx *gin.Context, du domain.User) error {
	key := c.key(du.Id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *UserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
