package service

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"auth_service/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionService struct {
	repo repository.RedisSessionRepos
}

func NewSessionService(repo repository.RedisSessionRepos) *SessionService {
	return &SessionService{repo: repo}
}
func (s *SessionService) GenerateSession(ctx context.Context, userID uuid.UUID) *ServiceResponse {
	sessionID := uuid.New().String()
	expiration := time.Now().Add(time.Hour * 24)

	session := model.Session{
		SessionID:      sessionID,
		UserID:         userID,
		ExpirationTime: expiration,
	}

	duration := time.Until(expiration)

	repoResponse := s.repo.SetSession(ctx, session, duration)

	if !repoResponse.Success {

		errMap := map[string]error{"Redis": repoResponse.Errors}
		return &ServiceResponse{Success: false, Errors: errMap}
	}

	redisData, ok := repoResponse.Data.(repository.RedisRepositoryResponseData)
	if !ok {

		errMap := map[string]error{"Data": erro.ErrorUnexpectedData}
		return &ServiceResponse{Success: false, Errors: errMap}
	}

	return &ServiceResponse{Success: true, UserId: redisData.UserID}
}

func (s *SessionService) Authorizate(ctx context.Context, sessionID string) *ServiceResponse {
	repoResponse := s.repo.GetSession(ctx, sessionID)
	if !repoResponse.Success {

		errMap := map[string]error{"Redis": repoResponse.Errors}
		return &ServiceResponse{Success: false, Errors: errMap}
	}

	redisData, ok := repoResponse.Data.(repository.RedisRepositoryResponseData)
	if !ok {

		errMap := map[string]error{"Data": erro.ErrorUnexpectedData}
		return &ServiceResponse{Success: false, Errors: errMap}
	}

	return &ServiceResponse{Success: true, UserId: redisData.UserID}
}
