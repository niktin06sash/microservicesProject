package repository

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"context"
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
		log.Printf("Hset error: %v", err)
		return &RepositoryResponse{Success: false, Errors: erro.ErrorSetSession}
	}

	err = redisrepo.Client.Expire(ctx, session.SessionID, expiration).Err()
	if err != nil {
		log.Printf("Expire error: %v", err)
		return &RepositoryResponse{Success: false, Errors: erro.ErrorSetSession}
	}

	responseData := RedisRepositoryResponseData{
		SessionId:      session.SessionID,
		ExpirationTime: session.ExpirationTime,
		UserID:         session.UserID,
	}
	log.Println("Successful session installation!")
	return &RepositoryResponse{Success: true, Data: responseData, Errors: nil}
}

func (redisrepo *AuthRedis) GetSession(ctx context.Context, sessionID string) *RepositoryResponse {
	result, err := redisrepo.Client.HGetAll(ctx, sessionID).Result()
	if err != nil {
		log.Printf("HGetAll error: %v", err)
		return &RepositoryResponse{Success: false, Errors: erro.ErrorGetSession}
	}

	if len(result) == 0 {
		return &RepositoryResponse{Success: false, Errors: erro.ErrorInvalidSessionID}
	}

	userIDString, ok := result["UserID"]
	if !ok {
		return &RepositoryResponse{Success: false, Errors: erro.ErrorGetUserIdSession}
	}

	expirationTimeString, ok := result["ExpirationTime"]
	if !ok {
		return &RepositoryResponse{Success: false, Errors: erro.ErrorGetExpirationTimeSession}
	}

	expirationTime, err := time.Parse(time.RFC3339, expirationTimeString)
	if err != nil {
		log.Printf("Time-parse error: %v", err)
		return &RepositoryResponse{Success: false, Errors: erro.ErrorSessionParse}
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		log.Printf("UUID-parse error: %v", err)
		return &RepositoryResponse{Success: false, Errors: erro.ErrorSessionParse}
	}

	responseData := RedisRepositoryResponseData{
		SessionId:      sessionID,
		ExpirationTime: expirationTime,
		UserID:         userID,
	}
	log.Println("Successful session receiving!")
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
