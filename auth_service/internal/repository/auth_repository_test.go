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
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	repos := NewAuthPostgres(db)

	testCases := []struct {
		name        string
		mock        func() uuid.UUID
		user        *model.Person
		expected    AuthenticationRepositoryResponse
		expectederr error
	}{
		{
			name: "success",
			mock: func() uuid.UUID {
				expectedUserID := uuid.New()
				rows := sqlmock.NewRows([]string{"userid"}).AddRow(expectedUserID)
				mock.ExpectQuery("INSERT INTO UserZ").
					WithArgs(sqlmock.AnyArg(), "testuser", "test@example.com", "securePassword123").
					WillReturnRows(rows)
				return expectedUserID
			},
			user: &model.Person{
				Id:       uuid.New(),
				Name:     "testuser",
				Email:    "test@example.com",
				Password: "securePassword123",
			},
			expected: AuthenticationRepositoryResponse{
				Success: true,
				Errors:  nil,
			},
			expectederr: nil,
		},
		{
			name: "duplicate email",
			mock: func() uuid.UUID {
				mock.ExpectQuery("INSERT INTO UserZ").
					WithArgs(sqlmock.AnyArg(), "testuser2", "test@example.com", "securePassword123").
					WillReturnError(sql.ErrNoRows)
				return uuid.Nil
			},
			user: &model.Person{
				Id:       uuid.New(),
				Name:     "testuser2",
				Email:    "test@example.com",
				Password: "securePassword123",
			},
			expected: AuthenticationRepositoryResponse{
				Success: false,
				Errors:  erro.ErrorUniqueEmail,
				UserId:  uuid.Nil,
			},
			expectederr: nil,
		},
		{name: "db error",
			mock: func() uuid.UUID {
				mock.ExpectQuery("INSERT INTO UserZ").
					WithArgs(sqlmock.AnyArg(), "testuser2", "test@example.com", "securePassword123").
					WillReturnError(errors.New("database error"))
				return uuid.Nil
			},
			user: &model.Person{
				Id:       uuid.New(),
				Name:     "testuser2",
				Email:    "test@example.com",
				Password: "securePassword123",
			},
			expected: AuthenticationRepositoryResponse{
				Success: false,
				Errors:  errors.New("database error"),
				UserId:  uuid.Nil,
			},
			expectederr: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedUserID := tc.mock()
			result := repos.CreateUser(context.Background(), tc.user)

			require.Equal(t, tc.expected.Success, result.Success)
			require.Equal(t, tc.expected.Errors, result.Errors)

			if tc.name == "success" {
				require.Equal(t, expectedUserID, result.UserId)
			} else {
				require.Equal(t, uuid.Nil, result.UserId)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repos := NewAuthPostgres(db)

	testCases := []struct {
		name         string
		mock         func(mock sqlmock.Sqlmock) uuid.UUID
		userEmail    string
		userPassword string
		expected     AuthenticationRepositoryResponse
		expectedErr  error
	}{
		{
			name: "success",
			mock: func(mock sqlmock.Sqlmock) uuid.UUID {
				userId := uuid.New()
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securePassword123"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{"userid", "userpassword"}).AddRow(userId, hashedPassword)

				mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
				return userId
			},
			userEmail:    "test@example.com",
			userPassword: "securePassword123",
			expected: AuthenticationRepositoryResponse{
				Success: true,
				Errors:  nil,
				UserId:  uuid.Nil,
			},
			expectedErr: nil,
		},
		{
			name: "email not registered",
			mock: func(mock sqlmock.Sqlmock) uuid.UUID {
				mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail = \\$1").
					WithArgs("test@example.com").
					WillReturnError(sql.ErrNoRows)
				return uuid.Nil
			},
			userEmail:    "test@example.com",
			userPassword: "securePassword123",
			expected: AuthenticationRepositoryResponse{
				Success: false,
				Errors:  erro.ErrorEmailNotRegister,
				UserId:  uuid.Nil,
			},
			expectedErr: nil,
		},
		{
			name: "invalid password",
			mock: func(mock sqlmock.Sqlmock) uuid.UUID {
				userId := uuid.New()
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securePassword123"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{"userid", "userpassword"}).AddRow(userId, hashedPassword)

				mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
				return userId
			},
			userEmail:    "test@example.com",
			userPassword: "wrongPassword",
			expected: AuthenticationRepositoryResponse{
				Success: false,
				Errors:  erro.ErrorInvalidPassword,
				UserId:  uuid.Nil,
			},
			expectedErr: nil,
		},
		{
			name: "db error",
			mock: func(mock sqlmock.Sqlmock) uuid.UUID {
				mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail = \\$1").
					WithArgs("test@example.com").
					WillReturnError(errors.New("database error"))
				return uuid.Nil
			},
			userEmail:    "test@example.com",
			userPassword: "securePassword123",
			expected: AuthenticationRepositoryResponse{
				Success: false,
				Errors:  errors.New("database error"),
				UserId:  uuid.Nil,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedUserID := tc.mock(mock)
			result := repos.GetUser(context.Background(), tc.userEmail, tc.userPassword)

			require.Equal(t, tc.expected.Success, result.Success)
			require.Equal(t, tc.expected.Errors, result.Errors)
			if tc.name == "success" {
				require.NotEqual(t, uuid.Nil, result.UserId)
				require.Equal(t, expectedUserID, result.UserId)
			} else {
				require.Equal(t, tc.expected.UserId, result.UserId)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
