package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type ProductUseCase interface {
	// variations
	SaveVariation(ctx context.Context, categoryID uint, variationNames []string) error
	SaveVariationOption(ctx context.Context, variationID uint, variationOptionValues []string) error

	FindAllVariationsAndItsValues(ctx context.Context, categoryID uint) ([]responses.Variation, error)

	// products
	FindAllProducts(ctx context.Context, pagination requests.Pagination) (products []responses.Product, err error)
	SaveProduct(ctx context.Context, product requests.Product) error
	UpdateProduct(ctx context.Context, product models.Product) error

	SaveProductItem(ctx context.Context, productID uint, productItem requests.ProductItem) error
	FindAllProductItems(ctx context.Context, productID uint) ([]responses.ProductItems, error)
}
