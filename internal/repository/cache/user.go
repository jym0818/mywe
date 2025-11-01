package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jym0818/mywe/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	Set(ctx context.Context, user domain.User) error
	Get(ctx context.Context, uid int64) (domain.User, error)
}

type userCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (u *userCache) Set(ctx context.Context, user domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := u.key(user.Id)
	return u.cmd.Set(ctx, key, data, u.expiration).Err()
}

func (u *userCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	data, err := u.cmd.Get(ctx, u.key(uid)).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *userCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func NewuserCache(cmd redis.Cmdable) UserCache {
	return &userCache{
		cmd:        cmd,
		expiration: time.Minute * 30,
	}
}
