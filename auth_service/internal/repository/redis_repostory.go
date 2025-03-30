package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRedis struct {
	Redis *redis.Client
}

func (client *AuthRedis) SetSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) *AuthenticationRepositoryResponse {
	return &AuthenticationRepositoryResponse{}
}
func (client *AuthRedis) GetSession(ctx context.Context, sessionID string) *AuthenticationRepositoryResponse {
	return &AuthenticationRepositoryResponse{}
}
func NewAuthRedis(client *redis.Client) *AuthRedis {
	return &AuthRedis{
		Redis: client,
	}
}
