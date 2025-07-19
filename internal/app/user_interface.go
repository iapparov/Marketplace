package app

import (
	"marketplace/internal/config"
)

type UserRepository interface {
	SaveNewUser(user User) error
	FindByLogin(login string) (User, error) // strings.ToLower(req.Login) допилить
	FindByUUID(uuid string) (User, error)
}

type UserServicer interface{
	RegisterUser(req SignUpRequest, config *config.Config) (User, error)
	LoginJwt(req JwtRequest, jwt *JwtProvider, config *config.Config) (JwtResponse, error)
	RefreshAccessToken(req RefreshJwtRequest, jwt *JwtProvider, config *config.Config) (JwtResponse, error)
}