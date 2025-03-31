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
	type testCase struct {
		name            string
		user            *model.Person
		mockSetup       func(mock sqlmock.Sqlmock, user *model.Person)
		expectedSuccess bool
		expectedError   error
		checkData       func(t *testing.T, data interface{}, user *model.Person)
	}

	user := &model.Person{
		Id:       uuid.New(),
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "password",
	}

	testCases := []testCase{
		{
			name: "Successful CreateUser",
			user: user,
			mockSetup: func(mock sqlmock.Sqlmock, user *model.Person) {
				mock.ExpectQuery("INSERT INTO UserZ").
					WithArgs(user.Id, user.Name, user.Email, user.Password).
					WillReturnRows(sqlmock.NewRows([]string{"userid"}).AddRow(user.Id))
			},
			expectedSuccess: true,
			expectedError:   nil,
			checkData: func(t *testing.T, data interface{}, user *model.Person) {
				assert.NotNil(t, data, "Data должен быть не nil")
				dataCasted, ok := data.(DBRepositoryResponseData)
				assert.True(t, ok, "Data должен быть типа DBRepositoryResponseData")
				assert.Equal(t, user.Id, dataCasted.UserId, "UserId должен совпадать")
			},
		},
		{
			name: "Duplicate Email",
			user: user,
			mockSetup: func(mock sqlmock.Sqlmock, user *model.Person) {
				mock.ExpectQuery("INSERT INTO UserZ").
					WithArgs(user.Id, user.Name, user.Email, user.Password).
					WillReturnError(sql.ErrNoRows)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorUniqueEmail,
			checkData:       nil, // No data to check on error
		},
		{
			name: "General Error",
			user: user,
			mockSetup: func(mock sqlmock.Sqlmock, user *model.Person) {
				mock.ExpectQuery("INSERT INTO UserZ").
					WithArgs(user.Id, user.Name, user.Email, user.Password).
					WillReturnError(errors.New("general database error"))
			},
			expectedSuccess: false,
			expectedError:   errors.New("general database error"),
			checkData:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewAuthPostgres(db)

			tc.mockSetup(mock, tc.user)

			response := repo.CreateUser(context.Background(), tc.user)

			assert.Equal(t, tc.expectedSuccess, response.Success, "Success должен совпадать")
			if tc.expectedError != nil {
				assert.Error(t, response.Errors, "Должна быть ошибка")

				if tc.name != "General Error" {
					assert.Equal(t, tc.expectedError, response.Errors, "Тип ошибки должен совпадать")
				} else {
					assert.Equal(t, tc.expectedError.Error(), response.Errors.Error(), "Текст ошибки должен совпадать") // Проверяем текст ошибки
				}

			} else {
				assert.NoError(t, response.Errors, "Ошибки быть не должно")
			}

			if tc.checkData != nil {
				tc.checkData(t, response.Data, tc.user)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
func TestAuthPostgres_GetUser(t *testing.T) {
	type testCase struct {
		name            string
		useremail       string
		userpassword    string
		mockSetup       func(mock sqlmock.Sqlmock, useremail string)
		expectedSuccess bool
		expectedError   error
		checkData       func(t *testing.T, data interface{})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	userId := uuid.New()

	testCases := []testCase{
		{
			name:         "Successful GetUser",
			useremail:    "test@example.com",
			userpassword: "password",
			mockSetup: func(mock sqlmock.Sqlmock, useremail string) {
				mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail =").
					WithArgs(useremail).
					WillReturnRows(sqlmock.NewRows([]string{"userid", "userpassword"}).AddRow(userId, string(hashedPassword)))
			},
			expectedSuccess: true,
			expectedError:   nil,
			checkData: func(t *testing.T, data interface{}) {
				dataCasted, ok := data.(DBRepositoryResponseData)
				assert.True(t, ok, "Data должен быть типа DBRepositoryResponseData")
				assert.Equal(t, userId, dataCasted.UserId, "UserId должен совпадать")
			},
		},
		{
			name:         "Email Not Register",
			useremail:    "test@example.com",
			userpassword: "password",
			mockSetup: func(mock sqlmock.Sqlmock, useremail string) {
				mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail =").
					WithArgs(useremail).
					WillReturnError(sql.ErrNoRows)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorEmailNotRegister,
			checkData:       nil,
		},
		{
			name:         "Invalid Password",
			useremail:    "test@example.com",
			userpassword: "wrongpassword",
			mockSetup: func(mock sqlmock.Sqlmock, useremail string) {
				mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail =").
					WithArgs(useremail).
					WillReturnRows(sqlmock.NewRows([]string{"userid", "userpassword"}).AddRow(userId, string(hashedPassword)))
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorInvalidPassword,
			checkData:       nil,
		},
		{
			name:         "General DB Error",
			useremail:    "test@example.com",
			userpassword: "password",
			mockSetup: func(mock sqlmock.Sqlmock, useremail string) {
				mock.ExpectQuery("SELECT userid, userpassword FROM userZ WHERE useremail =").
					WithArgs(useremail).
					WillReturnError(errors.New("general database error"))
			},
			expectedSuccess: false,
			expectedError:   errors.New("general database error"),
			checkData:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewAuthPostgres(db)

			tc.mockSetup(mock, tc.useremail)

			response := repo.GetUser(context.Background(), tc.useremail, tc.userpassword)

			assert.Equal(t, tc.expectedSuccess, response.Success, "Success должен совпадать")

			if tc.expectedError != nil {
				assert.Error(t, response.Errors, "Должна быть ошибка")
				if tc.name != "General DB Error" {
					assert.Equal(t, tc.expectedError, response.Errors, "Тип ошибки должен совпадать")
				} else {
					assert.Equal(t, tc.expectedError.Error(), response.Errors.Error(), "Текст ошибки должен совпадать")
				}
			} else {
				assert.NoError(t, response.Errors, "Ошибки быть не должно")
			}

			if tc.checkData != nil {
				tc.checkData(t, response.Data)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
