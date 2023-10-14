package handlers

import (
	"errors"
	"net/http"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/usecases"
	"online-shop-2N/pkg/utils"

	commonConstant "online-shop-2N/pkg/common/constants"

	"github.com/gin-gonic/gin"
)

// StripPaymentCheckout godoc
//
//	@Summary		Stripe checkout (User)
//	@Security		BearerAuth
//	@Description	API for user to create stripe payment
//	@Tags			User Payment
//	@Id				StripPaymentCheckout
//	@Param			shop_order_id	formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/stripe-checkout [post]
//	@Success		200	{object}	responses.responses{}	"successfully stripe payment order created"
//	@Failure		500	{object}	responses.responses{}	"Failed to create stripe order"
func (c *paymentHandler) StripPaymentCheckout(ctx *gin.Context) {

	shopOrderID, err := requests.GetFormValuesAsUint(ctx, "shop_order_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	UserID := utils.GetUserIdFromContext(ctx)

	stripeOrder, err := c.paymentUseCase.MakeStripeOrder(ctx, UserID, shopOrderID)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create stripe order", err, nil)
		return
	}

	stripeResponse := responses.OrderPayment{
		PaymentOrder: stripeOrder,
		PaymentType:  commonConstant.StripePayment,
	}

	ctx.JSON(http.StatusOK, stripeResponse)
}

// StripePaymentVeify godoc
//
//	@Summary		Stripe verify (User)
//	@Security		BearerAuth
//	@Description	API for user to callback backend after stripe payment for verification
//	@Tags			User Payment
//	@Id				StripePaymentVeify
//	@Param			stripe_payment_id	formData	string	true	"Stripe payment ID"
//	@Param			shop_order_id		formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/stripe-verify [post]
//	@Success		200	{object}	responses.responses{}	"Successfully stripe payment verified"
//	@Failure		402	{object}	responses.responses{}	"Payment not approved"
//	@Failure		500	{object}	responses.responses{}	"Failed to Approve order"
func (c *paymentHandler) StripePaymentVeify(ctx *gin.Context) {

	shopOrderID, err1 := requests.GetFormValuesAsUint(ctx, "shop_order_id")
	stripePaymentID, err2 := requests.GetFormValuesAsString(ctx, "stripe_payment_id")

	err := errors.Join(err1, err2)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	err = c.paymentUseCase.VerifyStripOrder(ctx, stripePaymentID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecases.ErrPaymentNotApproved) {
			statusCode = http.StatusPaymentRequired
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to verify stripe payment", err, nil)
		return
	}

	approveReq := requests.ApproveOrder{
		ShopOrderID: shopOrderID,
		PaymentType: commonConstant.StripePayment,
	}

	err = c.paymentUseCase.ApproveShopOrderAndClearCart(ctx, userID, approveReq)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to Approve order", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully stripe payment verified", nil)
}
