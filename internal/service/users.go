package service

import (
	"context"
	"database/sql"
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
	ErrEmailTaken         = errors.New("email already exist")
	ErrUsernameTaken      = errors.New("username already taken")
)

func NewUserService(s *store.Store, cfg config.AuthConfig) *UserService {
	return &UserService{
		store:      s,
		jwtSecret:  cfg.JWTSecret,
		accessTTL:  cfg.AccessTokenTTL,
		refreshTTL: cfg.RefreshTokenTTL,
	}
}

type LoginReturn struct {
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
		return database.User{}, fmt.Errorf("Can't hash password:%w", err)
	}

	existing, err := s.store.GetUserByEmail(ctx, email)
	if err != nil && existing.ID != uuid.Nil {
		return database.User{}, fmt.Errorf("failed to check email: %w", err)
	}
	if existing.ID != uuid.Nil {
		return database.User{}, ErrEmailTaken
	}

	existing, err = s.store.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return database.User{}, fmt.Errorf("failed to check username: %w", err)
	}
	if existing.ID != uuid.Nil {

		return database.User{}, ErrUsernameTaken
	}

	u, err := s.store.CreateUser(ctx, database.CreateUserParams{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return database.User{}, fmt.Errorf("Can't create user:%w", err)
	}

	return u, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.store.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("Can't delete user:%w", err)
	}
	return nil
}

func (s *UserService) generateTokens(ctx context.Context, user database.User) (LoginReturn, error) {
	accessToken, err := auth.MakeJWT(user.ID, s.jwtSecret, time.Hour)
	if err != nil {
		return LoginReturn{}, fmt.Errorf("cant create token: %w", err)
	}

	refreshKey, err := auth.MakeRefreshToken()
	if err != nil {
		return LoginReturn{}, fmt.Errorf("cant create refresh key: %w", err)
	}

	expiresAt := time.Now().UTC().Add(24 * time.Hour * 60)
	refreshToken, err := s.store.CreateRefreshTokenForUser(ctx, user.ID, refreshKey, expiresAt)
	if err != nil {
		return LoginReturn{}, fmt.Errorf("cant save refresh token: %w", err)
	}

	return LoginReturn{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (LoginReturn, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return LoginReturn{}, ErrNotFound
	}

	match, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		return LoginReturn{}, fmt.Errorf("cant check password hash: %w", err)
	}
	if !match {
		return LoginReturn{}, ErrInvalidCredentials
	}
	return s.generateTokens(ctx, user)
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (LoginReturn, error) {
	user, err := s.CreateUser(ctx, username, email, password)
	if err != nil {
		return LoginReturn{}, err
	}

	return s.generateTokens(ctx, user)
}
