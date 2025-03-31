package repository

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	t.Run("Successful CreateUser", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewAuthPostgres(db)

		user := &model.Person{
			Id:       uuid.New(),
			Name:     "testuser",
			Email:    "test@example.com",
			Password: "password",
		}

		mock.ExpectQuery("INSERT INTO UserZ").
			WithArgs(user.Id, user.Name, user.Email, user.Password).
			WillReturnRows(sqlmock.NewRows([]string{"userid"}).AddRow(user.Id))

		response := repo.CreateUser(context.Background(), user)

		assert.True(t, response.Success, "CreateUser должен вернуть Success = true")
		assert.Nil(t, response.Errors, "CreateUser должен вернуть Errors = nil")
		assert.NotNil(t, response.Data, "CreateUser должен вернуть Data != nil")

		data, ok := response.Data.(DBRepositoryResponseData)
		assert.True(t, ok, "Data должен быть типа DBRepositoryResponseData")
		assert.Equal(t, user.Id, data.UserId, "UserId должен совпадать")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Duplicate Email", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewAuthPostgres(db)

		user := &model.Person{
			Id:       uuid.New(),
			Name:     "testuser",
			Email:    "test@example.com",
			Password: "password",
		}

		mock.ExpectQuery("INSERT INTO UserZ").
			WithArgs(user.Id, user.Name, user.Email, user.Password).
			WillReturnError(sql.ErrNoRows)

		response := repo.CreateUser(context.Background(), user)

		assert.False(t, response.Success, "CreateUser должен вернуть Success = false")
		assert.NotNil(t, response.Errors, "CreateUser должен вернуть Errors != nil")
		assert.Equal(t, erro.ErrorUniqueEmail, response.Errors, "Должна быть ошибка: Unique Email")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("General Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewAuthPostgres(db)

		user := &model.Person{
			Id:       uuid.New(),
			Name:     "testuser",
			Email:    "test@example.com",
			Password: "password",
		}

		mock.ExpectQuery("INSERT INTO UserZ").
			WithArgs(user.Id, user.Name, user.Email, user.Password).
			WillReturnError(errors.New("general database error"))

		response := repo.CreateUser(context.Background(), user)

		assert.False(t, response.Success, "CreateUser должен вернуть Success = false")
		assert.NotNil(t, response.Errors, "CreateUser должен вернуть Errors != nil")
		assert.Equal(t, "general database error", response.Errors.Error(), "Текст ошибки должен совпадать") // Проверяем текст ошибки

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}

func TestAuthPostgres_GetUser(t *testing.T) {
	t.Run("Successful GetUser", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewAuthPostgres(db)

		useremail := "test@example.com"
		userpassword := "password"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userpassword), bcrypt.DefaultCost)

		userId := uuid.New()

		mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail =").
			WithArgs(useremail).
			WillReturnRows(sqlmock.NewRows([]string{"userid", "userpassword"}).AddRow(userId, string(hashedPassword)))

		response := repo.GetUser(context.Background(), useremail, userpassword)

		assert.True(t, response.Success, "GetUser должен вернуть Success = true")
		assert.Nil(t, response.Errors, "GetUser должен вернуть Errors = nil")
		assert.NotNil(t, response.Data, "GetUser должен вернуть Data != nil")

		data, ok := response.Data.(DBRepositoryResponseData)
		assert.True(t, ok, "Data должен быть типа DBRepositoryResponseData")
		assert.Equal(t, userId, data.UserId, "UserId должен совпадать")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Email Not Register", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewAuthPostgres(db)

		useremail := "test@example.com"
		userpassword := "password"

		mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail =").
			WithArgs(useremail).
			WillReturnError(sql.ErrNoRows)

		response := repo.GetUser(context.Background(), useremail, userpassword)

		assert.False(t, response.Success, "GetUser должен вернуть Success = false")
		assert.NotNil(t, response.Errors, "GetUser должен вернуть Errors != nil")
		assert.Equal(t, erro.ErrorEmailNotRegister, response.Errors, "Должна быть ошибка: Email Not Register")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewAuthPostgres(db)

		useremail := "test@example.com"
		userpassword := "wrongpassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

		userId := uuid.New()

		mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail =").
			WithArgs(useremail).
			WillReturnRows(sqlmock.NewRows([]string{"userid", "userpassword"}).AddRow(userId, string(hashedPassword)))

		response := repo.GetUser(context.Background(), useremail, userpassword)

		assert.False(t, response.Success, "GetUser должен вернуть Success = false")
		assert.NotNil(t, response.Errors, "GetUser должен вернуть Errors != nil")
		assert.Equal(t, erro.ErrorInvalidPassword, response.Errors, "Должна быть ошибка: Invalid Password")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("General DB Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewAuthPostgres(db)

		useremail := "test@example.com"
		userpassword := "password"

		mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail =").
			WithArgs(useremail).
			WillReturnError(errors.New("general database error"))

		response := repo.GetUser(context.Background(), useremail, userpassword)

		assert.False(t, response.Success, "GetUser должен вернуть Success = false")
		assert.NotNil(t, response.Errors, "GetUser должен вернуть Errors != nil")
		assert.Equal(t, "general database error", response.Errors.Error(), "Текст ошибки должен совпадать") // Проверяем текст ошибки

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}
