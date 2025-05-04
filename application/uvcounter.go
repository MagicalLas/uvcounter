package application

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UVCounterService struct {
	RedisClient *redis.Client
}

func (s *UVCounterService) GetUVCounter(id string) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.RedisClient.Get(ctx, fmt.Sprintf("sampleKey:%s", id))

	return 0
}
