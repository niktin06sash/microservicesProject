package service

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"auth_service/internal/repository"
	"context"
	"fmt"

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

func (as *AuthService) Registrate(user *model.Person, ctx context.Context) *ServiceResponse {
	errorvalidate := validatePerson(as, user, true)
	if errorvalidate != nil {
		return &ServiceResponse{Success: false, Errors: errorvalidate}
	}
	hashpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		hashmapa := make(map[string]error)
		hashmapa["HashPass"] = erro.ErrorHashPass
		return &ServiceResponse{Success: false, Errors: hashmapa}
	}
	user.Password = string(hashpass)
	userID := uuid.New()
	user.Id = userID
	response := as.repo.CreateUser(ctx, user)
	if !response.Success {
		errmapa := make(map[string]error)
		errmapa["CreateError"] = response.Errors
		return &ServiceResponse{Success: false, Errors: errmapa}
	}
	return &ServiceResponse{Success: true, UserId: response.UserId}
}
func (as *AuthService) Authenticate(user *model.Person, ctx context.Context) *ServiceResponse {
	errorvalidate := validatePerson(as, user, false)
	if errorvalidate != nil {
		return &ServiceResponse{Success: false, Errors: errorvalidate}
	}
	response := as.repo.GetUser(ctx, user.Email, user.Password)
	if !response.Success {
		errmapa := make(map[string]error)
		errmapa["AuthenticateError"] = response.Errors
		return &ServiceResponse{Success: false, Errors: errmapa}
	}
	return &ServiceResponse{Success: true, UserId: response.UserId}
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
