package service

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"auth_service/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.DBAuthenticateRepos
	validator *validator.Validate
}

func NewAuthService(repo repository.DBAuthenticateRepos) *AuthService {
	validator := validator.New()
	return &AuthService{repo: repo, validator: validator}
}

type AuthenticationServiceResponse struct {
	Success bool
	UserId  uuid.UUID
	Errors  map[string]error
}

func (as *AuthService) Registrate(user *model.Person, ctx context.Context) *AuthenticationServiceResponse {
	errorvalidate := validatePerson(as, user, true)
	if errorvalidate != nil {
		return &AuthenticationServiceResponse{Success: false, Errors: errorvalidate}
	}
	hashpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		hashmapa := make(map[string]error)
		hashmapa["HashPass"] = erro.ErrorHashPass
		return &AuthenticationServiceResponse{Success: false, Errors: hashmapa}
	}
	user.Password = string(hashpass)
	userID := uuid.New()
	user.Id = userID
	response := as.repo.CreateUser(ctx, user)
	if !response.Success {
		errmapa := make(map[string]error)
		errmapa["CreateError"] = response.Errors
		return &AuthenticationServiceResponse{Success: false, Errors: errmapa}
	}
	return &AuthenticationServiceResponse{Success: true, UserId: response.UserId}
}
func (as *AuthService) Authenticate(user *model.Person, ctx context.Context) *AuthenticationServiceResponse {
	errorvalidate := validatePerson(as, user, false)
	if errorvalidate != nil {
		return &AuthenticationServiceResponse{Success: false, Errors: errorvalidate}
	}
	response := as.repo.GetUser(ctx, user.Email, user.Password)
	if !response.Success {
		errmapa := make(map[string]error)
		errmapa["AuthenticateError"] = response.Errors
		return &AuthenticationServiceResponse{Success: false, Errors: errmapa}
	}
	return &AuthenticationServiceResponse{Success: true, UserId: response.UserId}
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
					erors[err.Field()] = erro.ErrorNotEmail
				case "min":
					errv := fmt.Errorf("%s is too short", err.Field())
					erors[err.Field()] = errv

				default:
					errv := fmt.Errorf("%s is Null", err.Field())
					erors[err.Field()] = errv
				}
			}
			return erors
		}
	}
	return nil
}

var sessions = make(map[string]Session)

type Session struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func (as *AuthService) GenerateSession(userID uuid.UUID) (string, time.Time) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24)
	session := Session{
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	sessions[sessionID] = session
	return sessionID, expiresAt
}
func (as *AuthService) Authorizate(sessionID string) *AuthenticationServiceResponse {
	session, ok := sessions[sessionID]
	if !ok {
		return &AuthenticationServiceResponse{Success: false, UserId: uuid.Nil}
	}
	if time.Now().After(session.ExpiresAt) {

		return &AuthenticationServiceResponse{Success: false, UserId: uuid.Nil}
	}

	return &AuthenticationServiceResponse{Success: true, UserId: session.UserID}
}

/*func (as *AuthService) DeleteSession(token string) error {
	//return nil
}*/
