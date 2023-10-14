package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
)

type StockRepository interface {
	FindAll(ctx context.Context, pagination requests.Pagination) (stocks []responses.Stock, err error)
	Update(ctx context.Context, updateValues requests.UpdateStock) error
}
