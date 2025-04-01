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
	CreateUser(ctx context.Context, user *model.Person) *RepositoryResponse
	GetUser(ctx context.Context, useremail, password string) *RepositoryResponse
	DeleteUser(ctx context.Context, userId uuid.UUID, password string) *RepositoryResponse
	BeginTx(ctx context.Context) (*sql.Tx, error)
	RollbackTx(ctx context.Context, tx *sql.Tx) error
	CommitTx(ctx context.Context, tx *sql.Tx) error
}
type RedisSessionRepos interface {
	SetSession(ctx context.Context, session model.Session, expiration time.Duration) *RepositoryResponse
	GetSession(ctx context.Context, sessionID string) *RepositoryResponse
	DeleteSession(ctx context.Context, sessionID string) *RepositoryResponse
}
type Repository struct {
	DBAuthenticateRepos
	RedisSessionRepos
}
type RepositoryResponse struct {
	Success bool
	Data    interface{}
	Errors  error
}

type DBRepositoryResponseData struct {
	UserId uuid.UUID
}

type RedisRepositoryResponseData struct {
	SessionId      string
	ExpirationTime time.Time
	UserID         uuid.UUID
}

func NewRepository(db *sql.DB, client *redis.Client) *Repository {
	return &Repository{
		DBAuthenticateRepos: NewAuthPostgres(db),
		RedisSessionRepos:   NewAuthRedis(client),
	}
}
