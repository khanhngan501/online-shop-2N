package handlers

import (
	"errors"
	"net/http"
	"online-shop-2N/pkg/api/handlers/interfaces"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/config"
	"online-shop-2N/pkg/models"
	"online-shop-2N/pkg/services/tokens"
	"online-shop-2N/pkg/usecases"

	usecaseInterface "online-shop-2N/pkg/usecases/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

const (
	authorizationHeaderKey = "Authorization"
	authorizationType      = "Bearer"
)

type AuthHandler struct {
	authUseCase usecaseInterface.AuthUseCase
	config      config.Config
}

func NewAuthHandler(authUsecase usecaseInterface.AuthUseCase, config config.Config) interfaces.AuthHandler {
	return &AuthHandler{
		authUseCase: authUsecase,
		config:      config,
	}
}

// UserLogin godoc
//
//	@Summary		Login with password (User)
//	@Description	API for user to login with email | phone | user_name with password
//	@Id				UserLogin
//	@Tags			User Authentication
//	@Param			inputs	body	requests.Login{}	true	"Login Details"
//	@Router			/auth/sign-in [post]
//	@Success		200	{object}	responses.responses{data=responses.TokenResponse}	"Successfully logged in"
//	@Failure		400	{object}	responses.responses{}								"Invalid inputs"
//	@Failure		403	{object}	responses.responses{}								"User blocked by admin"
//	@Failure		401	{object}	responses.responses{}								"User not exist with given login credentials"
//	@Failure		500	{object}	responses.responses{}								"Failed to login"
func (c *AuthHandler) UserLogin(ctx *gin.Context) {

	var body requests.Login

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	userID, err := c.authUseCase.UserLogin(ctx, body)

	if err != nil {

		var statusCode int

		switch {
		case errors.Is(err, usecases.ErrEmptyLoginCredentials):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrUserNotExist):
			statusCode = http.StatusNotFound
		case errors.Is(err, usecases.ErrUserBlocked):
			statusCode = http.StatusForbidden
		case errors.Is(err, usecases.ErrUserNotVerified):
			statusCode = http.StatusUnauthorized
		case errors.Is(err, usecases.ErrWrongPassword):
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}

		responses.ErrorResponse(ctx, statusCode, "Failed to login", err, nil)
		return
	}

	// common functionality for admin and user
	c.setupTokenAndResponse(ctx, tokens.User, userID)
}

// UserLoginOtpSend godoc
//
//	@Summary		Login with Otp send (User)
//	@Description	API for user to send otp for login enter email | phone | user_name : otp will send to user registered number
//	@Id				UserLoginOtpSend
//	@Tags			User Authentication
//	@Param			inputs	body	requests.OTPLogin{}	true	"Login credentials"
//	@Router			/auth/sign-in/otp/send [post]
//	@Success		200	{object}	responses.responses{responses.OTPResponse{}}	"Successfully otp send to user's registered number"
//	@Failure		400	{object}	responses.responses{}							"Invalid Otp"
//	@Failure		403	{object}	responses.responses{}							"User blocked by admin"
//	@Failure		401	{object}	responses.responses{}							"User not exist with given login credentials"
//	@Failure		500	{object}	responses.responses{}							"Failed to send otp"
func (u *AuthHandler) UserLoginOtpSend(ctx *gin.Context) {

	var body requests.OTPLogin
	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	//check all input field is empty
	if body.Email == "" && body.Phone == "" && body.UserName == "" {
		err := errors.New("enter at least user_name or email or phone")
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	otpID, err := u.authUseCase.UserLoginOtpSend(ctx, body)

	if err != nil {
		var statusCode int

		switch {
		case errors.Is(err, usecases.ErrEmptyLoginCredentials):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrUserNotExist):
			statusCode = http.StatusForbidden
		case errors.Is(err, usecases.ErrUserBlocked):
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to send otp", err, nil)
		return
	}

	otpRes := responses.OTPResponse{
		OtpID: otpID,
	}
	responses.SuccessResponse(ctx, http.StatusOK, "Successfully otp send to user's registered number", otpRes)
}

// UserLoginOtpVerify godoc
//
//	@summary		Login with Otp verify (User)
//	@description	API for user to verify otp
//	@id				UserLoginOtpVerify
//	@tags			User Authentication
//	@param			inputs	body	requests.OTPVerify{}	true	"Otp Verify Details"
//	@Router			/auth/sign-in/otp/verify [post]
//	@Success		200	{object}	responses.responses{data=responses.TokenResponse}	"Successfully user logged in"
//	@Failure		400	{object}	responses.responses{}								"Invalid inputs"
//	@Failure		401	{object}	responses.responses{}								"Otp not matched"
//	@Failure		410	{object}	responses.responses{}								"Otp Expired"
//	@Failure		500	{object}	responses.responses{}								"Failed to verify otp
func (c *AuthHandler) UserLoginOtpVerify(ctx *gin.Context) {

	var body requests.OTPVerify
	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	// get the user using loginOtp usecases
	userID, err := c.authUseCase.LoginOtpVerify(ctx, body)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, usecases.ErrOtpExpired):
			statusCode = http.StatusGone
		case errors.Is(err, usecases.ErrInvalidOtp):
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to verify otp", err, nil)
		return
	}

	c.setupTokenAndResponse(ctx, tokens.User, userID)
}

