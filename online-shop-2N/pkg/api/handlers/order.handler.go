package handlers

import (
	"errors"
	"net/http"
	"online-shop-2N/pkg/api/handlers/interfaces"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/usecases"
	usecaseInterface "online-shop-2N/pkg/usecases/interfaces"
	"online-shop-2N/pkg/utils"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase usecaseInterface.OrderUseCase
}

func NewOrderHandler(orderUseCase usecaseInterface.OrderUseCase) interfaces.OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

// GetAllOrderStatuses godoc
//
//	@Summary		Get all order statuses (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all available order statuses
//	@Id				GetAllOrderStatuses
//	@Tags			Admin Orders
//	@Router			/admin/orders/statuses [get]
//	@Success		200	{object}	responses.Response{}	"Successfully retrieved all order statuses"
//	@Success		204	{object}	responses.Response{}	"No order statuses found"
//	@Failure		500	{object}	responses.Response{}	"Failed to find all order statuses"
func (c *OrderHandler) GetAllOrderStatuses(ctx *gin.Context) {

	orderStatuses, err := c.orderUseCase.FindAllOrderStatuses(ctx)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to find all order statuses", err, nil)
		return
	}

	if orderStatuses == nil {
		responses.SuccessResponse(ctx, 200, "No order statuses found")
		return
	}

	responses.SuccessResponse(ctx, 200, "Successfully retrieved all order statuses", orderStatuses)
}

