package web

import (
	"encoding/json"
	"marketplace/internal/app"
	"marketplace/internal/config"
	"net/http"
	"go.uber.org/zap"
)

type UserHandler struct{
	config *config.Config
	logger *zap.Logger
	jwt *app.JwtProvider
	app app.UserServicer
}

func NewUserHandler(app app.UserServicer, config *config.Config, jwt *app.JwtProvider, logger *zap.Logger) *UserHandler{
	return &UserHandler{
		app: app, 
		config: config,
		jwt: jwt,
		logger: logger,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login_req app.JwtRequest
	if err := json.NewDecoder(r.Body).Decode(&login_req); err != nil{
		h.logger.Warn("invalid login request body", zap.Error(err))
		http.Error(w, "bad login request", http.StatusBadRequest)
		return
	}
	login_resp, err := h.app.LoginJwt(login_req, h.jwt, h.config)
	if err != nil {
		h.logger.Warn("login failed", zap.Error(err), zap.String("login", login_req.Login))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.logger.Info("login successful", zap.String("login", login_req.Login))
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(login_resp)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user_req app.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&user_req); err != nil{
		h.logger.Warn("invalid registration body", zap.Error(err))
		http.Error(w, "bad registration request", http.StatusBadRequest)
		return
	}

	user, err := h.app.RegisterUser(user_req, h.config)

	if err !=nil {
		h.logger.Warn("registration failed", zap.Error(err), zap.String("login", user_req.Login))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user_resp app.SignUpResponse
	user_resp.Login = user.Login
	user_resp.UUID = user.UUID

	h.logger.Info("registration successful", zap.String("login", user_resp.Login))
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user_resp)
}

func (h *UserHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	var refresh_req app.RefreshJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&refresh_req); err != nil {
		h.logger.Warn("bad refresh token request", zap.Error(err))
		http.Error(w, "bad refresh request", http.StatusBadRequest)
		return
	}

	refresh_resp, err := h.app.RefreshAccessToken(refresh_req, h.jwt, h.config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.logger.Info("refresh token successful")
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(refresh_resp)
}
