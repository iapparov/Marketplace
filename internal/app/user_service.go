package app

import (
	"errors"
	"marketplace/internal/config"
	"unicode/utf8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepository;
}

func NewUserService(repo UserRepository) *UserService{
	return &UserService{
		repo:repo,
	}
}

func (s *UserService) RegisterUser(req SignUpRequest) (User, error) {

	_, err := s.repo.FindByLogin(req.Login)
	if err == nil {
		return User{}, errors.New("user already exists")
	}
	
	if req.Login == "" || req.Password == "" {
		return User{}, errors.New("login and password cannot be empty")
	}

	_ , err = isValidLogin(req.Login)
	if err != nil {
		return User{}, err
	}
	_ , err = isValidPassword(req.Password)
	if err != nil {
		return User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user := User{
		UUID:     uuid.New(),
		Login:    req.Login,
		Password: string(hashedPassword),
	}

	return user, nil
}

func isValidLogin(login string) (bool, error) {
	if utf8.RuneCountInString(login) < 3 || utf8.RuneCountInString(login) > 20 {
		return false, errors.New("invalid login length . Must be between 3 and 20 characters")
	}
	for _, r := range login {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r != ' ') {
			return false, errors.New("invalid login characters. Must contain only letters, digits, underscores, or hyphens and must not contain spaces")
		}
	}
	return true, nil
}

func isValidPassword(password string) (bool, error) {
	if utf8.RuneCountInString(password) < 6 || utf8.RuneCountInString(password) > 20 {
		return false, errors.New("invalid password length. Must be between 6 and 20 characters")
	}
	if !utf8.ValidString(password) {
		return false, errors.New("invalid password characters. Must contain only valid UTF-8 characters")
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpace := false
	for _, r := range password {
		if utf8.RuneCountInString(password) < 6 || utf8.RuneCountInString(password) > 20 {
			return false, errors.New("invalid password length. Must be between 6 and 20 characters")
		}
		if !utf8.ValidString(password) {
			return false, errors.New("invalid password characters. Must contain only valid UTF-8 characters")
		}
		if r >= 'A' && r <= 'Z' {
			hasUpper = true
		} else if r >= 'a' && r <= 'z' {
			hasLower = true
		} else if r >= '0' && r <= '9' {
			hasDigit = true
		} else if r == ' ' {
			hasSpace = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || hasSpace {
		return false, errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and must not contain spaces")
	}
	return hasUpper && hasLower && hasDigit, nil
}

func (s *UserService) LoginJwt(req JwtRequest, jwt JwtProvider, config *config.Config) (JwtResponse, error) {
	user, err := s.repo.FindByLogin(req.Login)
	if err != nil {
		return JwtResponse{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return JwtResponse{}, errors.New("unauthorized")
	}
	accessToken, err := jwt.GenerateAccessToken(user, config)
	if err != nil {
		return JwtResponse{}, errors.New("failed to generate access token")
	}
	refreshToken, err := jwt.GenerateRefreshToken(user, config)
	if err != nil {
		return JwtResponse{}, errors.New("failed to generate access token")
	}

	return JwtResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
		RefreshToken: refreshToken,

	}, nil
}

func (s *UserService) RefreshAccessToken(req RefreshJwtRequest, config *config.Config) (JwtResponse, error){
	var jwt JwtProvider
	id, err := jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return JwtResponse{}, errors.New("invalid refresh token")
	}
	flag, userinfo := s.repo.FindByUUID(id.String())
	if !flag {
		return JwtResponse{}, errors.New("user not found")
	}
	var user User
	user.UUID = id
	user.Login = userinfo[0]
	user.Password = userinfo[1]
	accessToken, err := jwt.GenerateAccessToken(user, config)
	if err != nil {
		return JwtResponse{}, errors.New("failed to generate access token")
	}

	return JwtResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
		RefreshToken: req.RefreshToken,

	}, nil

}
func (s *UserService) RefreshRefreshToken(req RefreshJwtRequest, oldAccessToken string, config *config.Config) (JwtResponse, error){
	var jwt JwtProvider
	id, err := jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return JwtResponse{}, errors.New("invalid refresh token")
	}
	flag, userinfo := s.repo.FindByUUID(id.String())
	if !flag {
		return JwtResponse{}, errors.New("user not found")
	}
	var user User
	user.UUID = id
	user.Login = userinfo[0]
	user.Password = userinfo[1]
	refreshToken, err := jwt.GenerateRefreshToken(user, config)
	if err != nil {
		return JwtResponse{}, errors.New("failed to generate access token")
	}

	return JwtResponse{
		Type:        "Bearer",
		AccessToken: oldAccessToken,
		RefreshToken: refreshToken,

	}, nil
}