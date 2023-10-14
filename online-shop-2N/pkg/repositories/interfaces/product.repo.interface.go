package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type ProductRepository interface {
	Transactions(ctx context.Context, trxFn func(repo ProductRepository) error) error

	// variation
	IsVariationNameExistForCategory(ctx context.Context, name string, categoryID uint) (bool, error)
	SaveVariation(ctx context.Context, categoryID uint, variationName string) error
	FindAllVariationsByCategoryID(ctx context.Context, categoryID uint) ([]responses.Variation, error)

	// variation values
	IsVariationValueExistForVariation(ctx context.Context, value string, variationID uint) (exist bool, err error)
	SaveVariationOption(ctx context.Context, variationID uint, variationValue string) error
	FindAllVariationOptionsByVariationID(ctx context.Context, variationID uint) ([]responses.VariationOption, error)

	FindAllVariationValuesOfProductItem(ctx context.Context, productItemID uint) ([]responses.ProductVariationValue, error)
	//product
	FindProductByID(ctx context.Context, productID uint) (product models.Product, err error)
	IsProductNameExistForOtherProduct(ctx context.Context, name string, productID uint) (bool, error)
	IsProductNameExist(ctx context.Context, productName string) (exist bool, err error)

	FindAllProducts(ctx context.Context, pagination requests.Pagination) ([]responses.Product, error)
	SaveProduct(ctx context.Context, product models.Product) error
	UpdateProduct(ctx context.Context, product models.Product) error

	// product items
	FindProductItemByID(ctx context.Context, productItemID uint) (models.ProductItem, error)
	FindAllProductItems(ctx context.Context, productID uint) ([]responses.ProductItems, error)
	FindVariationCountForProduct(ctx context.Context, productID uint) (variationCount uint, err error) // to check the product config already exist
	FindAllProductItemIDsByProductIDAndVariationOptionID(ctx context.Context, productID, variationOptionID uint) ([]uint, error)
	SaveProductConfiguration(ctx context.Context, productItemID, variationOptionID uint) error
	SaveProductItem(ctx context.Context, productItem models.ProductItem) (productItemID uint, err error)
	// product item image
	FindAllProductItemImages(ctx context.Context, productItemID uint) (images []string, err error)
	SaveProductItemImage(ctx context.Context, productItemID uint, image string) error
}
