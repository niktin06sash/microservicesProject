package service

import (
	"auth_service/internal/model"
	"auth_service/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type UserAuthentication interface {
	Registrate(user *model.Person, ctx context.Context) *ServiceResponse
	Authenticate(user *model.Person, ctx context.Context) *ServiceResponse
}
type SessionManager interface {
	GenerateSession(ctx context.Context, userId uuid.UUID) *ServiceResponse
	Authorizate(ctx context.Context, sessionID string) *ServiceResponse
}
type Service struct {
	UserAuthentication
	SessionManager
}
type ServiceResponse struct {
	Success    bool
	UserId     uuid.UUID
	Errors     map[string]error
	Expiration *time.Duration
}

func NewService(repos *repository.Repository) *Service {

	return &Service{
<<<<<<< HEAD
		UserAuthentication: NewAuthService(repos.DBAuthenticateRepos),
		SessionManager:     NewSessionService(repos.RedisSessionRepos),
=======
		Authorization: NewAuthService(repos.DBAuthenticateRepos),
>>>>>>> new_branch
	}
}
