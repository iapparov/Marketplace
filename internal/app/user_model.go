package app

import (
	"github.com/google/uuid"
)

type User struct {
	UUID     uuid.UUID	`json:"uuid"`
	Login    string		`json:"login"`
	Password string		`json:"password"`
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UUID     uuid.UUID	`json:"uuid"`
	Login    string		`json:"login"`
}


type UserService struct {
	repo UserRepository
}