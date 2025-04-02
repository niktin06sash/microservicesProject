package service

import (
	"auth_service/internal/kafka"
	"auth_service/internal/model"
	"auth_service/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type UserAuthentication interface {
	RegistrateAndLogin(ctx context.Context, user *model.Person) *ServiceResponse
	AuthenticateAndLogin(ctx context.Context, user *model.Person) *ServiceResponse
	Authorization(ctx context.Context, sessionID string) *ServiceResponse
	Logout(ctx context.Context, sessionID string, userId uuid.UUID) *ServiceResponse
	DeleteAccount(ctx context.Context, sessionID string, userid uuid.UUID, password string) *ServiceResponse
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

func NewService(repos *repository.Repository, kafkaProd kafka.KafkaProducer) *Service {

	return &Service{

		UserAuthentication: NewAuthService(repos.DBAuthenticateRepos, repos.RedisSessionRepos, kafkaProd),
	}
}
