package service

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/auth"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
)

type AuthService struct {
	BaseService
	repository repository.Repository
	config     *config.Config
	logger     *slog.Logger
}

func NewAuthService(r repository.Repository, c *config.Config, l *slog.Logger) *AuthService {
	return &AuthService{repository: r, config: c, logger: l}
}

// possible password encryption
func (s *AuthService) encryptPassword(salt []byte, password string) (string, error) {

	passwordHash := auth.HashPassword(password, salt)

	return passwordHash, nil
}

// possible password validation
func (s *AuthService) passwordIsValid(password string) (bool, error) {
	return password != "", nil
}

// possible custom login validation
func (s *AuthService) loginIsValid(login string) (bool, error) {
	return login != "", nil
}

func newUser(login string, password string, salt string) (*models.User, error) {
	if login == "" {
		return nil, errors.New("empty login")
	}
	if password == "" {
		return nil, errors.New("empty password")
	}
	return &models.User{ID: "", Login: login, Password: password, Salt: salt}, nil
}

func (s *AuthService) Register(ctx context.Context, login string, password string) (string, error) {

	loginIsValid, err := s.loginIsValid(login)
	if err != nil {
		return "", err
	}
	if !loginIsValid {
		return "", common.ErrorInvalidLoginFormat
	}

	passwordIsValid, err := s.passwordIsValid(password)
	if err != nil {
		return "", err
	}
	if !passwordIsValid {
		return "", common.ErrorInvalidPasswordFormat
	}

	salt, err := auth.GenerateSalt(16)

	if err != nil {
		return "", err
	}

	//ok, adding user
	encryptedPassword, err := s.encryptPassword(salt, password)
	if err != nil {
		return "", err
	}

	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	u, err := newUser(login, encryptedPassword, saltBase64)
	if err != nil {
		return "", err
	}

	tx, err := s.repository.UnitOfWork().Begin(ctx)
	if err != nil {
		return "", err
	}
	defer s.EndTransaction(tx, &err)

	user, err := s.repository.AddUser(ctx, u)
	if err != nil {
		return "", err
	}

	t, err := auth.GenerateToken(user.ID, s.config.SecretKey, &s.config.TokenValidityDuration)
	if err != nil {
		return "", err
	}

	return t, nil

}

func (s *AuthService) validatePassword(password string, existingUser *models.User) (bool, error) {

	salt, err := base64.StdEncoding.DecodeString(existingUser.Salt)
	if err != nil {
		return false, err
	}

	encryptedPassword, err := s.encryptPassword(salt, password)
	if err != nil {
		return false, err
	}

	return encryptedPassword == existingUser.Password, nil

}

func (s *AuthService) Login(ctx context.Context, login string, password string) (string, error) {

	existingLogin, err := s.repository.FindUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	//ok, adding user
	passwordIsOk, err := s.validatePassword(password, &existingLogin)
	if err != nil {
		return "", err
	}

	if !passwordIsOk {
		return "", common.ErrorInvalidLoginPassword
	}

	t, err := auth.GenerateToken(existingLogin.ID, s.config.SecretKey, &s.config.TokenValidityDuration)
	if err != nil {
		return "", err
	}

	return t, nil
}
