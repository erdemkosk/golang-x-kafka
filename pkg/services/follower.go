package services

import "context"

type FollowerService struct {
	followers map[string][]string
}

func CreateFollowerService() FollowerService {
	return FollowerService{
		followers: map[string][]string{
			"erdem": {"john", "maria", "hanna"},
			"kosk":  {"hanna"},
		},
	}
}

func (f FollowerService) Followers(ctx context.Context, user string) ([]string, error) {
	return f.followers[user], nil
}
