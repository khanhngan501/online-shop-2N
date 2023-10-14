package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type CouponUseCase interface {
	// coupon
	AddCoupon(ctx context.Context, coupon models.Coupon) error
	GetAllCoupons(ctx context.Context, pagination requests.Pagination) (coupons []models.Coupon, err error)
	UpdateCoupon(ctx context.Context, coupon models.Coupon) error

	//user side coupons
	GetCouponsForUser(ctx context.Context, userID uint, pagination requests.Pagination) (coupons []responses.UserCoupon, err error)

	GetCouponByCouponCode(ctx context.Context, couponCode string) (coupon models.Coupon, err error)
	ApplyCouponToCart(ctx context.Context, userID uint, couponCode string) (discountPrice uint, err error)
}
