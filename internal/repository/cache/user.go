package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ra1n6ow/webook/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	cmd    redis.Cmdable
	expire time.Duration
}

func NewUserCache(cmd redis.Cmdable, expire time.Duration) *UserCache {
	return &UserCache{
		cmd:    cmd,
		expire: expire,
	}
}

func (c *UserCache) Key(id int64) string {
	return fmt.Sprintf("users:info:%d", id)
}

func (c *UserCache) Set(ctx context.Context, u domain.User) error {
	bytes, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := c.Key(u.Id)
	return c.cmd.Set(ctx, key, bytes, c.expire).Err()
}

func (c *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := c.Key(id)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal([]byte(data), &user)
	return user, err
}
