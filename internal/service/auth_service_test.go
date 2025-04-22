package service

import (
	"context"
	"testing"
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/auth"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/logging"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {

	ctx := context.Background()

	repo, err := repository.NewInMemoryRepository()
	require.NoError(t, err)

	config := &config.Config{SecretKey: "secretkey", TokenValidityDuration: 1 * time.Minute}
	logger := logging.NewLogger()

	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"OK", args{ctx: ctx, login: "login", password: "password"}, false},
		{"Empty Login", args{ctx: ctx, login: "", password: "password"}, true},
		{"Empty Login", args{ctx: ctx, login: "login", password: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthService{
				repository: repo,
				config:     config,
				logger:     logger,
			}
			got, err := s.Register(tt.args.ctx, tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				userID, err := auth.GetUserIDFromToken(got, config.SecretKey)
				require.NoError(t, err)
				require.NotZero(t, userID)
			}

		})
	}
}

func TestAuthService_Login(t *testing.T) {

	ctx := context.Background()

	repo, err := repository.NewInMemoryRepository()
	require.NoError(t, err)

	_, err = repo.AddUser(ctx, &models.User{Login: "login", Password: "password"})
	require.NoError(t, err)

	config := &config.Config{SecretKey: "secretkey", TokenValidityDuration: 1 * time.Minute}
	logger := logging.NewLogger()

	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"OK", args{ctx, "login", "password"}, false},
		{"Wrongpassword", args{ctx, "login", "wrongpassword"}, true},
		{"UnknownLogin", args{ctx, "unknownlogin", "password"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthService{
				repository: repo,
				config:     config,
				logger:     logger,
			}
			got, err := s.Login(tt.args.ctx, tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				userID, err := auth.GetUserIDFromToken(got, config.SecretKey)
				require.NoError(t, err)
				require.NotZero(t, userID)
			}
		})
	}
}
