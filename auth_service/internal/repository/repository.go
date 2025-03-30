package repository

import (
	"auth_service/internal/model"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type DBAuthenticateRepos interface {
<<<<<<< HEAD
	CreateUser(ctx context.Context, user *model.Person) *RepositoryResponse
	GetUser(ctx context.Context, useremail, password string) *RepositoryResponse
}
type RedisSessionRepos interface {
	SetSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) *RepositoryResponse
	GetSession(ctx context.Context, sessionID string) *RepositoryResponse
	//DeleteSession(ctx context.Context, sessionID string) error
}
type Repository struct {
	DBAuthenticateRepos
	RedisSessionRepos
}
type RepositoryResponse struct {
	Success bool
	UserId  uuid.UUID
	Errors  error
	Time    *time.Duration
}

=======
	CreateUser(ctx context.Context, user *model.Person) *AuthenticationRepositoryResponse
	GetUser(ctx context.Context, useremail, password string) *AuthenticationRepositoryResponse
}
type RedisSessionRepos interface {
	SetSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) *AuthenticationRepositoryResponse
	GetSession(ctx context.Context, sessionID string) *AuthenticationRepositoryResponse
	//DeleteSession(ctx context.Context, sessionID string) error
}
type AuthenticationRepositoryResponse struct {
	Success bool
	UserId  uuid.UUID
	Errors  error
}

type Repository struct {
	DBAuthenticateRepos
	RedisSessionRepos
}

>>>>>>> new_branch
func NewRepository(db *sql.DB, client *redis.Client) *Repository {
	return &Repository{
		DBAuthenticateRepos: NewAuthPostgres(db),
		RedisSessionRepos:   NewAuthRedis(client),
	}
}
