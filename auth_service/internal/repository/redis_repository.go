package repository

import (
	"auth_service/internal/model"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthRedis struct {
	Client *redis.Client
}

func (redisrepo *AuthRedis) SetSession(ctx context.Context, session model.Session, expiration time.Duration) *RepositoryResponse {
	err := redisrepo.Client.HSet(ctx, session.SessionID, map[string]interface{}{
		"UserID":         session.UserID.String(),
		"ExpirationTime": session.ExpirationTime.Format(time.RFC3339),
	}).Err()

	if err != nil {
		log.Printf("Redis set error: %v", err)
		return &RepositoryResponse{Success: false, Errors: fmt.Errorf("failed to set session: %w", err)}
	}

	err = redisrepo.Client.Expire(ctx, session.SessionID, expiration).Err()
	if err != nil {
		log.Printf("Redis expire error: %v", err)
		return &RepositoryResponse{Success: false, Errors: fmt.Errorf("failed to set expiration: %w", err)}
	}

	responseData := RedisRepositoryResponseData{
		SessionId:      session.SessionID,
		ExpirationTime: session.ExpirationTime,
		UserID:         session.UserID,
	}

	return &RepositoryResponse{Success: true, Data: responseData, Errors: nil}
}

func (redisrepo *AuthRedis) GetSession(ctx context.Context, sessionID string) *RepositoryResponse {
	result, err := redisrepo.Client.HGetAll(ctx, sessionID).Result()
	if err != nil {
		log.Printf("Redis get error: %v", err)
		return &RepositoryResponse{Success: false, Errors: fmt.Errorf("failed to get session: %w", err)}
	}

	if len(result) == 0 {
		return &RepositoryResponse{Success: false, Errors: fmt.Errorf("session not found")}
	}

	userIDString, ok := result["UserID"]
	if !ok {
		return &RepositoryResponse{Success: false, Errors: fmt.Errorf("UserID not found in session")}
	}

	expirationTimeString, ok := result["ExpirationTime"]
	if !ok {
		return &RepositoryResponse{Success: false, Errors: fmt.Errorf("ExpirationTime not found in session")}
	}

	expirationTime, err := time.Parse(time.RFC3339, expirationTimeString)
	if err != nil {
		return &RepositoryResponse{Success: false, Errors: fmt.Errorf("failed to parse ExpirationTime: %w", err)}
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return &RepositoryResponse{Success: false, Errors: fmt.Errorf("failed to parse userID: %w", err)}
	}

	responseData := RedisRepositoryResponseData{
		SessionId:      sessionID,
		ExpirationTime: expirationTime,
		UserID:         userID,
	}
	return &RepositoryResponse{Success: true, Data: responseData, Errors: nil}
}

/*
	func (redisrepo *AuthRedis) DeleteSession(ctx context.Context, sessionID string) error {
		return nil
	}
*/
func NewAuthRedis(client *redis.Client) *AuthRedis {
	return &AuthRedis{Client: client}
}
