package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type OfferRepository interface {
	Transactions(ctx context.Context, trxFn func(repo OfferRepository) error) error

	// offer
	FindOfferByID(ctx context.Context, offerID uint) (models.Offer, error)
	FindOfferByName(ctx context.Context, offerName string) (models.Offer, error)
	FindAllOffers(ctx context.Context, pagination requests.Pagination) ([]models.Offer, error)
	SaveOffer(ctx context.Context, offer requests.Offer) error
	DeleteOffer(ctx context.Context, offerID uint) error

	// to calculate the discount price and update
	UpdateProductsDiscountByCategoryOfferID(ctx context.Context, categoryOfferID uint) error
	UpdateProductItemsDiscountByCategoryOfferID(ctx context.Context, categoryOfferID uint) error
	UpdateProductsDiscountByProductOfferID(ctx context.Context, productOfferID uint) error
	UpdateProductItemsDiscountByProductOfferID(ctx context.Context, productOfferID uint) error

	// to remove the discount product price
	RemoveProductsDiscountByCategoryOfferID(ctx context.Context, categoryOfferID uint) error
	RemoveProductItemsDiscountByCategoryOfferID(ctx context.Context, categoryOfferID uint) error
	RemoveProductsDiscountByProductOfferID(ctx context.Context, productOfferID uint) error
	RemoveProductItemsDiscountByProductOfferID(ctx context.Context, productOfferID uint) error

	// offer category
	FindOfferCategoryCategoryID(ctx context.Context, categoryID uint) (models.OfferCategory, error)
	FindAllOfferCategories(ctx context.Context, pagination requests.Pagination) ([]responses.OfferCategory, error)

	SaveCategoryOffer(ctx context.Context, categoryOffer requests.OfferCategory) (categoryOfferID uint, err error)
	DeleteCategoryOffer(ctx context.Context, categoryOfferID uint) error
	UpdateCategoryOffer(ctx context.Context, categoryOfferID, offerID uint) error

	// offer products
	FindOfferProductByProductID(ctx context.Context, productID uint) (models.OfferProduct, error)
	FindAllOfferProducts(ctx context.Context, pagination requests.Pagination) ([]responses.OfferProduct, error)

	SaveOfferProduct(ctx context.Context, offerProduct models.OfferProduct) (productOfferId uint, err error)
	DeleteOfferProduct(ctx context.Context, productOfferID uint) error
	UpdateOfferProduct(ctx context.Context, productOfferID, offerID uint) error

	DeleteAllProductOffersByOfferID(ctx context.Context, offerID uint) error
	DeleteAllCategoryOffersByOfferID(ctx context.Context, offerID uint) error
}
