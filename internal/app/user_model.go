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

type UserService struct {
	repo UserRepository
}



// type UserService interface {
// 	RegisterUser(req SignUpRequest) (User, error)
// 	LoginJwt(req JwtRequest, jwt JwtProvider) (JwtResponse, error)
// 	RefreshAccessToken(req RefreshJwtRequest) (JwtResponse, error)
// 	RefreshRefreshToken(req RefreshJwtRequest, oldAccessToken string) (JwtResponse, error)
// }