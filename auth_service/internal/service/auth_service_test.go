package service_test

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"auth_service/internal/repository"
	mock_repository "auth_service/internal/repository/mocks"
	"auth_service/internal/service"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Registrate(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockAuthorizationRepos, user *model.Person, ctx context.Context) *repository.AuthenticationRepositoryResponse

	testCases := []struct {
		name          string
		user          *model.Person
		mockBehavior  mockBehavior
		expected      *service.AuthenticationServiceResponse
		expectedError bool
	}{
		{
			name: "Success",
			user: &model.Person{
				Name:     "Valid Name",
				Email:    "valid@example.com",
				Password: "ValidPassword123",
			},
			mockBehavior: func(r *mock_repository.MockAuthorizationRepos, user *model.Person, ctx context.Context) *repository.AuthenticationRepositoryResponse {
				repoResponse := &repository.AuthenticationRepositoryResponse{
					Success: true,
					UserId:  uuid.New(),
					Errors:  nil,
				}
				r.EXPECT().CreateUser(ctx, gomock.Any()).Return(repoResponse)
				return repoResponse
			},
			expected: &service.AuthenticationServiceResponse{
				Success: true,
				Errors:  nil,
			},
			expectedError: false,
		},
		{
			name: "Invalid User Data",
			user: &model.Person{
				Name:     "",
				Email:    "invalid-email",
				Password: "short",
			},
			mockBehavior: func(r *mock_repository.MockAuthorizationRepos, user *model.Person, ctx context.Context) *repository.AuthenticationRepositoryResponse {
				return nil
			},
			expected: &service.AuthenticationServiceResponse{
				Success: false,
				Errors: map[string]error{
					"Email":    erro.ErrorNotEmail,
					"Name":     fmt.Errorf("%s is Null", "Name"),
					"Password": fmt.Errorf("%s is too short", "Password"),
				},
			},
			expectedError: true,
		},
		{
			name: "CreateUser Fails",
			user: &model.Person{
				Name:     "Valid Name",
				Email:    "valid@example.com",
				Password: "ValidPassword123",
			},
			mockBehavior: func(r *mock_repository.MockAuthorizationRepos, user *model.Person, ctx context.Context) *repository.AuthenticationRepositoryResponse {
				repoResponse := &repository.AuthenticationRepositoryResponse{
					Success: false,
					UserId:  uuid.Nil,
					Errors:  errors.New("CreateUser Failed"),
				}
				r.EXPECT().CreateUser(ctx, gomock.Any()).Return(repoResponse)
				return repoResponse
			},
			expected: &service.AuthenticationServiceResponse{
				Success: false,
				Errors:  map[string]error{"CreateError": errors.New("CreateUser Failed")},
			},
			expectedError: true,
		},
		/*{
			name: "Hash Pass fails",
			user: &model.Person{
				Name:     "Valid Name",
				Email:    "valid@example.com",
				Password: "ValidPassword123",
			},
			mockBehavior: func(r *mock_repository.MockAuthorizationRepos, user *model.Person, ctx context.Context) *repository.AuthenticationRepositoryResponse {
				return nil
			},
			expected: &service.AuthenticationServiceResponse{
				Success: false,
				Errors:  map[string]error{"HashPass": erro.ErrorHashPass},
			},
			expectedError: true,
		},*/
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Use mocks.NewMockAuthorizationRepos
			mockRepo := mock_repository.NewMockAuthorizationRepos(ctrl)

			repoResponse := tc.mockBehavior(mockRepo, tc.user, context.Background())

			authService := service.NewAuthService(mockRepo)
			result := authService.Registrate(tc.user, context.Background())

			if tc.expectedError {
				require.NotNil(t, result.Errors, "Expected errors, but got none")
				require.Equal(t, tc.expected.Success, result.Success, "Success flag mismatch")

				for field, expectedErr := range tc.expected.Errors {
					actualErr, ok := result.Errors[field]
					require.True(t, ok, "Expected error for field %s, but got none", field)
					require.EqualError(t, actualErr, expectedErr.Error(), "Error message mismatch for field %s", field)
				}

			} else {
				require.Nil(t, result.Errors, "Expected no errors, but got some")
				require.True(t, result.Success, "Expected success, but got failure")
				if repoResponse != nil {
					require.Equal(t, repoResponse.UserId, result.UserId, "UserID mismatch")
				}
			}
		})
	}
}

