package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	commonConstant "online-shop-2N/pkg/common/constants"
	"online-shop-2N/pkg/models"
)

type PaymentRepository interface {
	FindPaymentMethodByID(ctx context.Context, paymentMethodID uint) (paymentMethods models.PaymentMethod, err error)
	FindPaymentMethodByType(ctx context.Context, paymentType commonConstant.PaymentType) (paymentMethod models.PaymentMethod, err error)
	FindAllPaymentMethods(ctx context.Context) ([]models.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, paymentMethod requests.PaymentMethodUpdate) error
}