// UserSignUp godoc
//
//	@Summary		Signup (User)
//	@Description	API for user to register a new account
//	@Id				UserSignUp
//	@Tags			User Authentication
//	@Param			input	body	requests.UserSignUp{}	true	"Input Fields"
//	@Router			/auth/sign-up [post]
//	@Success		200	{object}	responses.responses{data=responses.OTPResponse}	"Successfully account created and otp send to registered number"
//	@Failure		400	{object}	responses.responses{}								"Invalid input"
//	@Failure		409	{object}	responses.responses{}								"A verified user already exist with given user credentials"
//	@Failure		500	{object}	responses.responses{}								"Failed to signup"
func (c *AuthHandler) UserSignUp(ctx *gin.Context) {

	var body requests.UserSignUp

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	var user models.User
	if err := copier.Copy(&user, body); err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "failed to copy details", err, nil)
		return
	}

	err := c.authUseCase.UserSignUp(ctx, user)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecases.ErrUserAlreadyExit) {
			statusCode = http.StatusConflict
		}

		responses.ErrorResponse(ctx, statusCode, "Failed to signup", err, nil)
		return
	}

	// otpRes := responses.OTPResponse{
	// 	OtpID: otpID,
	// }

	responses.SuccessResponse(ctx, http.StatusCreated,
		"Successfully account created")
}

// UserSignUpVerify godoc
//
//	@summary		UserSingUp verify OTP  (User)
//	@description	API for user to verify otp on sign up
//	@id				UserSignUpVerify
//	@tags			User Authentication
//	@param			inputs	body	requests.OTPVerify{}	true	"Otp Verify Details"
//	@Router			/auth/sign-up/verify [post]
//	@Success		200	{object}	responses.responses{data=responses.TokenResponse}	"Successfully otp verified for user sign up"
//	@Failure		400	{object}	responses.responses{}								"Invalid inputs"
//	@Failure		401	{object}	responses.responses{}								"Otp not matched"
//	@Failure		410	{object}	responses.responses{}								"Otp Expired"
//	@Failure		500	{object}	responses.responses{}								"Failed to verify otp"
func (c *AuthHandler) UserSignUpVerify(ctx *gin.Context) {

	var body requests.OTPVerify
	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	// get the user using loginOtp usecases
	userID, err := c.authUseCase.SingUpOtpVerify(ctx, body)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, usecases.ErrOtpExpired):
			statusCode = http.StatusGone
		case errors.Is(err, usecases.ErrInvalidOtp):
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to verify otp", err, nil)
		return
	}

	c.setupTokenAndResponse(ctx, tokens.User, userID)
}

// AdminLogin godoc
//
//	@Summary		Login with password (Admin)
//	@Description	API for admin to login with password
//	@Id				AdminLogin
//	@Tags			Admin Authentication
//	@Param			input	body	requests.Login{}	true	"Login credentials"
//	@Router			/admin/auth/sign-in [post]
//	@Success		200	{object}	responses.responses{data=responses.TokenResponse}	"Successfully logged in"
//	@Failure		400	{object}	responses.responses{}								"Invalid input"
//	@Failure		401	{object}	responses.responses{}								"Wrong password"
//	@Failure		404	{object}	responses.responses{}								"Admin not exist with this details"
//	@Failure		500	{object}	responses.responses{}								"Failed to login"
func (c *AuthHandler) AdminLogin(ctx *gin.Context) {

	var body requests.Login

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	adminID, err := c.authUseCase.AdminLogin(ctx, body)
	if err != nil {

		var statusCode int

		switch {
		case errors.Is(err, usecases.ErrEmptyLoginCredentials):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrUserNotExist):
			statusCode = http.StatusNotFound
		case errors.Is(err, usecases.ErrWrongPassword):
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}

		responses.ErrorResponse(ctx, statusCode, "Failed to login", err, nil)
		return
	}

	// setup tokens common part
	c.setupTokenAndResponse(ctx, tokens.Admin, adminID)
}

