package api

import (
	"log"
	"net/http"
)

func (handler *Handler) NonAuthorizedMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			} else {

				log.Printf("Error reading cookie: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		sessionID := cookie.Value
		response := handler.services.Authorization(r.Context(), sessionID)
		if response.Success {
			stringMap := convertErrorToString(response)
			badResponse(w, stringMap, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}