// SaveOrder godoc
//
//	@Summary		Save Order (User)
//	@Security		BearerAuth
//	@Description	API for user save an order
//	@Tags			User Orders
//	@Id				SaveOrder
//	@Param			address_id	formData	string	true	"Address ID"
//	@Router			/carts/place-order [post]
//	@Success		200	{object}	responses.Response{}	"successfully order placed"
//	@Success		204	{object}	responses.Response{}	"Cart is empty"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
//	@Failure		409	{object}	responses.Response{}	"Can't place order out of stock product on cart"
//	@Failure		500	{object}	responses.Response{}	"Failed to save order"
func (c *OrderHandler) SaveOrder(ctx *gin.Context) {

	addressID, err := requests.GetFormValuesAsUint(ctx, "address_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	shopOrderID, err := c.orderUseCase.SaveOrder(ctx, userID, addressID)

	if err != nil {
		var statusCode int

		switch {
		case errors.Is(err, usecases.ErrEmptyCart):
			statusCode = http.StatusNoContent
		case errors.Is(err, usecases.ErrOutOfStockOnCart):
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to save order", err, nil)
		return
	}

	data := gin.H{
		"shop_order_id": shopOrderID,
	}
	responses.SuccessResponse(ctx, http.StatusOK, "Successfully order saved", data)
}

// GetUserOrder godoc
//
//	@summary		Get user orders (User)
//	@Security		BearerAuth
//	@description	API to get order for user user orders
//	@id				GetUserOrder
//	@tags			User Orders
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count Of Order"
//	@Router			/orders [get]
//	@Success		200	{object}	responses.Response{}	"Successfully retrieved all user orders"
//	@Success		204	{object}	responses.Response{}	"No shop orders for user"
//	@Failure		500	{object}	responses.Response{}	"Failed to retrieve all user orders"
func (c *OrderHandler) GetUserOrder(ctx *gin.Context) {

	userId := utils.GetUserIdFromContext(ctx)
	pagination := requests.GetPagination(ctx)

	orders, err := c.orderUseCase.FindUserShopOrder(ctx, userId, pagination)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve all user shop orders", err, nil)
		return
	}

	if orders == nil {
		responses.SuccessResponse(ctx, http.StatusNoContent, "No shop orders found", nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all user orders", orders)
}

// GetAllShopOrders godoc
//
//	@Summary		Get all orders (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all orders
//	@Id				GetAllShopOrders
//	@Tags			Admin Orders
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/orders/all [get]
//	@Success		200	{object}	responses.Response{}	"Successfully retrieved all shop orders"
//	@Success		204	{object}	responses.Response{}	"No shop order found"
//	@Failure		500	{object}	responses.Response{}	"Failed to find all shop orders"
func (c *OrderHandler) GetAllShopOrders(ctx *gin.Context) {

	pagination := requests.GetPagination(ctx)

	shopOrders, err := c.orderUseCase.FindAllShopOrders(ctx, pagination)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all shop orders", err, nil)
		return
	}

	if len(shopOrders) == 0 {
		responses.SuccessResponse(ctx, http.StatusNoContent, "No shop order found", nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all shop orders", shopOrders)
}

// GetAllOrderItemsUser godoc
//
//	@Summary		Get all order items (User)
//	@Security		BearerAuth
//	@Description	API for user to get all order items of a specific order
//	@Id				GetAllOrderItemsUser
//	@Tags			User Orders
//	@Param			shop_order_id	path	int	true	"Shop Order ID"
//	@Param			page_number		query	int	false	"Page Number"
//	@Param			count			query	int	false	"Count Of Order"
//	@Router			/orders/{shop_order_id}/items  [get]
//	@Success		200	{object}	responses.Response{}	"Successfully found order items"
//	@Failure		500	{object}	responses.Response{}	"Failed to find order items"
func (c *OrderHandler) GetAllOrderItemsUser() func(ctx *gin.Context) {
	return c.findAllOrderItems()
}

// GetAllOrderItemsAdmin godoc
//
//	@Summary		Get all order items (Admin)
//	@Security		BearerAuth
//	@Description	API for user to get all order items of a specific order
//	@Id				GetAllOrderItemsAdmin
//	@Tags			Admin Orders
//	@Param			shop_order_id	path	int	true	"Shop Order ID"
//	@Param			page_number		query	int	false	"Page Number"
//	@Param			count			query	int	false	"Count"
//	@Router			/admin/orders/{shop_order_id}/items [get]
//	@Success		200	{object}	responses.Response{}	"Successfully found order items"
//	@Success		204	{object}	responses.Response{}	"No order items found"
//	@Failure		500	{object}	responses.Response{}	"Failed to find order items"
func (c *OrderHandler) GetAllOrderItemsAdmin() func(ctx *gin.Context) {
	return c.findAllOrderItems()
}

func (c *OrderHandler) findAllOrderItems() func(ctx *gin.Context) {

	return func(ctx *gin.Context) {
		shopOrderID, err := requests.GetParamAsUint(ctx, "shop_order_id")
		if err != nil {
			responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		}
		pagination := requests.GetPagination(ctx)

		orderItems, err := c.orderUseCase.FindOrderItems(ctx, shopOrderID, pagination)

		if err != nil {
			responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find order items", err, nil)
			return
		}

		if orderItems == nil {
			responses.SuccessResponse(ctx, http.StatusNoContent, "No order items found", nil)
			return
		}

		responses.SuccessResponse(ctx, http.StatusOK, "Successfully found order items", orderItems)
	}
}

// CancelOrder godoc
//
//	@Summary		Cancel order (User)
//	@Security		BearerAuth
//	@Description	Api for user to cancel a order
//	@Id				CancelOrder
//	@Tags			User Orders
//	@Param			shop_order_id	path	int	true	"Shop Order ID"
//	@Router			/orders/{shop_order_id}/cancel [post]
//	@Success		200	{object}	responses.Response{}	"Successfully order cancelled"
//	@Failure		400	{object}	responses.Response{}	"Invalid inputs"
//	@Failure		500	{object}	responses.Response{}	"Failed to cancel order"
func (c *OrderHandler) CancelOrder(ctx *gin.Context) {

	shopOrderID, err := requests.GetParamAsUint(ctx, "shop_order_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
	}

	err = c.orderUseCase.CancelOrder(ctx, shopOrderID)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to cancel order", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "successfully order cancelled", nil)
}

// UpdateOrderStatus godoc
//
//	@Summary		Change order status (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to change order status
//	@Id				UpdateOrderStatus
//	@Tags			Admin Orders
//	@Param			input	body	requests.UpdateOrder{}	true	"input field"
//	@Router			/admin/orders/ [put]
//	@Success		200	{object}	responses.Response{}	"Successfully order status updated"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c *OrderHandler) UpdateOrderStatus(ctx *gin.Context) {

	var body requests.UpdateOrder

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.orderUseCase.UpdateOrderStatus(ctx, body.ShopOrderID, body.OrderStatusID)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update order status", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully order status updated", nil)
}

// SubmitReturnRequest godoc
//
//	@Summary		Return requests (User)
//	@Security		BearerAuth
//	@Description	API for user to requests a return for delivered order
//	@Id				SubmitReturnRequest
//	@Tags			User Orders
//	@Param			input	body	requests.Return	true	"Input Fields"
//	@Router			/orders/return [post]
//	@Success		200	{object}	responses.Response{}	"Successfully return requests submitted for order"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c OrderHandler) SubmitReturnRequest(ctx *gin.Context) {

	var body requests.Return
	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.orderUseCase.SubmitReturnRequest(ctx, body)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to submit return requests", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully return requests submitted for order", nil)
}

