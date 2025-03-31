package repository

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {

	args := append([]interface{}{ctx, key}, values...)
	result := m.Called(args...)
	return result.Get(0).(*redis.IntCmd)
}
func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.MapStringStringCmd)
}

func TestMain(m *testing.M) {
	// setup
	log.SetOutput(os.Stdout)
	// run tests
	exitCode := m.Run()
	// teardown
	os.Exit(exitCode)
}

func TestAuthRedis_SetSession(t *testing.T) {
	type testCase struct {
		name            string
		session         model.Session
		expiration      time.Duration
		mockSetup       func(mock *MockRedisClient, session model.Session, expiration time.Duration)
		expectedSuccess bool
		expectedError   error
		checkData       func(t *testing.T, data interface{}, session model.Session)
	}

	session := model.Session{
		SessionID:      "session-id",
		UserID:         uuid.New(),
		ExpirationTime: time.Now().Add(time.Hour),
	}
	expiration := time.Hour

	testCases := []testCase{
		{
			name:       "Successful SetSession",
			session:    session,
			expiration: expiration,
			mockSetup: func(mocks *MockRedisClient, session model.Session, expiration time.Duration) {
				hSetCmd := redis.NewIntCmd(context.Background())
				hSetCmd.SetVal(1)

				mocks.On("HSet", context.Background(), session.SessionID, mock.AnythingOfType("map[string]interface{}")).Return(hSetCmd)

				expireCmd := redis.NewBoolCmd(context.Background())
				expireCmd.SetVal(true)
				mocks.On("Expire", context.Background(), session.SessionID, expiration).Return(expireCmd)
			},
			expectedSuccess: true,
			expectedError:   nil,
			checkData: func(t *testing.T, data interface{}, session model.Session) {
				dataCasted, ok := data.(RedisRepositoryResponseData)
				assert.True(t, ok, "Data должен быть типа RedisRepositoryResponseData")
				assert.Equal(t, session.SessionID, dataCasted.SessionId, "SessionId должен совпадать")
				assert.Equal(t, session.ExpirationTime.Format(time.RFC3339), dataCasted.ExpirationTime.Format(time.RFC3339), "ExpirationTime должен совпадать")
				assert.Equal(t, session.UserID, dataCasted.UserID, "UserID должен совпадать")
			},
		},
		{
			name:       "HSet Error",
			session:    session,
			expiration: expiration,
			mockSetup: func(mocks *MockRedisClient, session model.Session, expiration time.Duration) {
				hSetCmd := redis.NewIntCmd(context.Background())
				hSetCmd.SetErr(errors.New("hset error"))
				mocks.On("HSet", context.Background(), session.SessionID, mock.Anything).Return(hSetCmd)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorSetSession,
			checkData:       nil,
		},
		{
			name:       "Expire Error",
			session:    session,
			expiration: expiration,
			mockSetup: func(mocks *MockRedisClient, session model.Session, expiration time.Duration) {
				hSetCmd := redis.NewIntCmd(context.Background())
				hSetCmd.SetVal(1)
				mocks.On("HSet", context.Background(), session.SessionID, mock.Anything).Return(hSetCmd)

				expireCmd := redis.NewBoolCmd(context.Background())
				expireCmd.SetErr(errors.New("expire error"))
				mocks.On("Expire", context.Background(), session.SessionID, expiration, mock.Anything).Return(expireCmd)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorSetSession,
			checkData:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRedisClient := new(MockRedisClient)
			repo := &AuthRedis{Client: (*redis.Client)(unsafe.Pointer(mockRedisClient))}

			tc.mockSetup(mockRedisClient, tc.session, tc.expiration)

			response := repo.SetSession(context.Background(), tc.session, tc.expiration)

			assert.Equal(t, tc.expectedSuccess, response.Success, "Success должен совпадать")
			assert.Equal(t, tc.expectedError, response.Errors, "Тип ошибки должен совпадать")

			if tc.checkData != nil {
				tc.checkData(t, response.Data, tc.session)
			}

			mockRedisClient.AssertExpectations(t)
		})
	}
}

func TestAuthRedis_GetSession(t *testing.T) {
	type testCase struct {
		name            string
		sessionID       string
		mockSetup       func(mock *MockRedisClient, sessionID string)
		expectedSuccess bool
		expectedError   error
		checkData       func(t *testing.T, data interface{})
	}

	sessionID := "session-id"
	userID := uuid.New()
	expirationTime := time.Now().Add(time.Hour)

	testCases := []testCase{
		{
			name:      "Successful GetSession",
			sessionID: sessionID,
			mockSetup: func(mock *MockRedisClient, sessionID string) {
				hGetAllCmd := redis.NewMapStringStringCmd(context.Background())
				hGetAllCmd.SetVal(map[string]string{
					"UserID":         userID.String(),
					"ExpirationTime": expirationTime.Format(time.RFC3339),
				})
				mock.On("HGetAll", context.Background(), sessionID).Return(hGetAllCmd)
			},
			expectedSuccess: true,
			expectedError:   nil,
			checkData: func(t *testing.T, data interface{}) {
				dataCasted, ok := data.(RedisRepositoryResponseData)
				assert.True(t, ok, "Data должен быть типа RedisRepositoryResponseData")
				assert.Equal(t, sessionID, dataCasted.SessionId, "SessionId должен совпадать")
				assert.Equal(t, expirationTime.Format(time.RFC3339), dataCasted.ExpirationTime.Format(time.RFC3339), "ExpirationTime должен совпадать")
				assert.Equal(t, userID, dataCasted.UserID, "UserID должен совпадать")
			},
		},
		{
			name:      "HGetAll Error",
			sessionID: sessionID,
			mockSetup: func(mock *MockRedisClient, sessionID string) {
				hGetAllCmd := redis.NewMapStringStringCmd(context.Background())
				hGetAllCmd.SetErr(errors.New("hgetall error"))
				mock.On("HGetAll", context.Background(), sessionID).Return(hGetAllCmd)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorGetSession,
			checkData:       nil,
		},
		{
			name:      "Session Not Found",
			sessionID: sessionID,
			mockSetup: func(mock *MockRedisClient, sessionID string) {
				hGetAllCmd := redis.NewMapStringStringCmd(context.Background())
				hGetAllCmd.SetVal(map[string]string{})
				mock.On("HGetAll", context.Background(), sessionID).Return(hGetAllCmd)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorInvalidSessionID,
			checkData:       nil,
		},
		{
			name:      "UserID Not Found",
			sessionID: sessionID,
			mockSetup: func(mock *MockRedisClient, sessionID string) {
				hGetAllCmd := redis.NewMapStringStringCmd(context.Background())
				hGetAllCmd.SetVal(map[string]string{
					"ExpirationTime": expirationTime.Format(time.RFC3339),
				})
				mock.On("HGetAll", context.Background(), sessionID).Return(hGetAllCmd)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorGetUserIdSession,
			checkData:       nil,
		},
		{
			name:      "ExpirationTime Not Found",
			sessionID: sessionID,
			mockSetup: func(mock *MockRedisClient, sessionID string) {
				hGetAllCmd := redis.NewMapStringStringCmd(context.Background())
				hGetAllCmd.SetVal(map[string]string{
					"UserID": userID.String(),
				})
				mock.On("HGetAll", context.Background(), sessionID).Return(hGetAllCmd)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorGetExpirationTimeSession,
			checkData:       nil,
		},
		{
			name:      "Invalid UserID Format",
			sessionID: sessionID,
			mockSetup: func(mock *MockRedisClient, sessionID string) {
				hGetAllCmd := redis.NewMapStringStringCmd(context.Background())
				hGetAllCmd.SetVal(map[string]string{
					"UserID":         "invalid-uuid",
					"ExpirationTime": expirationTime.Format(time.RFC3339),
				})
				mock.On("HGetAll", context.Background(), sessionID).Return(hGetAllCmd)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorSessionParse,
			checkData:       nil,
		},
		{
			name:      "Invalid ExpirationTime Format",
			sessionID: sessionID,
			mockSetup: func(mock *MockRedisClient, sessionID string) {
				hGetAllCmd := redis.NewMapStringStringCmd(context.Background())
				hGetAllCmd.SetVal(map[string]string{
					"UserID":         userID.String(),
					"ExpirationTime": "invalid-time-format",
				})
				mock.On("HGetAll", context.Background(), sessionID).Return(hGetAllCmd)
			},
			expectedSuccess: false,
			expectedError:   erro.ErrorSessionParse,
			checkData:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRedisClient := new(MockRedisClient)
			repo := &AuthRedis{Client: (*redis.Client)(unsafe.Pointer(mockRedisClient))} // <-  небезопасный способ

			tc.mockSetup(mockRedisClient, tc.sessionID)

			response := repo.GetSession(context.Background(), tc.sessionID)

			assert.Equal(t, tc.expectedSuccess, response.Success, "Success должен совпадать")
			assert.Equal(t, tc.expectedError, response.Errors, "Тип ошибки должен совпадать")

			if tc.checkData != nil {
				tc.checkData(t, response.Data)
			}

			mockRedisClient.AssertExpectations(t)
		})
	}
}
