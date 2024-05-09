package user

import (
	"context"

	"github.com/michaelwongycn/job-portal/domain/request"
)

type UserUsecase interface {
	Login(ctx context.Context, req request.UserLoginRequest) (*string, *string, error)
	Register(ctx context.Context, req request.UserRegisterRequest) (*string, *string, error)
	Logout(ctx context.Context, req request.UserLogoutRequest) error
	RefreshToken(ctx context.Context, req request.UserRefreshTokenRequest) (*string, *string, error)
}
