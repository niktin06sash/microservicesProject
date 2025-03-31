package repository

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthPostgres struct {
	Db *sql.DB
}

func (repoap *AuthPostgres) CreateUser(ctx context.Context, user *model.Person) *RepositoryResponse {
	var createdUserID uuid.UUID

	err := repoap.Db.QueryRowContext(ctx,
		"INSERT INTO UserZ (userid, username, useremail, userpassword) values ($1, $2, $3, $4) ON CONFLICT (useremail) DO NOTHING RETURNING userid;",
		user.Id, user.Name, user.Email, user.Password).Scan(&createdUserID)

	if err != nil {
		log.Printf("CreateUser Error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return &RepositoryResponse{Success: false, Errors: erro.ErrorUniqueEmail}
		}
		return &RepositoryResponse{Success: false, Errors: err}
	}

	responseData := DBRepositoryResponseData{
		UserId: createdUserID,
	}
	log.Println("Successful create person!")
	return &RepositoryResponse{Success: true, Data: responseData, Errors: nil}
}

func (repoap *AuthPostgres) GetUser(ctx context.Context, useremail, userpassword string) *RepositoryResponse {
	var hashpass string
	var userId uuid.UUID
	err := repoap.Db.QueryRowContext(ctx, "SELECT userid, userpassword FROM userZ WHERE useremail = $1", useremail).Scan(&userId, &hashpass)

	if err != nil {
		log.Printf("GetUser Error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return &RepositoryResponse{Success: false, Errors: erro.ErrorEmailNotRegister}
		}
		return &RepositoryResponse{Success: false, Errors: err}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(userpassword))
	if err != nil {
		log.Printf("CompareHashAndPassword Error: %v", err)
		return &RepositoryResponse{Success: false, Errors: erro.ErrorInvalidPassword}
	}

	responseData := DBRepositoryResponseData{
		UserId: userId,
	}
	log.Println("Successful get person!")
	return &RepositoryResponse{Success: true, Data: responseData, Errors: nil}
}
func (r *AuthPostgres) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.Db.BeginTx(ctx, nil)
}

func (r *AuthPostgres) RollbackTx(ctx context.Context, tx *sql.Tx) error {
	return tx.Rollback()
}

func (r *AuthPostgres) CommitTx(ctx context.Context, tx *sql.Tx) error {
	return tx.Commit()
}
func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{Db: db}
}
