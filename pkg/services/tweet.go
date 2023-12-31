package services

import (
	"context"

	"github.com/erdemkosk/golang-x-kafka/pkg/models"
	"github.com/erdemkosk/golang-x-kafka/pkg/repositories"
	"github.com/google/uuid"
)

type TweetService struct {
	redisRepo repositories.Redis[models.Tweet]
}

func CreateTweetService(redisInstance repositories.Redis[models.Tweet]) *TweetService {
	return &TweetService{redisRepo: redisInstance}
}

func (tweetService TweetService) Save(ctx context.Context, tweet models.Tweet) (models.Tweet, error) {
	tweet.UID = uuid.New().String()
	error := tweetService.redisRepo.Save(ctx, tweet)

	return tweet, error
}

func (tweetService TweetService) Get(ctx context.Context, uid string) (models.Tweet, error) {
	uid = "tweet:" + uid
	tweet, error := tweetService.redisRepo.Get(ctx, uid)

	return tweet, error
}

func (tweetService TweetService) MGet(ctx context.Context, uid ...string) ([]models.Tweet, error) {
	ids := make([]string, len(uid))
	for i, s := range uid {
		ids[i] = "tweet:" + s
	}
	return tweetService.redisRepo.MGet(ctx, ids...)
}
