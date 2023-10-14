package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type AdminUseCase interface {
	SignUp(ctx context.Context, admin models.Admin) error

	FindAllUser(ctx context.Context, pagination requests.Pagination) (users []responses.User, err error)
	BlockOrUnBlockUser(ctx context.Context, blockDetails requests.BlockUser) error

	GetFullSalesReport(ctx context.Context, requestData requests.SalesReport) (salesReport []responses.SalesReport, err error)
}
