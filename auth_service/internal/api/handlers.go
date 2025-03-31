package api

import (
	"auth_service/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	jsonResponseType = "application/json"
)

type Handler struct {
	services *service.Service
}
type HTTPResponse struct {
	Success bool              `json:"success"`
	Errors  map[string]string `json:"errors"`
	UserID  uuid.UUID         `json:"data"`
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/reg", h.NonAuthorizedMiddleware(h.Registration)).Methods("POST")
	m.HandleFunc("/auth", h.NonAuthorizedMiddleware(h.Authentication)).Methods("POST")
	m.HandleFunc("/check-session", h.Authorization).Methods("GET")
	return m
}
