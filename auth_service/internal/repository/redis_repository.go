package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRedis struct {
	Client *redis.Client
}

func (redisrepo *AuthRedis) SetSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) *RepositoryResponse {
	return &RepositoryResponse{}
}
func (redisrepo *AuthRedis) GetSession(ctx context.Context, sessionID string) *RepositoryResponse {
	return &RepositoryResponse{}
}

/*
	func (redisrepo *AuthRedis) DeleteSession(ctx context.Context, sessionID string) error {
		return nil
	}
*/
func NewAuthRedis(client *redis.Client) *AuthRedis {
	return &AuthRedis{Client: client}
}
