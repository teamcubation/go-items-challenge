package application

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"github.com/teamcubation/go-items-challenge/internal/ports/out"
	"github.com/teamcubation/go-items-challenge/internal/utils"
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
		return nil, presenter.ErrFetchingUser
	}

	if userFound != nil {
		return nil, presenter.ErrUsernameExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, presenter.ErrHashingPassword
	}

	newUser.Username = strings.ToUpper(string(newUser.Username[0])) + strings.ToLower(newUser.Username[1:])
	newUser.Password = string(hashedPassword)

	if err := srv.repo.CreateUser(ctx, newUser); err != nil {
		return nil, presenter.ErrCreatingUser
	}

	return newUser, nil
}

func (srv *authService) Login(ctx context.Context, creds user.Credentials) (string, error) {
	userFound, err := srv.repo.GetUserByUsername(ctx, creds.Username)
	if err != nil {
		return "", presenter.ErrFetchingUser
	}

	if userFound == nil || userFound.Username != creds.Username {
		return "", presenter.ErrUsernameNotFound
	}

	if !utils.CheckPasswordHash(creds.Password, userFound.Password) {
		return "", presenter.ErrHashingPassword
	}

	// Generate token (assuming you have a function to generate JWT tokens)
	token, err := utils.GenerateToken(userFound.ID)
	if err != nil {
		return "", presenter.ErrTokenGeneration
	}

	return token, nil
}