// GetAllOrderReturns godoc
//
//	@Summary		Get all order returns (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all order returns
//	@Id				GetAllOrderReturns
//	@Tags			Admin Orders
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count Of Order"
//	@Router			/admin/orders/returns [get]
//	@Success		200	{object}	responses.Response{}	"Successfully found all order returns"
//	@Failure		500	{object}	responses.Response{}	"Failed to find all order returns"
func (c *OrderHandler) GetAllOrderReturns(ctx *gin.Context) {

	pagination := requests.GetPagination(ctx)

	orderReturns, err := c.orderUseCase.FindAllOrderReturns(ctx, pagination)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all order returns", err, nil)
		return
	}

	if len(orderReturns) == 0 {
		responses.SuccessResponse(ctx, http.StatusOK, "No order returns found", nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully found all order returns", orderReturns)
}

// GetAllPendingReturns godoc
//
//	@Summary		Get all pending returns (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all pending returns
//	@Id				GetAllPendingReturns
//	@Tags			Admin Orders
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count Of Order"
//	@Router			/admin/orders/returns/pending [get]
//	@Success		200	{object}	responses.Response{}	"Successfully found all pending orders return requests"
//	@Failure		500	{object}	responses.Response{}	"Failed to find all pending order return requests"
func (c *OrderHandler) GetAllPendingReturns(ctx *gin.Context) {

	pagination := requests.GetPagination(ctx)

	orderReturns, err := c.orderUseCase.FindAllPendingOrderReturns(ctx, pagination)
	if err != nil {
		responses.ErrorResponse(ctx, 500, "Failed to find all pending order return requests", err, nil)
		return
	}

	if len(orderReturns) == 0 {
		responses.SuccessResponse(ctx, 200, "No pending order returns requests found", nil)
		return
	}

	responses.SuccessResponse(ctx, 200, "Successfully found all pending orders return requests", orderReturns)
}

// UpdateReturnRequest godoc
//
//	@summary		Change return requests status (Admin)
//	@Security		BearerAuth
//	@description	API for admin to change status of return requested orders
//	@id				UpdateReturnRequest
//	@tags			Admin Orders
//	@Param			input	body	requests.UpdateOrderReturn{}	true	"Input Fields"
//	@Router			/admin/orders/returns/pending [put]
//	@Success		200	{object}	responses.Response{}	"successfully order_response updated"
//	@Failure		500	{object}	responses.Response{}	"invalid input"
func (c *OrderHandler) UpdateReturnRequest(ctx *gin.Context) {

	var body requests.UpdateOrderReturn

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.orderUseCase.UpdateReturnDetails(ctx, body)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update order return", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "successfully order return updated")
}
