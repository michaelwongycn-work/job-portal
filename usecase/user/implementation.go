package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/michaelwongycn/job-portal/domain/request"
	"github.com/michaelwongycn/job-portal/lib/auth"
	"github.com/michaelwongycn/job-portal/lib/cache"
	"github.com/michaelwongycn/job-portal/lib/encrypt"
	"github.com/michaelwongycn/job-portal/repository/appDB"
)

type userImpl struct {
	appDB                appDB.AppDBInterface
	refreshTokenDuration time.Duration
}

func NewUserImpl(appDB appDB.AppDBInterface, refreshTokenDuration time.Duration) UserUsecase {
	return &userImpl{
		appDB:                appDB,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (u *userImpl) Login(ctx context.Context, req request.UserLoginRequest) (*string, *string, error) {
	currTime := time.Now()

	encryptedPassword, err := encrypt.Hash(req.Password)
	if err != nil {
		return nil, nil, err
	}

	user, err := u.appDB.GetUserByEmailAndPassword(ctx, req.Email, encryptedPassword)
	if err != nil {
		return nil, nil, err
	}

	oldUserToken, err := u.appDB.GetUserToken(ctx, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}

	accessToken, refreshToken, err := auth.CreateToken(currTime, user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}

	err = u.appDB.InsertUserToken(ctx, user.ID, accessToken, refreshToken, currTime.Add(time.Minute*u.refreshTokenDuration).Unix())
	if err != nil {
		return nil, nil, err
	}

	if oldUserToken != nil {
		cache.DeleteCache(oldUserToken.AccessToken)
	}

	cache.SetCache(accessToken, refreshToken)
	return &accessToken, &refreshToken, nil
}

func (u *userImpl) Register(ctx context.Context, req request.UserRegisterRequest) (*string, *string, error) {
	currTime := time.Now()

	encryptedPassword, err := encrypt.Hash(req.Password)
	if err != nil {
		return nil, nil, err
	}

	userId, err := u.appDB.InsertUser(ctx, req.Email, encryptedPassword, req.Role)
	if err != nil {
		return nil, nil, err
	}

	accessToken, refreshToken, err := auth.CreateToken(currTime, userId, req.Role)
	if err != nil {
		return nil, nil, err
	}

	err = u.appDB.InsertUserToken(ctx, userId, accessToken, refreshToken, currTime.Add(time.Minute*u.refreshTokenDuration).Unix())
	if err != nil {
		return nil, nil, err
	}

	cache.SetCache(accessToken, refreshToken)
	return &accessToken, &refreshToken, nil
}

func (u *userImpl) Logout(ctx context.Context, req request.UserLogoutRequest) error {
	cache.DeleteCache(req.AccessToken)
	return u.appDB.DeleteUserToken(ctx, req.UserId)
}

func (u *userImpl) RefreshToken(ctx context.Context, req request.UserRefreshTokenRequest) (*string, *string, error) {
	userToken, err := u.appDB.GetUserToken(ctx, req.UserId)
	if err != nil {
		return nil, nil, err
	}

	currTime := time.Now()

	if req.RefreshToken != userToken.RefreshToken || userToken.ExpirationTime < currTime.Unix() {
		return nil, nil, errors.New("invalid refresh token")
	}

	newAccessToken, newRefreshToken, err := auth.CreateToken(currTime, req.UserId, req.Role)
	if err != nil {
		return nil, nil, err
	}

	err = u.appDB.InsertUserToken(ctx, req.UserId, newAccessToken, newRefreshToken, currTime.Add(time.Minute*u.refreshTokenDuration).Unix())
	if err != nil {
		return nil, nil, err
	}

	cache.DeleteCache(userToken.AccessToken)
	cache.SetCache(newAccessToken, newRefreshToken)
	return &newAccessToken, &newRefreshToken, nil
}
