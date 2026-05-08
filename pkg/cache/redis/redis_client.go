package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

// BlacklistToken lưu token vào redis với TTL bằng thời gian còn lại của token
func (r *RedisClient) BlacklistToken(ctx context.Context, tokenID string, expiration time.Duration) error {
	return r.client.Set(ctx, "blacklist:"+tokenID, "true", expiration).Err()
}

// IsTokenBlacklisted kiểm tra xem token có trong danh sách đen không
func (r *RedisClient) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	val, err := r.client.Get(ctx, "blacklist:"+tokenID).Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == "true", err
}
