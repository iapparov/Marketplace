package web

import (
	"net/http"
)

type UserHandler struct{

}

func NewUserHandler() *UserHandler{
	return &UserHandler{
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
}

