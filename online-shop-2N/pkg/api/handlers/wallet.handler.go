package handlers

import (
	"net/http"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetUserWallet godoc
//
//	@Summary		Get user wallet  (User)
//	@Security		BearerAuth
//	@Description	API for user to get user wallet
//	@Id				GetUserWallet
//	@Tags			User Profile
//	@Router			/account/wallet [get]
//	@Success		200	{object}	responses.Response{}	"Successfully retrieve user wallet"
//	@Failure		500	{object}	responses.Response{}	"Failed to retrieve user wallet"
func (c *OrderHandler) GetUserWallet(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	wallet, err := c.orderUseCase.FindUserWallet(ctx, userID)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve user wallet", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully retrieve user wallet", wallet)
}

// GetUserWalletTransactions godoc
//
//	@Summary		Get user wallet  (User)
//	@Security		BearerAuth
//	@Description	API for user to get user wallet transaction
//	@Id				GetUserWalletTransactions
//	@Tags			User Profile
//	@Router			/account/wallet/transactions [get]
//	@Success		200	{object}	responses.Response{}	"Successfully retrieved user wallet transactions"
//	@Success		204	{object}	responses.Response{}	"No wallet transaction for user"
//	@Failure		500	{object}	responses.Response{}	"Failed to retrieve user wallet transactions"
func (c *OrderHandler) GetUserWalletTransactions(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	pagination := requests.GetPagination(ctx)

	transactions, err := c.orderUseCase.FindUserWalletTransactions(ctx, userID, pagination)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve user wallet transactions", err, nil)
		return
	}

	if len(transactions) == 0 {
		responses.SuccessResponse(ctx, http.StatusNoContent, "No user wallet transaction found", nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved user wallet transactions", transactions)
}
