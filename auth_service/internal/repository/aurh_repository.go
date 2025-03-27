package repository

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthPostgres struct {
	Db *sql.DB
}
type AuthenticationRepositoryResponse struct {
	Success bool
	UserId  uuid.UUID
	Errors  error
}

func (repoap *AuthPostgres) CreateUser(ctx context.Context, user *model.Person) *AuthenticationRepositoryResponse {
	var createdUserID uuid.UUID

	err := repoap.Db.QueryRowContext(ctx,
		"INSERT INTO UserZ (userid, username, useremail, userpassword) values ($1, $2, $3, $4) ON CONFLICT (useremail) DO NOTHING RETURNING userid;",
		user.Id, user.Name, user.Email, user.Password).Scan(&createdUserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &AuthenticationRepositoryResponse{Success: false, Errors: erro.ErrorUniqueEmail}
		}
		return &AuthenticationRepositoryResponse{Success: false, Errors: err}
	}

	return &AuthenticationRepositoryResponse{Success: true, Errors: nil, UserId: createdUserID}
}
func (repoap *AuthPostgres) GetUser(ctx context.Context, useremail, userpassword string) *AuthenticationRepositoryResponse {
	var hashpass string
	var userId uuid.UUID
	err := repoap.Db.QueryRowContext(ctx, "SELECT userid, userpassword FROM userZ WHERE useremail = $1", useremail).Scan(&userId, &hashpass)

	if err == sql.ErrNoRows {
		return &AuthenticationRepositoryResponse{UserId: uuid.Nil, Success: false, Errors: erro.ErrorEmailNotRegister}
	}
	if err != nil {
		return &AuthenticationRepositoryResponse{UserId: uuid.Nil, Success: false, Errors: err}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(userpassword))
	if err != nil {

		return &AuthenticationRepositoryResponse{UserId: uuid.Nil, Success: false, Errors: erro.ErrorInvalidPassword}
	}

	return &AuthenticationRepositoryResponse{UserId: userId, Success: true, Errors: nil}
}
func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{Db: db}
}
