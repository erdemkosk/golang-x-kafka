package repositories

import (
	"context"
	"encoding/json"

	"github.com/erdemkosk/golang-x-kafka/pkg/models"
	"github.com/go-redis/redis"
)

type Redis[T models.Keyer] struct { //Burada T bu interface implemente etmiş (yani içindeki fonksiyonu imp eden fonksiyonu ) olan struct bekliyor
	rdb *redis.Client
}

func NewRedis[T models.Keyer](rdb *redis.Client) Redis[T] {
	r := Redis[T]{rdb: rdb}
	return r
}

func (r Redis[T]) Save(ctx context.Context, k T) error {
	b, _ := json.Marshal(k)
	return r.rdb.WithContext(ctx).Set(k.Key(), b, 0).Err()
}

func (r Redis[T]) Get(ctx context.Context, key string) (T, error) {
	var t T
	b, err := r.rdb.Get(key).Bytes()
	if err != nil {
		return t, err
	}
	json.Unmarshal(b, &t)
	return t, nil
}

func (r Redis[T]) MGet(ctx context.Context, key ...string) ([]T, error) {
	bb, err := r.rdb.WithContext(ctx).MGet(key...).Result()
	if err != nil {
		return nil, err
	}
	result := make([]T, len(key))
	for i, b := range bb {
		json.Unmarshal([]byte(b.(string)), &result[i])
	}
	return result, nil
}
