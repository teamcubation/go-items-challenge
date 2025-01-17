package application

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"github.com/teamcubation/go-items-challenge/internal/ports/out"
	"github.com/teamcubation/go-items-challenge/internal/utils"
)

var (
	ErrUsernameExists   = fmt.Errorf("username already exists")
	ErrUsernameNotFound = fmt.Errorf("username not found")
)

type authService struct {
	repo out.UserRepository
}

func NewAuthService(repo out.UserRepository) *authService {
	return &authService{repo: repo}
}

func (srv *authService) RegisterUser(ctx context.Context, newUser *user.User) (*user.User, error) {
	lowerUsername := strings.ToLower(newUser.Username)
	userFound, err := srv.repo.GetUserByUsername(ctx, lowerUsername)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	if userFound != nil {
		return nil, ErrUsernameExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	newUser.Username = strings.ToUpper(string(newUser.Username[0])) + strings.ToLower(newUser.Username[1:])
	newUser.Password = string(hashedPassword)

	if err := srv.repo.CreateUser(ctx, newUser); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return newUser, nil
}

func (srv *authService) Login(ctx context.Context, creds user.Credentials) (string, error) {
	userFound, err := srv.repo.GetUserByUsername(ctx, creds.Username)

	if userFound.Username != creds.Username {
		return "", fmt.Errorf("invalid user: %s", creds.Username)
	}

	if err != nil {
		return "", fmt.Errorf("error fetching user: %w", err)
	}

	if userFound == nil {
		return "", ErrUsernameNotFound
	}

	if !utils.CheckPasswordHash(creds.Password, userFound.Password) {
		return "", fmt.Errorf("invalid password: %s", creds.Username)
	}

	// Generate token (assuming you have a function to generate JWT tokens)
	token, err := utils.GenerateToken(userFound.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %s", creds.Username)
	}

	return token, nil
}
