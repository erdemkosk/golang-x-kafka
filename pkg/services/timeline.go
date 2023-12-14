package services

import (
	"context"

	"github.com/go-redis/redis"
)

type TimelineService struct {
	rdb *redis.Client
}

func CreateTimelineService(rdb *redis.Client) *TimelineService {
	return &TimelineService{rdb: rdb}
}

func (t TimelineService) Push(ctx context.Context, user string, tweet ...interface{}) error {
	return t.rdb.WithContext(ctx).RPush("timeline:"+user, tweet...).Err()
}

func (t TimelineService) Latest(ctx context.Context, user string, count int64) ([]string, error) {
	return t.rdb.WithContext(ctx).LRange("timeline:"+user, -1*count, -1).Result()
}
