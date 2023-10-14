package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
)

type CategoryUseCase interface {
	FindAllCategories(ctx context.Context, pagination requests.Pagination) ([]responses.Category, error)
	SaveCategory(ctx context.Context, categoryName string) error
	SaveSubCategory(ctx context.Context, subCategory requests.SubCategory) error
}
