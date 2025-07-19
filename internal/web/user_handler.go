package web

import (
	"encoding/json"
	"marketplace/internal/app"
	"marketplace/internal/config"
	"net/http"
)

type UserHandler struct{
	config *config.Config
	jwt *app.JwtProvider
	app app.UserServicer
}

func NewUserHandler(app app.UserServicer, config *config.Config, jwt *app.JwtProvider) *UserHandler{
	return &UserHandler{
		app: app, 
		config: config,
		jwt: jwt,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login_req app.JwtRequest
	if err := json.NewDecoder(r.Body).Decode(&login_req); err != nil{
		http.Error(w, "bad login request", http.StatusBadRequest)
		return
	}
	login_resp, err := h.app.LoginJwt(login_req, h.jwt, h.config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(login_resp)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user_req app.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&user_req); err != nil{
		http.Error(w, "bad registration request", http.StatusBadRequest)
		return
	}

	user, err := h.app.RegisterUser(user_req, h.config)

	if err !=nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var user_resp app.SignUpResponse
	user_resp.Login = user.Login
	user_resp.UUID = user.UUID

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user_resp)
}

func (h *UserHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	var refresh_req app.RefreshJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&refresh_req); err != nil {
		http.Error(w, "bad refresh request", http.StatusBadRequest)
		return
	}

	refresh_resp, err := h.app.RefreshAccessToken(refresh_req, h.jwt, h.config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(refresh_resp)
}