func TestAuthService_Authenticate(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockAuthorizationRepos, useremail, userpassword string, ctx context.Context) *repository.AuthenticationRepositoryResponse

	type args struct {
		UserEmail    string
		UserPassword string
	}

	testCases := []struct {
		name          string
		args          args
		mockBehavior  mockBehavior
		expected      *service.AuthenticationServiceResponse
		expectedError bool
	}{
		{
			name: "Success",
			args: args{
				UserEmail:    "testemail@gmail.com",
				UserPassword: "testpassword",
			},
			mockBehavior: func(r *mock_repository.MockAuthorizationRepos, useremail, userpassword string, ctx context.Context) *repository.AuthenticationRepositoryResponse {
				repoResponse := &repository.AuthenticationRepositoryResponse{
					Success: true,
					UserId:  uuid.New(),
					Errors:  nil,
				}
				r.EXPECT().GetUser(ctx, useremail, userpassword).Return(repoResponse) // Specify arguments
				return repoResponse
			},
			expected: &service.AuthenticationServiceResponse{
				Success: true,
				UserId:  uuid.Nil, // or expected UUID if known
				Errors:  nil,
			},
			expectedError: false,
		},
		{
			name: "Invalid User Data",
			args: args{
				UserEmail:    "testemailgmail.com",
				UserPassword: "tes",
			},
			mockBehavior: func(r *mock_repository.MockAuthorizationRepos, useremail, userpassword string, ctx context.Context) *repository.AuthenticationRepositoryResponse {
				return nil
			},
			expected: &service.AuthenticationServiceResponse{
				Success: false,
				Errors: map[string]error{
					"Email":    erro.ErrorNotEmail,
					"Password": fmt.Errorf("%s is too short", "Password"),
				},
			},
			expectedError: true,
		},
		{
			name: "Invalid Email",
			args: args{
				UserEmail:    "wrongtestemail@gmail.com",
				UserPassword: "testpassword",
			},
			mockBehavior: func(r *mock_repository.MockAuthorizationRepos, useremail, userpassword string, ctx context.Context) *repository.AuthenticationRepositoryResponse {
				repoResponse := &repository.AuthenticationRepositoryResponse{
					Success: false,
					UserId:  uuid.Nil,
					Errors:  erro.ErrorEmailNotRegister,
				}
				r.EXPECT().GetUser(ctx, useremail, userpassword).Return(repoResponse)
				return repoResponse
			},
			expected: &service.AuthenticationServiceResponse{
				Success: false,
				UserId:  uuid.Nil,
				Errors:  map[string]error{"AuthenticateError": erro.ErrorEmailNotRegister},
			},
			expectedError: true,
		},
		{
			name: "Invalid Password",
			args: args{
				UserEmail:    "invalid@example.com",
				UserPassword: "wrongpassword",
			},
			mockBehavior: func(r *mock_repository.MockAuthorizationRepos, useremail, userpassword string, ctx context.Context) *repository.AuthenticationRepositoryResponse {
				repoResponse := &repository.AuthenticationRepositoryResponse{
					Success: false,
					UserId:  uuid.Nil,
					Errors:  erro.ErrorInvalidPassword,
				}
				r.EXPECT().GetUser(ctx, useremail, userpassword).Return(repoResponse)
				return repoResponse
			},
			expected: &service.AuthenticationServiceResponse{
				Success: false,
				UserId:  uuid.Nil,
				Errors:  map[string]error{"AuthenticateError": erro.ErrorInvalidPassword},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockAuthorizationRepos(ctrl)
			tc.mockBehavior(mockRepo, tc.args.UserEmail, tc.args.UserPassword, context.Background())

			authService := service.NewAuthService(mockRepo)
			authResponse := authService.Authenticate(&model.Person{Email: tc.args.UserEmail, Password: tc.args.UserPassword}, context.Background())

			if tc.expectedError {
				require.NotNil(t, authResponse.Errors, "Expected errors, but got none")
				for key, expectedErr := range tc.expected.Errors {
					actualErr, ok := authResponse.Errors[key]
					require.True(t, ok, "Expected error for key %s", key)
					require.EqualError(t, actualErr, expectedErr.Error())

				}

			} else {
				require.Nil(t, authResponse.Errors, "Expected no errors, but got some")
				require.Equal(t, tc.expected.Success, authResponse.Success, "Success flag mismatch")
			}

		})
	}
}
