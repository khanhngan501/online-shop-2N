package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type OfferUseCase interface {

	// offer
	SaveOffer(ctx context.Context, offer requests.Offer) error
	RemoveOffer(ctx context.Context, offerID uint) error
	FindAllOffers(ctx context.Context, pagination requests.Pagination) ([]models.Offer, error)

	// offer category
	SaveCategoryOffer(ctx context.Context, offerCategory requests.OfferCategory) error
	FindAllCategoryOffers(ctx context.Context, pagination requests.Pagination) ([]responses.OfferCategory, error)
	RemoveCategoryOffer(ctx context.Context, categoryOfferID uint) error
	ChangeCategoryOffer(ctx context.Context, categoryOfferID, offerID uint) error

	// offer product
	SaveProductOffer(ctx context.Context, offerProduct models.OfferProduct) error
	FindAllProductOffers(ctx context.Context, pagination requests.Pagination) ([]responses.OfferProduct, error)
	RemoveProductOffer(ctx context.Context, productOfferID uint) error
	ChangeProductOffer(ctx context.Context, productOfferID, offerID uint) error
}
