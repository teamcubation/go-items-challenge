package in

import (
	"context"

	"github.com/teamcubation/go-items-challenge/internal/domain/user"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user *user.User) (*user.User, error)
	Login(ctx context.Context, crd user.Credentials) (string, error)
}
