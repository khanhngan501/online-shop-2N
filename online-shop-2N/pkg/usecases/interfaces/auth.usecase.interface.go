package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/models"
	"online-shop-2N/pkg/services/tokens"
)

type AuthUseCase interface {
	//user
	UserSignUp(ctx context.Context, signUpDetails models.User) (err error)
	SingUpOtpVerify(ctx context.Context, otpVerifyDetails requests.OTPVerify) (userID uint, err error)
	GoogleLogin(ctx context.Context, user models.User) (userID uint, err error)
	UserLogin(ctx context.Context, loginDetails requests.Login) (userID uint, err error)
	UserLoginOtpSend(ctx context.Context, loginDetails requests.OTPLogin) (otpID string, err error)
	LoginOtpVerify(ctx context.Context, otpVerifyDetails requests.OTPVerify) (userID uint, err error)

	// admin
	AdminLogin(ctx context.Context, loginDetails requests.Login) (adminID uint, err error)
	// token
	GenerateAccessToken(ctx context.Context, tokenParams GenerateTokenParams) (tokenString string, err error)
	GenerateRefreshToken(ctx context.Context, tokenParams GenerateTokenParams) (tokenString string, err error)
	VerifyAndGetRefreshTokenSession(ctx context.Context, refreshToken string, usedFor tokens.UserType) (models.RefreshSession, error)
}

type GenerateTokenParams struct {
	UserID   uint
	UserType tokens.UserType
}
