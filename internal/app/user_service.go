package app

import (
	"errors"
	"fmt"
	"marketplace/internal/config"
	"strings"
	"unicode/utf8"
	"regexp"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)



func NewUserService(repo UserRepository) *UserService{
	return &UserService{
		repo:repo,
	}
}

func (s *UserService) RegisterUser(req SignUpRequest, config *config.Config) (User, error) {

	_, err := s.repo.FindByLogin(req.Login)
	if err == nil {
		return User{}, errors.New("user already exists")
	}
	
	if req.Login == "" || req.Password == "" {
		return User{}, errors.New("login and password cannot be empty")
	}

	_ , err = isValidLogin(req.Login, config)
	if err != nil {
		return User{}, err
	}
	_ , err = isValidPassword(req.Password, config)
	if err != nil {
		return User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	if config.Username.CaseSensitive {
		req.Login = strings.ToLower(req.Login)
	}

	user := User{
		UUID:     uuid.New(),
		Login:    req.Login,
		Password: string(hashedPassword),
	}
	err = s.repo.SaveNewUser(user)
	if err != nil {
		return User{}, fmt.Errorf("failed to save user: %w", err)
	}	

	return user, nil
}

func isValidLogin(login string, config *config.Config) (bool, error) {
	if utf8.RuneCountInString(login) < config.Username.MinLength || utf8.RuneCountInString(login) > config.Username.MaxLength {
		return false, fmt.Errorf("invalid login length . Must be between %d and %d characters", config.Username.MinLength, config.Username.MaxLength)
	}

	escapedChars := regexp.QuoteMeta(config.Username.AllowedCharacters)
	loginRegexp := regexp.MustCompile(`^[` + escapedChars + `]+$`)
	if !loginRegexp.MatchString(login) {
		return false, errors.New("invalid login characters. Must contain only letters, digits, underscores, or hyphens and must not contain spaces")
	}
	return true, nil
}

func isValidPassword(password string, config *config.Config) (bool, error) {
	if utf8.RuneCountInString(password) < config.Password.MinLength || utf8.RuneCountInString(password) > config.Password.MaxLength {
		return false, fmt.Errorf("invalid password length. Must be between %d and %d characters", config.Password.MinLength, config.Password.MaxLength)
	}
	if !utf8.ValidString(password) {
		return false, errors.New("invalid password characters. Must contain only valid UTF-8 characters")
	}
	hasUpper := !config.Password.RequireUpper
	hasLower := !config.Password.RequireLower
	hasDigit := !config.Password.RequireDigit
	hasSpace := false
	for _, r := range password {
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

func (s *UserService) LoginJwt(req JwtRequest, jwt *JwtProvider, config *config.Config) (JwtResponse, error) {
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

func (s *UserService) RefreshAccessToken(req RefreshJwtRequest, jwt *JwtProvider, config *config.Config) (JwtResponse, error){
	claims, err := jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return JwtResponse{}, errors.New("invalid refresh token")
	}
	user, err := s.repo.FindByUUID(claims["uuid"].(string))
	if err != nil {
		return JwtResponse{}, errors.New("user not found")
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