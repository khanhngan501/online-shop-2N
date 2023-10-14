package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type AdminRepository interface {
	FindAdminByEmail(ctx context.Context, email string) (models.Admin, error)
	FindAdminByUserName(ctx context.Context, userName string) (models.Admin, error)
	SaveAdmin(ctx context.Context, admin models.Admin) error

	FindAllUser(ctx context.Context, pagination requests.Pagination) (users []responses.User, err error)

	CreateFullSalesReport(ctc context.Context, reqData requests.SalesReport) (salesReport []responses.SalesReport, err error)

	//stock side
	FindStockBySKU(ctx context.Context, sku string) (stock responses.Stock, err error)
}
