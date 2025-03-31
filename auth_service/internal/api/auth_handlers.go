package api

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"auth_service/internal/service"
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
		jsonResponse := badResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	datafromperson, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll Error: %v", err)
		maparesponse["ReadAll"] = erro.ErrorReadAll.Error()
		jsonResponse := badResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	var newperk model.Person
	err = json.Unmarshal(datafromperson, &newperk)
	if err != nil {
		log.Printf("Unmarshal Error: %v", err)
		maparesponse["Unmarshal"] = erro.ErrorUnmarshal.Error()
		jsonResponse := badResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	regresponse := h.services.RegistrateAndLogin(&newperk, r.Context())
	if !regresponse.Success {
		stringMap := convertErrorToString(regresponse)
		log.Printf("Error during user registration: %v", regresponse.Errors)
		jsonResponse := badResponse(w, stringMap)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
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
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Person with id: %v has successfully registered", regresponse.UserId)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) Authentication(w http.ResponseWriter, r *http.Request) {
	maparesponse := make(map[string]string)
	if r.Method != http.MethodPost {
		log.Printf("Invalid request method(expected Post but it was sent %v)", r.Method)
		maparesponse["Method"] = erro.ErrorNotPost.Error()
		jsonResponse := badResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	datafromperson, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll Error: %v", err)
		maparesponse["ReadAll"] = erro.ErrorReadAll.Error()
		jsonResponse := badResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	var newperk model.Person
	err = json.Unmarshal(datafromperson, &newperk)
	if err != nil {
		log.Printf("Unmarshal Error: %v", err)
		maparesponse["Unmarshal"] = erro.ErrorUnmarshal.Error()
		jsonResponse := badResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	auresponse := h.services.AuthenticateAndLogin(&newperk, r.Context())
	if !auresponse.Success {
		stringMap := convertErrorToString(auresponse)
		log.Printf("Error during user authentication: %v", auresponse.Errors)
		jsonResponse := badResponse(w, stringMap)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}

	addCookie(w, auresponse.SessionId, auresponse.ExpirationTime)
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  auresponse.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		log.Printf("Marshal Error: %v", err)
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Person with id: %v has successfully authenticated", auresponse.UserId)
	fmt.Fprint(w, string(jsonResponse))

}
func (h *Handler) Authorization(w http.ResponseWriter, r *http.Request) {
	maparesponse := make(map[string]string)
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		log.Println("The person's session was not found")
		maparesponse["SessionId"] = erro.ErrorInvalidSessionID.Error()
		jsonResponse := badResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	response := h.services.Authorization(r.Context(), sessionID)
	if !response.Success {
		stringMap := convertErrorToString(response)

		jsonResponse := badResponse(w, stringMap)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  response.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Person with id: %v has successfully authorizated", response.UserId)
	fmt.Fprint(w, string(jsonResponse))
}
func badResponse(w http.ResponseWriter, vc map[string]string) []byte {
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusMethodNotAllowed)
	response := HTTPResponse{
		Success: false,
		Errors:  vc,
		UserID:  uuid.Nil,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal Error: %v", err)
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return nil
	}
	return jsonResponse
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
