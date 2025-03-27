package repository

import (
	"auth_service/internal/model"
	"context"
	"database/sql"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type AuthorizationRepos interface {
	CreateUser(ctx context.Context, user *model.Person) *AuthenticationRepositoryResponse
	GetUser(ctx context.Context, useremail, password string) *AuthenticationRepositoryResponse
}
type Repository struct {
	AuthorizationRepos
}

func NewAuthRepository(db *sql.DB) *Repository {
	return &Repository{
		AuthorizationRepos: NewAuthPostgres(db),
	}
}
