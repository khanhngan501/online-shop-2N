package handlers

import (
	"net/http"
	handlerInterface "online-shop-2N/pkg/api/handlers/interfaces"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	commonConstant "online-shop-2N/pkg/common/constants"
	"online-shop-2N/pkg/usecases/interfaces"
	"online-shop-2N/pkg/utils"

	"github.com/gin-gonic/gin"
)

type paymentHandler struct {
	paymentUseCase interfaces.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase interfaces.PaymentUseCase) handlerInterface.PaymentHandler {
	return &paymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

// CartOrderPaymentSelectPage godoc
//
//	@Summary		Render Payment Page (User)
//	@Security		BearerAuth
//	@Description	API for user to render payment select page
//	@Id				CartOrderPaymentSelectPage
//	@Tags			User Payment
//	@Router			/carts/checkout/payment-select-page [get]
//	@Success		200	{object}	responses.responses{}	"Successfully rendered payment page"
//	@Failure		500	{object}	responses.responses{}	"Failed to render payment page"
func (c *paymentHandler) CartOrderPaymentSelectPage(ctx *gin.Context) {

	Payments, err := c.paymentUseCase.FindAllPaymentMethods(ctx)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to render payment page", err, nil)
		return
	}

	ctx.HTML(200, "paymentForm.html", Payments)
}

// UpdatePaymentMethod godoc
//
//	@Summary		Update payment method (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to change maximum price or block or unblock the payment method
//	@tags			Admin Payment
//	@id				UpdatePaymentMethod
//	@Router			/admin/payment-method  [put]
//	@Success		200	{object}	responses.responses{}	"Successfully payment method updated"
//	@Success		400	{object}	responses.responses{}	"Invalid inputs"
//	@Failure		500	{object}	responses.responses{}	"Failed to update payment method"
func (c *paymentHandler) UpdatePaymentMethod(ctx *gin.Context) {

	var body requests.PaymentMethodUpdate

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.paymentUseCase.UpdatePaymentMethod(ctx, body)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update payment method", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully payment method updated", nil)
}

// GetAllPaymentMethodsAdmin godoc
//
//	@summary		Get payment methods (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all payment methods
//	@tags			Admin Payment
//	@id				GetAllPaymentMethodsAdmin
//	@Router			/admin/payment-methods [get]
//	@Success		200	{object}	responses.responses{}	"Failed to retrieve payment methods"
//	@Failure		500	{object}	responses.responses{}	"Successfully retrieved all payment methods"
func (c *paymentHandler) GetAllPaymentMethodsAdmin() func(ctx *gin.Context) {
	return c.findAllPaymentMethods()
}

// GetAllPaymentMethodsUser godoc
//
//	@summary		Get payment methods (User)
//	@Security		BearerAuth
//	@Description	API for user to get all payment methods
//	@tags			User Payment
//	@id				GetAllPaymentMethodsUser
//	@Router			/payment-methods [get]
//	@Success		200	{object}	responses.responses{}	"Failed to retrieve payment methods"
//	@Failure		500	{object}	responses.responses{}	"Successfully retrieved all payment methods"
func (c *paymentHandler) GetAllPaymentMethodsUser() func(ctx *gin.Context) {
	return c.findAllPaymentMethods()
}

func (c *paymentHandler) findAllPaymentMethods() func(ctx *gin.Context) {

	return func(ctx *gin.Context) {

		paymentMethods, err := c.paymentUseCase.FindAllPaymentMethods(ctx)
		if err != nil {
			responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve payment methods", err, nil)
			return
		}

		if paymentMethods == nil {
			responses.SuccessResponse(ctx, http.StatusOK, "No payment methods found")
			return
		}

		responses.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all payment methods", paymentMethods)
	}
}

// PaymentCOD godoc
//
//	@summary		Place order  for COD (User)
//	@Security		BearerAuth
//	@Description	API for user to place order for cash on delivery
//	@tags			User Payment
//	@id				PaymentCOD
//	@Param			shop_order_id	formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/cod [post]
//	@Success		200	{object}	responses.responses{}	"successfully order placed for COD"
//	@Failure		500	{object}	responses.responses{}	"Failed place order for COD"
func (c *paymentHandler) PaymentCOD(ctx *gin.Context) {

	shopOrderID, err := requests.GetFormValuesAsUint(ctx, "shop_order_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	UserID := utils.GetUserIdFromContext(ctx)

	approveReq := requests.ApproveOrder{
		ShopOrderID: shopOrderID,
		PaymentType: commonConstant.CodPayment,
	}

	// approve the order and clear the user cart
	err = c.paymentUseCase.ApproveShopOrderAndClearCart(ctx, UserID, approveReq)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to approve order and clear cart", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully order placed for cod")
}
