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

// RazorpayCheckout godoc
//
//	@Summary		Razorpay checkout (User)
//	@Security		BearerAuth
//	@Description	API for user to create stripe payment
//	@Security		ApiKeyAuth
//	@Tags			User Payment
//	@Id				RazorpayCheckout
//	@Param			shop_order_id	formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/razorpay-checkout [post]
//	@Success		200	{object}	responses.responses{}	"successfully razorpay payment order created"
//	@Failure		500	{object}	responses.responses{}	"Failed to make razorpay order"
func (c *paymentHandler) RazorpayCheckout(ctx *gin.Context) {

	shopOrderID, err := requests.GetFormValuesAsUint(ctx, "shop_order_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	UserID := utils.GetUserIdFromContext(ctx)

	razorpayOrder, err := c.paymentUseCase.MakeRazorpayOrder(ctx, UserID, shopOrderID)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to make razorpay order ", err, nil)
		return
	}

	razorPayRes := responses.OrderPayment{
		PaymentType:  commonConstant.RazopayPayment,
		PaymentOrder: razorpayOrder,
	}
	ctx.JSON(http.StatusOK, razorPayRes)
	// responses.SuccessResponse(ctx, http.StatusCreated, "Successfully razor pay order created", razorPayRes)
}

// RazorpayVerify godoc
//
//	@Summary		Razorpay verify (User)
//	@Security		BearerAuth
//	@Description	API for razorpay to callback backend for payment verification
//	@tags			User Payment
//	@id				RazorpayVerify
//	@Param			razorpay_order_id	formData	string	true	"Razorpay payment id"
//	@Param			razorpay_payment_id	formData	string	true	"Razorpay payment id"
//	@Param			razorpay_signature	formData	string	false	"Razorpay signature"
//	@Param			shop_order_id		formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/razorpay-verify [post]
//	@Success		200	{object}	responses.responses{}	"Successfully razorpay payment verified"
//	@Failure		402	{object}	responses.responses{}	"Payment not approved"
//	@Failure		500	{object}	responses.responses{}	"Failed to Approve order"
func (c *paymentHandler) RazorpayVerify(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	razorpayOrderID, err2 := requests.GetFormValuesAsString(ctx, "razorpay_order_id")
	razorpayPaymentID, err1 := requests.GetFormValuesAsString(ctx, "razorpay_payment_id")
	razorpaySignature, err3 := requests.GetFormValuesAsString(ctx, "razorpay_signature")

	shopOrderID, err4 := requests.GetFormValuesAsUint(ctx, "shop_order_id")

	err := errors.Join(err1, err2, err3, err4)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	verifyReq := requests.RazorpayVerify{
		OrderID:   razorpayOrderID,
		PaymentID: razorpayPaymentID,
		Signature: razorpaySignature,
	}

	err = c.paymentUseCase.VerifyRazorPay(ctx, verifyReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecases.ErrPaymentNotApproved) {
			statusCode = http.StatusPaymentRequired
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to verify razorpay payment", err, nil)
		return
	}

	approveReq := requests.ApproveOrder{
		ShopOrderID: shopOrderID,
		PaymentType: commonConstant.RazopayPayment,
	}

	err = c.paymentUseCase.ApproveShopOrderAndClearCart(ctx, userID, approveReq)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to Approve order", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully razorpay payment verified", nil)
}
