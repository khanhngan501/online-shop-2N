package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
)

type StockUseCase interface {
	GetAllStockDetails(ctx context.Context, pagination requests.Pagination) (stocks []responses.Stock, err error)
	UpdateStockBySKU(ctx context.Context, updateDetails requests.UpdateStock) error
}
