package repositories

import (
	"context"
	"online-shop-2N/pkg/models"
	"online-shop-2N/pkg/repositories/interfaces"

	"gorm.io/gorm"
)

type authDatabase struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) interfaces.AuthRepository {
	return &authDatabase{
		DB: db,
	}
}

// FindOtpSession implements interfaces.AuthRepository.
func (c *authDatabase) FindOtpSession(ctx context.Context, otpID string) (otpSession models.OtpSession, err error) {

	query := `SELECT * FROM otp_sessions WHERE otp_id = $1`

	err = c.DB.Raw(query, otpID).Scan(&otpSession).Error

	return otpSession, err
}

// FindRefreshSessionByTokenID implements interfaces.AuthRepository.
func (c *authDatabase) FindRefreshSessionByTokenID(ctx context.Context, tokenID string) (refreshSession models.RefreshSession, err error) {
	query := `SELECT * FROM refresh_sessions WHERE token_id = $1`

	err = c.DB.Raw(query, tokenID).Scan(&refreshSession).Error

	return
}

// SaveOtpSession implements interfaces.AuthRepository.
func (c *authDatabase) SaveOtpSession(ctx context.Context, otpSession models.OtpSession) error {

	query := `INSERT INTO otp_sessions (otp_id, user_id, phone ,expire_at) 
	VALUES ($1, $2, $3, $4)`
	err := c.DB.Exec(query, otpSession.OtpID, otpSession.UserID, otpSession.Phone, otpSession.ExpireAt).Error
	return err
}

// SaveRefreshSession implements interfaces.AuthRepository.
func (c *authDatabase) SaveRefreshSession(ctx context.Context, refreshSession models.RefreshSession) error {
	query := `INSERT INTO refresh_sessions (token_id, user_id, refresh_token, expire_at) 
VALUES ($1, $2, $3, $4)`
	err := c.DB.Exec(query, refreshSession.TokenID, refreshSession.UserID, refreshSession.RefreshToken, refreshSession.ExpireAt).Error

	return err
}
