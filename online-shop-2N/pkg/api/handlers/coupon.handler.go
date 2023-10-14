package handlers

import (
	"net/http"
	"online-shop-2N/pkg/api/handlers/interfaces"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
	usecase "online-shop-2N/pkg/usecases/interfaces"
	"online-shop-2N/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type CouponHandler struct {
	couponUseCase usecase.CouponUseCase
}

func NewCouponHandler(couponUseCase usecase.CouponUseCase) interfaces.CouponHandler {
	return &CouponHandler{couponUseCase: couponUseCase}
}

// SaveCoupon godoc
//
//	@Summary		Add coupons (Admin)
//	@Description	API for admin to add a new coupon
//	@Security		BearerAuth
//	@Tags			Admin Coupon
//	@Id				SaveCoupon
//	@Param			inputs	body	requests.Coupon{}	true	"Input Fields"
//	@Router			/admin/coupons [post]
//	@Success		200	{object}	responses.Response{}	"successfully coupon added"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c *CouponHandler) SaveCoupon(ctx *gin.Context) {

	var body requests.Coupon

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var coupon models.Coupon

	copier.Copy(&coupon, &body)

	err := c.couponUseCase.AddCoupon(ctx, coupon)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add coupon", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully coupon added")
}

// GetAllCouponsAdmin godoc
//
//	@Summary		Get all coupons (Admin)
//	@Description	API for admin to get all coupons
//	@Security		BearerAuth
//	@Tags			Admin Coupon
//	@Id				GetAllCouponsAdmin
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/coupons [get]
//	@Success		200	{object}	responses.Response{}	"successfully go all the coupons
//	@Failure		500	{object}	responses.Response{}	"failed to get all coupons"
func (c *CouponHandler) GetAllCouponsAdmin(ctx *gin.Context) {

	pagination := requests.GetPagination(ctx)

	coupons, err := c.couponUseCase.GetAllCoupons(ctx, pagination)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get all coupons", err, nil)
		return
	}

	if len(coupons) == 0 {
		responses.SuccessResponse(ctx, http.StatusOK, "No Coupons found", nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully found coupons", coupons)
}

// GetAllCouponsForUser godoc
//
//	@Summary		Get all user coupons (User)
//	@Description	API for user to get all coupons
//	@Security		BearerAuth
//	@tags			User Profile
//	@id				GetAllCouponsForUser
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count Of Order"
//	@Router			/account/coupons [get]
//	@Success		200	{object}	responses.Response{}	""Successfully	found	all	coupons	for	user"
//	@Failure		500	{object}	responses.Response{}	"Failed to find all user"
func (c *CouponHandler) GetAllCouponsForUser(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	pagination := requests.GetPagination(ctx)

	coupons, err := c.couponUseCase.GetCouponsForUser(ctx, userID, pagination)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to find all user", err, nil)
		return
	}

	if len(coupons) == 0 {
		responses.SuccessResponse(ctx, http.StatusOK, "No coupons found", nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully found all coupons for user", coupons)
}

// UpdateCoupon godoc
//
//	@Summary		Update Coupon (Admin)
//	@Description	API for admin update coupon details
//	@Security		BearerAuth
//	@Tags			Admin Coupon
//	@Id				UpdateCoupon
//	@Param			inputs	body	requests.EditCoupon{}	true	"Input Field"
//	@Router			/admin/coupons [put]
//	@Success		200	{object}	responses.Response{}	"Successfully updated the coupon"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c *CouponHandler) UpdateCoupon(ctx *gin.Context) {

	var body requests.EditCoupon

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var coupon models.Coupon

	copier.Copy(&coupon, &body)

	err := c.couponUseCase.UpdateCoupon(ctx, coupon)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update coupon", err, coupon)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully updated the coupon", coupon)
}

// ApplyCouponToCart godoc
//
//	@Summary		Apply coupon
//	@Description	API for user to apply a coupon on cart
//	@Security		BearerAuth
//	@Tags			User Cart
//	@Id				ApplyCouponToCart
//	@Param			inputs	body	requests.ApplyCoupon{}	true	"Input Field"
//	@Router			/carts/apply-coupon [patch]
//	@Success		200	{object}	responses.Response{}	"Successfully coupon applied to user cart"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c *CouponHandler) ApplyCouponToCart(ctx *gin.Context) {

	var body requests.ApplyCoupon

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	discountPrice, err := c.couponUseCase.ApplyCouponToCart(ctx, userID, body.CouponCode)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to apply the coupon code", err, nil)
		return
	}

	data := gin.H{"discount_amount": discountPrice}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully coupon applied to user cart", data)
}
