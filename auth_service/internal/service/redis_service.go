package service

import (
	"auth_service/internal/repository"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type SessionService struct {
	repo repository.RedisSessionRepos
}
type Session struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func NewSessionService(repo repository.RedisSessionRepos) *SessionService {
	return &SessionService{repo: repo}
}
func (s *SessionService) GenerateSession(ctx context.Context, userID uuid.UUID) *ServiceResponse {
	sessionID := uuid.New().String()
	expiration := time.Now().Add(time.Hour * 24)

	duration := time.Until(expiration)

	repoResponse := s.repo.SetSession(ctx, sessionID, userID.String(), duration)
	if !repoResponse.Success {
		log.Printf("Ошибка при сохранении сессии в Redis: %v", repoResponse.Errors)
		return &ServiceResponse{}
	}

	return &ServiceResponse{}
}

func (s *SessionService) Authorizate(ctx context.Context, sessionID string) *ServiceResponse {
	repoResponse := s.repo.GetSession(ctx, sessionID)
	if !repoResponse.Success {
		log.Printf("Ошибка при получении сессии из Redis: %v", repoResponse.Errors)
		return &ServiceResponse{}
	}

	return &ServiceResponse{}
}
