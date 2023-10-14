package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
)

type CategoryRepository interface {
	Transactions(ctx context.Context, trxFn func(repo CategoryRepository) error) error

	// category
	IsCategoryNameExist(ctx context.Context, categoryName string) (bool, error)
	FindAllMainCategories(ctx context.Context, pagination requests.Pagination) ([]responses.Category, error)
	SaveCategory(ctx context.Context, categoryName string) error

	// sub category
	IsSubCategoryNameExist(ctx context.Context, categoryName string, categoryID uint) (bool, error)
	FindAllSubCategories(ctx context.Context, categoryID uint) ([]responses.SubCategory, error)
	SaveSubCategory(ctx context.Context, categoryID uint, categoryName string) error
}
