package api

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"auth_service/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	maparesponse := make(map[string]string)
	if r.Method != http.MethodPost {

		maparesponse["Method"] = erro.ErrorNotPost.Error()
		jsonResponse := BadResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	datafromperson, err := io.ReadAll(r.Body)
	if err != nil {

		maparesponse["ReadAll"] = erro.ErrorReadAll.Error()
		jsonResponse := BadResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	var newperk model.Person
	err = json.Unmarshal(datafromperson, &newperk)
	if err != nil {
		maparesponse["Unmarshal"] = erro.ErrorUnmarshal.Error()
		jsonResponse := BadResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	regresponse := h.services.Registrate(&newperk, r.Context())
	if !regresponse.Success {
		stringMap := convertErrorToString(regresponse)

		jsonResponse := BadResponse(w, stringMap)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	sessionresponse := h.services.GenerateSession(context.TODO(), regresponse.UserId)
	AddCookie(w, sessionID, time)
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  regresponse.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(jsonResponse))
}
func (h *Handler) Authentication(w http.ResponseWriter, r *http.Request) {
	maparesponse := make(map[string]string)
	if r.Method != http.MethodPost {
		maparesponse["Method"] = erro.ErrorNotPost.Error()
		jsonResponse := BadResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	datafromperson, err := io.ReadAll(r.Body)
	if err != nil {
		maparesponse["ReadAll"] = erro.ErrorReadAll.Error()
		jsonResponse := BadResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	var newperk model.Person
	err = json.Unmarshal(datafromperson, &newperk)
	if err != nil {
		maparesponse["Unmarshal"] = erro.ErrorUnmarshal.Error()
		jsonResponse := BadResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	response := h.services.Authenticate(&newperk, r.Context())
	if !response.Success {
		stringMap := convertErrorToString(response)

		jsonResponse := BadResponse(w, stringMap)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	_ = h.services.GenerateSession(r.Context(), newperk.Id)
	//AddCookie(w, responseredis., time)
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  response.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(jsonResponse))

}
func (h *Handler) Authorization(w http.ResponseWriter, r *http.Request) {
	maparesponse := make(map[string]string)
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		maparesponse["SessionId"] = erro.ErrorInvalidSessionID.Error()
		jsonResponse := BadResponse(w, maparesponse)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	response := h.services.Authorizate(r.Context(), sessionID)
	if !response.Success {
		stringMap := convertErrorToString(response)

		jsonResponse := BadResponse(w, stringMap)
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
	fmt.Fprint(w, string(jsonResponse))
}
func BadResponse(w http.ResponseWriter, vc map[string]string) []byte {
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusMethodNotAllowed)
	response := HTTPResponse{
		Success: false,
		Errors:  vc,
		UserID:  uuid.Nil,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return nil
	}
	return jsonResponse
}
func AddCookie(w http.ResponseWriter, sessionID string, duration time.Time) {
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
