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
	RegistrateAndLogin(user *model.Person, ctx context.Context) *ServiceResponse
	AuthenticateAndLogin(user *model.Person, ctx context.Context) *ServiceResponse
	Authorization(ctx context.Context, sessionID string) *ServiceResponse
}
type Service struct {
	UserAuthentication
}
type ServiceResponse struct {
	Success        bool
	UserId         uuid.UUID
	SessionId      string
	ExpirationTime time.Time
	Errors         map[string]error
}

func NewService(repos *repository.Repository) *Service {

	return &Service{

		UserAuthentication: NewAuthService(repos.DBAuthenticateRepos, repos.RedisSessionRepos),
	}
}
