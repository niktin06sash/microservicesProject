package service

import (
	"auth_service/internal/model"
	"auth_service/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Authorization interface {
	Registrate(user *model.Person, ctx context.Context) *AuthenticationServiceResponse
	Authenticate(user *model.Person, ctx context.Context) *AuthenticationServiceResponse
	GenerateSession(userId uuid.UUID) (string, time.Time)
	Authorizate(sessionID string) *AuthenticationServiceResponse
	//DeleteSession(token string) error
}
type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {

	return &Service{
		Authorization: NewAuthService(repos.DBAuthenticateRepos),
	}
}
