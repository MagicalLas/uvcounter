package application

import "github.com/redis/go-redis/v9"

type UVCounterService struct {
	RedisClient *redis.Client
}

func (s UVCounterService) GetUVCounter(id string) int64 {
	return 0
}