// access and refresh tokens generating for user and admin is same so created
// a common function for it.(differentiate user by user type )
func (c *AuthHandler) setupTokenAndResponse(ctx *gin.Context, tokenUser tokens.UserType, userID uint) {

	tokenParams := usecaseInterface.GenerateTokenParams{
		UserID:   userID,
		UserType: tokenUser,
	}

	accessToken, err := c.authUseCase.GenerateAccessToken(ctx, tokenParams)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate access tokens", err, nil)
		return
	}

	refreshToken, err := c.authUseCase.GenerateRefreshToken(ctx, usecaseInterface.GenerateTokenParams{
		UserID:   userID,
		UserType: tokenUser,
	})
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate refresh tokens", err, nil)
		return
	}

	authorizationValue := authorizationType + " " + accessToken
	ctx.Header(authorizationHeaderKey, authorizationValue)

	ctx.Header("access_token", accessToken)
	ctx.Header("refresh_token", refreshToken)

	tokenRes := responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully logged in", tokenRes)
}

// UserRenewAccessToken godoc
//
//	@Summary		Renew Access tokens (User)
//	@Description	API for user to renew access tokens using refresh tokens
//	@Security		ApiKeyAuth
//	@Id				UserRenewAccessToken
//	@Tags			User Authentication
//	@Param			input	body	requests.RefreshToken{}	true	"Refresh tokens"
//	@Router			/auth/renew-access-tokens [post]
//	@Success		200	{object}	responses.responses{}	"Successfully generated access tokens using refresh tokens"
//	@Failure		400	{object}	responses.responses{}	"Invalid input"
//	@Failure		401	{object}	responses.responses{}	"Invalid refresh tokens"
//	@Failure		404	{object}	responses.responses{}	"No session found for the given refresh tokens"
//	@Failure		410	{object}	responses.responses{}	"Refresh tokens expired"
//	@Failure		403	{object}	responses.responses{}	"Refresh tokens blocked"
//	@Failure		500	{object}	responses.responses{}	"Failed generate access tokens"
func (c *AuthHandler) UserRenewAccessToken() gin.HandlerFunc {
	return c.renewAccessToken(tokens.User)
}

// AdminRenewAccessToken godoc
//
//	@Summary		Renew Access tokens (Admin)
//	@Description	API for admin to renew access tokens using refresh tokens
//	@Security		ApiKeyAuth
//	@Id				AdminRenewAccessToken
//	@Tags			Admin Authentication
//	@Param			input	body	requests.RefreshToken{}	true	"Refresh tokens"
//	@Router			/admin/auth/renew-access-tokens [post]
//	@Success		200	{object}	responses.responses{}	"Successfully generated access tokens using refresh tokens"
//	@Failure		400	{object}	responses.responses{}	"Invalid input"
//	@Failure		401	{object}	responses.responses{}	"Invalid refresh tokens"
//	@Failure		404	{object}	responses.responses{}	"No session found for the given refresh tokens"
//	@Failure		410	{object}	responses.responses{}	"Refresh tokens expired"
//	@Failure		403	{object}	responses.responses{}	"Refresh tokens blocked"
//	@Failure		500	{object}	responses.responses{}	"Failed generate access tokens"
func (c *AuthHandler) AdminRenewAccessToken() gin.HandlerFunc {
	return c.renewAccessToken(tokens.Admin)
}

// common functionality of renewing access tokens for user and admin
func (c *AuthHandler) renewAccessToken(tokenUser tokens.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var body requests.RefreshToken

		if err := ctx.ShouldBindJSON(&body); err != nil {
			responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
			return
		}

		refreshSession, err := c.authUseCase.VerifyAndGetRefreshTokenSession(ctx, body.RefreshToken, tokenUser)

		if err != nil {
			var statusCode int

			switch {
			case errors.Is(err, usecases.ErrInvalidRefreshToken):
				statusCode = http.StatusUnauthorized
			case errors.Is(err, usecases.ErrRefreshSessionNotExist):
				statusCode = http.StatusNotFound
			case errors.Is(err, usecases.ErrRefreshSessionExpired):
				statusCode = http.StatusGone
			case errors.Is(err, usecases.ErrRefreshSessionBlocked):
				statusCode = http.StatusForbidden
			default:
				statusCode = http.StatusInternalServerError
			}
			responses.ErrorResponse(ctx, statusCode, "Failed verify refresh tokens", err, nil)
			return
		}

		accessTokenParams := usecaseInterface.GenerateTokenParams{
			UserID:   refreshSession.UserID,
			UserType: tokenUser,
		}

		accessToken, err := c.authUseCase.GenerateAccessToken(ctx, accessTokenParams)

		if err != nil {
			responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed generate access tokens", err, nil)
			return
		}
		cookieName := "auth-" + string(tokenUser)
		ctx.SetCookie(cookieName, accessToken, 15*60, "", "", false, true)

		accessTokenRes := responses.TokenResponse{
			AccessToken: accessToken,
		}
		responses.SuccessResponse(ctx, http.StatusOK, "Successfully generated access tokens using refresh tokens", accessTokenRes)
	}
}
