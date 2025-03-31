package service

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"auth_service/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	dbrepo    repository.DBAuthenticateRepos
	redisrepo repository.RedisSessionRepos
	validator *validator.Validate
}

func NewAuthService(repo repository.DBAuthenticateRepos, redis repository.RedisSessionRepos) *AuthService {
	validator := validator.New()
	return &AuthService{dbrepo: repo, validator: validator, redisrepo: redis}
}

func (as *AuthService) RegistrateAndLogin(user *model.Person, ctx context.Context) *ServiceResponse {
	registrateMap := make(map[string]error)
	var err error
	var tx *sql.Tx

	defer func() {
		if err != nil {
			if tx != nil {
				rbErr := as.dbrepo.RollbackTx(ctx, tx)
				if rbErr != nil {
					log.Printf("Error when rolling back a transaction: %v", rbErr)
				}
			} else {
				log.Println("Transaction is nil, rollback not needed")
			}
		}
	}()

	errorvalidate := validatePerson(as, user, true)
	if errorvalidate != nil {
		log.Printf("Validate error %v", errorvalidate)

		return &ServiceResponse{Success: false, Errors: errorvalidate}
	}

	tx, err = as.dbrepo.BeginTx(ctx)
	if err != nil {
		log.Printf("TransactionError %v", err)
		registrateMap["TransactionError"] = erro.ErrorStartTransaction
		return &ServiceResponse{Success: false, Errors: registrateMap}
	}

	hashpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("HashPassError %v", err)
		registrateMap["HashPassError"] = erro.ErrorHashPass
		return &ServiceResponse{Success: false, Errors: registrateMap}
	}
	user.Password = string(hashpass)
	userID := uuid.New()
	user.Id = userID
	response := as.dbrepo.CreateUser(ctx, user)

	if !response.Success {
		log.Printf("Error when creating person in the database %v", response.Errors)
		registrateMap["RegistrateError"] = response.Errors
		return &ServiceResponse{Success: false, Errors: registrateMap}
	}

	dbData, ok := response.Data.(repository.DBRepositoryResponseData)
	if !ok {
		log.Printf("Unexpected data type from repository: %T", response.Data)
		registrateMap["UnexpectedData"] = erro.ErrorUnexpectedData
		return &ServiceResponse{Success: false, Errors: registrateMap}
	}

	createdUserID := dbData.UserId

	sessionID := uuid.New().String()
	expirationTime := time.Now().Add(time.Hour * 24)
	duration := time.Until(expirationTime)
	session := model.Session{
		SessionID:      sessionID,
		UserID:         createdUserID,
		ExpirationTime: expirationTime,
	}

	redisResponse := as.redisrepo.SetSession(ctx, session, duration)
	if !redisResponse.Success {
		log.Printf("Error when creating a session in Redis: %v", redisResponse.Errors)
		registrateMap["SetSessionError"] = redisResponse.Errors
		return &ServiceResponse{Success: false, Errors: registrateMap}
	}

	err = as.dbrepo.CommitTx(ctx, tx)
	if err != nil {
		log.Printf("Transaction commit error: %v", err)
		registrateMap["CommitError"] = erro.ErrorCommitTransaction
		return &ServiceResponse{Success: false, Errors: registrateMap}
	}

	redisData, ok := redisResponse.Data.(repository.RedisRepositoryResponseData)
	if !ok {
		log.Printf("Unexpected data type from repository: %T", response.Data)
		registrateMap["UnexpectedData"] = erro.ErrorUnexpectedData
		return &ServiceResponse{Success: false, Errors: registrateMap}
	}
	log.Println("The session was created successfully and the user is registrated!")
	return &ServiceResponse{
		Success:        true,
		UserId:         dbData.UserId,
		SessionId:      redisData.SessionId,
		ExpirationTime: redisData.ExpirationTime,
	}
}
func (as *AuthService) AuthenticateAndLogin(user *model.Person, ctx context.Context) *ServiceResponse {
	errorvalidate := validatePerson(as, user, false)
	if errorvalidate != nil {
		log.Printf("Validate error %v", errorvalidate)
		return &ServiceResponse{Success: false, Errors: errorvalidate}
	}
	authenticateMap := make(map[string]error)
	response := as.dbrepo.GetUser(ctx, user.Email, user.Password)
	if !response.Success {
		log.Printf("Failed to authenticate user: %v", response.Errors)
		authenticateMap["AuthenticateError"] = response.Errors
		return &ServiceResponse{Success: false, Errors: authenticateMap}
	}

	dbData, ok := response.Data.(repository.DBRepositoryResponseData)
	if !ok {
		log.Printf("Unexpected data type from repository: %T", response.Data)
		authenticateMap["UnexpectedData"] = erro.ErrorUnexpectedData
		return &ServiceResponse{Success: false, Errors: authenticateMap}
	}

	userID := dbData.UserId

	sessionID := uuid.New().String()
	expirationTime := time.Now().Add(24 * time.Hour)
	session := model.Session{
		SessionID:      sessionID,
		UserID:         userID,
		ExpirationTime: expirationTime,
	}

	duration := time.Until(expirationTime)
	redisResponse := as.redisrepo.SetSession(ctx, session, duration)

	if !redisResponse.Success {
		log.Printf("Error when creating a session in Redis: %v", redisResponse.Errors)
		authenticateMap["SetSessionError"] = redisResponse.Errors
		return &ServiceResponse{Success: false, Errors: authenticateMap}
	}

	redisData, ok := redisResponse.Data.(repository.RedisRepositoryResponseData)
	if !ok {
		log.Printf("Unexpected data type from repository: %T", response.Data)
		authenticateMap["UnexpectedData"] = erro.ErrorUnexpectedData
		return &ServiceResponse{Success: false, Errors: authenticateMap}
	}
	log.Println("The session was created successfully and the user is authenticated!")
	return &ServiceResponse{
		Success:        true,
		UserId:         redisData.UserID,
		SessionId:      sessionID,
		ExpirationTime: expirationTime,
	}
}
func (as *AuthService) Authorization(ctx context.Context, sessionID string) *ServiceResponse {
	authorizateMap := make(map[string]error)
	repoResponse := as.redisrepo.GetSession(ctx, sessionID)
	if !repoResponse.Success {
		log.Printf("Error when getting a session from Redis: %v", repoResponse.Errors)
		authorizateMap["GetSessionError"] = repoResponse.Errors
		return &ServiceResponse{Success: false, Errors: authorizateMap}
	}

	redisData, ok := repoResponse.Data.(repository.RedisRepositoryResponseData)
	if !ok {
		log.Printf("Unexpected data type from repository: %T", repoResponse.Data)
		authorizateMap["UnexpectedData"] = erro.ErrorUnexpectedData
		return &ServiceResponse{Success: false, Errors: authorizateMap}
	}
	log.Println("The session has been confirmed and the user has successfully logged in")
	return &ServiceResponse{
		Success:        true,
		UserId:         redisData.UserID,
		SessionId:      redisData.SessionId,
		ExpirationTime: redisData.ExpirationTime,
	}
}
func validatePerson(as *AuthService, user *model.Person, flag bool) map[string]error {
	personToValidate := *user
	if !flag {
		personToValidate.Name = "qwertyuiopasdfghjklzxcvbn"
	}
	err := as.validator.Struct(&personToValidate)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			erors := make(map[string]error)
			for _, err := range validationErrors {

				switch err.Tag() {

				case "email":
					log.Println("Email format error")
					erors[err.Field()] = erro.ErrorNotEmail
				case "min":
					errv := fmt.Errorf("%s is too short", err.Field())
					log.Println(err.Field() + " format error")
					erors[err.Field()] = errv

				default:
					errv := fmt.Errorf("%s is Null", err.Field())
					log.Println(err.Field() + " format error")
					erors[err.Field()] = errv
				}
			}
			return erors
		}
	}
	return nil
}
