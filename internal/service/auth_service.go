package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/auth"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
)

type AuthService struct {
	repository repository.Repository
	config     *config.Config
	logger     *slog.Logger
}

func NewAuthService(r repository.Repository, c *config.Config, l *slog.Logger) *AuthService {
	return &AuthService{repository: r, config: c, logger: l}
}

// possible password encryption
func (s *AuthService) encryptPassword(password string) (string, error) {
	return password, nil
}

// possible password validation
func (s *AuthService) passwordIsValid(password string) (bool, error) {
	return password != "", nil
}

// possible custom login validation
func (s *AuthService) loginIsValid(login string) (bool, error) {
	return login != "", nil
}

func newUser(login string, password string) (*models.User, error) {
	if login == "" {
		return nil, errors.New("empty login")
	}
	if password == "" {
		return nil, errors.New("empty password")
	}
	return &models.User{ID: "", Login: login, Password: password}, nil
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

	//ok, adding user
	encryptedPassword, err := s.encryptPassword(password)
	if err != nil {
		return "", err
	}

	u, err := newUser(login, encryptedPassword)
	if err != nil {
		return "", err
	}

	err = s.repository.BeginTransaction(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		if p := recover(); p != nil {
			s.repository.RollbackTransaction(ctx)
			panic(p)
		} else if err != nil {
			s.repository.RollbackTransaction(ctx)
		} else {
			err = s.repository.CommitTransaction(ctx)
		}
	}()

	user, err := s.repository.AddUser(ctx, u)
	if err != nil {
		return "", err
	}

	// err = s.repository.CommitTransaction(ctx)
	// if err != nil {
	// 	return err
	// }

	t, err := auth.GenerateToken(user.ID, s.config.SecretKey, &s.config.TokenValidityDuration)
	if err != nil {
		return "", err
	}

	return t, nil

}

func (s *AuthService) validatePassword(password string, existingUser *models.User) (bool, error) {

	encryptedPassword, err := s.encryptPassword(password)
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
