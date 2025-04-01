package api

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"auth_service/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {

	maparesponse := make(map[string]string)
	if r.Method != http.MethodPost {
		log.Printf("Invalid request method(expected Post but it was sent %v)", r.Method)
		maparesponse["Method"] = erro.ErrorNotPost.Error()
		badResponse(w, maparesponse, http.StatusMethodNotAllowed)

		return
	}
	datafromperson, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll Error: %v", err)
		maparesponse["ReadAll"] = erro.ErrorReadAll.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)

		return
	}
	var newperk model.Person
	err = json.Unmarshal(datafromperson, &newperk)
	if err != nil {
		log.Printf("Unmarshal Error: %v", err)
		maparesponse["Unmarshal"] = erro.ErrorUnmarshal.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)

		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	regresponse := h.services.RegistrateAndLogin(ctx, &newperk)
	if !regresponse.Success {
		stringMap := convertErrorToString(regresponse)
		log.Printf("Error during user registration: %v", regresponse.Errors)
		badResponse(w, stringMap, http.StatusBadRequest)

		return
	}

	addCookie(w, regresponse.SessionId, regresponse.ExpirationTime)
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  regresponse.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		log.Printf("Marshal Error: %v", err)
		maparesponse["Marshal"] = erro.ErrorMarshal.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)

		return
	}
	log.Printf("Person with id: %v has successfully registered", regresponse.UserId)
	addCookie(w, regresponse.SessionId, regresponse.ExpirationTime)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) Authentication(w http.ResponseWriter, r *http.Request) {

	maparesponse := make(map[string]string)
	if r.Method != http.MethodPost {
		log.Printf("Invalid request method(expected Post but it was sent %v)", r.Method)
		maparesponse["Method"] = erro.ErrorNotPost.Error()
		badResponse(w, maparesponse, http.StatusMethodNotAllowed)

		return
	}
	datafromperson, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll Error: %v", err)
		maparesponse["ReadAll"] = erro.ErrorReadAll.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)

		return
	}
	var newperk model.Person
	err = json.Unmarshal(datafromperson, &newperk)
	if err != nil {
		log.Printf("Unmarshal Error: %v", err)
		maparesponse["Unmarshal"] = erro.ErrorUnmarshal.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)

		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	auresponse := h.services.AuthenticateAndLogin(ctx, &newperk)
	if !auresponse.Success {
		stringMap := convertErrorToString(auresponse)
		log.Printf("Error during user authentication: %v", auresponse.Errors)
		badResponse(w, stringMap, http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  auresponse.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		log.Printf("Marshal Error: %v", err)
		maparesponse["Marshal"] = erro.ErrorMarshal.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)

		return
	}
	log.Printf("Person with id: %v has successfully authenticated", auresponse.UserId)
	addCookie(w, auresponse.SessionId, auresponse.ExpirationTime)
	fmt.Fprint(w, string(jsonResponse))

}
func (h *Handler) Authorization(w http.ResponseWriter, r *http.Request) {
	maparesponse := make(map[string]string)
	if r.Method != http.MethodGet {
		log.Printf("Invalid request method(expected Get but it was sent %v)", r.Method)
		maparesponse["Method"] = erro.ErrorNotGet.Error()
		badResponse(w, maparesponse, http.StatusMethodNotAllowed)

		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("The person's session was not found")
		if err == http.ErrNoCookie {
			maparesponse["SessionId"] = erro.ErrorInvalidSessionID.Error()
		} else {
			log.Printf("Error reading cookie: %v", err)
			maparesponse["SessionId"] = erro.ErrorInternalServer.Error()
		}
		badResponse(w, maparesponse, http.StatusUnauthorized)

		return
	}
	sessionID := cookie.Value
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	response := h.services.Authorization(ctx, sessionID)
	if !response.Success {
		stringMap := convertErrorToString(response)

		badResponse(w, stringMap, http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  response.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		log.Printf("Marshal Error: %v", err)
		maparesponse["Marshal"] = erro.ErrorMarshal.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)

		return
	}
	log.Printf("Person with id: %v has successfully authorizated", response.UserId)
	fmt.Fprint(w, string(jsonResponse))
}
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	maparesponse := make(map[string]string)
	userID, ok := getUserIDFromRequestContext(r)
	if !ok {
		log.Println("Error getting the UserId from the request context")
		maparesponse["UserId"] = erro.ErrorGetUserId.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("Invalid request method(expected Post but it was sent %v)", r.Method)
		maparesponse["Method"] = erro.ErrorNotPost.Error()
		badResponse(w, maparesponse, http.StatusMethodNotAllowed)

		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("Unexpected error getting session cookie (should have been validated by middleware): %v", err)
	}
	sessionID := cookie.Value
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	response := h.services.Logout(ctx, sessionID)
	if !response.Success {
		stringMap := convertErrorToString(response)

		badResponse(w, stringMap, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		log.Printf("Marshal Error: %v", err)
		maparesponse["Marshal"] = erro.ErrorMarshal.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)

		return
	}
	log.Printf("Person with id: %v has successfully logged out", userID)
	deleteCookie(w)
	fmt.Fprint(w, string(jsonResponse))

}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	maparesponse := make(map[string]string)
	userID, ok := getUserIDFromRequestContext(r)
	if !ok {
		log.Println("Error getting the UserId from the request context")
		maparesponse["UserId"] = erro.ErrorGetUserId.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodDelete {
		log.Printf("Invalid request method(expected Delete but it was sent %v)", r.Method)
		maparesponse["Method"] = erro.ErrorNotDelete.Error()
		badResponse(w, maparesponse, http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("Unexpected error getting session cookie (should have been validated by middleware): %v", err)
	}
	sessionID := cookie.Value
	password, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll Error: %v", err)
		maparesponse["ReadAll"] = erro.ErrorReadAll.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)

		return
	}
	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	response := h.services.DeleteAccount(ctx, sessionID, userID, string(password))
	if !response.Success {
		stringMap := convertErrorToString(response)
		badResponse(w, stringMap, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		log.Printf("Marshal Error: %v", err)
		maparesponse["Marshal"] = erro.ErrorMarshal.Error()
		badResponse(w, maparesponse, http.StatusInternalServerError)
		return
	}
	log.Printf("Person with id: %v has successfully delete account with all data", userID)
	deleteCookie(w)
	fmt.Fprint(w, string(jsonResponse))

}
func badResponse(w http.ResponseWriter, vc map[string]string, statusCode int) {
	response := HTTPResponse{
		Success: false,
		Errors:  vc,
		UserID:  uuid.Nil,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal Error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(statusCode)

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("Write Error: %v", err)
		return
	}
}
func addCookie(w http.ResponseWriter, sessionID string, duration time.Time) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		//Secure:   true, // Рекомендуется использовать HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  duration,
	}

	http.SetCookie(w, cookie)
}
func deleteCookie(w http.ResponseWriter) {

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		//Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
}
func convertErrorToString(mapa *service.ServiceResponse) map[string]string {
	stringMap := make(map[string]string)
	for key, err := range mapa.Errors {
		if err != nil {
			stringMap[key] = err.Error()
		} else {
			stringMap[key] = ""
		}
	}
	return stringMap
}
func getUserIDFromRequestContext(r *http.Request) (uuid.UUID, bool) {
	return getUserIDFromContext(r.Context())
}
func getUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {

		return uuid.Nil, false
	}
	return userID, true
}
