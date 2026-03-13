package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ihyaulhaq/url-shotener-BE/internal/auth"
	"github.com/ihyaulhaq/url-shotener-BE/internal/config"
	"github.com/ihyaulhaq/url-shotener-BE/internal/database"
	"github.com/ihyaulhaq/url-shotener-BE/internal/store"
)

type UserService struct {
	store      *store.Store
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

var (
	ErrNotFound           = errors.New("record not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func NewUserService(s *store.Store, cfg config.AuthConfig) *UserService {
	return &UserService{
		store:      s,
		jwtSecret:  cfg.JWTSecret,
		accessTTL:  cfg.AccessTokenTTL,
		refreshTTL: cfg.RefreshTokenTTL,
	}
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (database.User, error) {

	u, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return database.User{}, fmt.Errorf("error:%v", err)
	}

	return u, nil
}

func (s *UserService) CreateUser(ctx context.Context, username, email, password string) (database.User, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return database.User{}, fmt.Errorf("Can't hash password:%v", err)
	}

	existing, err := s.store.GetUserByEmail(ctx, email)
	if err == nil && existing.ID != uuid.Nil {
		return database.User{}, errors.New("email already taken")
	}
	if existing.ID != uuid.Nil {
		return database.User{}, errors.New("email already taken")
	}
	existing, _ = s.store.GetUserByUsername(ctx, username)
	if existing.ID != uuid.Nil {

		return database.User{}, errors.New("username already taken")
	}

	u, err := s.store.CreateUser(ctx, database.CreateUserParams{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return database.User{}, fmt.Errorf("Can't create user:%v", err)
	}

	return u, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.store.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("Can't delete user:%v", err)
	}
	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (LoginResult, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return LoginResult{}, ErrNotFound
	}

	match, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		return LoginResult{}, fmt.Errorf("cant check password hash: %v", err)
	}
	if !match {
		return LoginResult{}, ErrInvalidCredentials
	}

	accessToken, err := auth.MakeJWT(user.ID, s.jwtSecret, s.accessTTL)
	if err != nil {
		return LoginResult{}, fmt.Errorf("cant create token: %v", err)
	}

	refreshKey, err := auth.MakeRefreshToken()
	if err != nil {
		return LoginResult{}, fmt.Errorf("cant create refresh key: %v", err)
	}

	expiresAt := time.Now().UTC().Add(s.refreshTTL)
	refreshToken, err := s.store.CreateRefreshTokenForUser(ctx, user.ID, refreshKey, expiresAt)
	if err != nil {
		return LoginResult{}, fmt.Errorf("cant save refresh token: %v", err)
	}

	return LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}
