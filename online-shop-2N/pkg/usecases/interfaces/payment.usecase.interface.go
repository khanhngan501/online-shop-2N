package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type PaymentUseCase interface {
	FindAllPaymentMethods(ctx context.Context) ([]models.PaymentMethod, error)
	FindPaymentMethodByID(ctx context.Context, paymentMethodID uint) (models.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, paymentMethod requests.PaymentMethodUpdate) error

	// razorpay
	MakeRazorpayOrder(ctx context.Context, userID, shopOrderID uint) (razorpayOrder responses.RazorpayOrder, err error)
	VerifyRazorPay(ctx context.Context, verifyReq requests.RazorpayVerify) error
	// stipe
	MakeStripeOrder(ctx context.Context, userID, shopOrderID uint) (stipeOrder responses.StripeOrder, err error)
	VerifyStripOrder(ctx context.Context, stripePaymentID string) error

	ApproveShopOrderAndClearCart(ctx context.Context, userID uint, approveDetails requests.ApproveOrder) error
}
