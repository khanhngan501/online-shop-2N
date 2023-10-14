package handlers

import (
	"errors"
	"net/http"
	"online-shop-2N/pkg/api/handlers/interfaces"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
	"online-shop-2N/pkg/usecases"
	usecaseInterface "online-shop-2N/pkg/usecases/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type offerHandler struct {
	offerUseCase usecaseInterface.OfferUseCase
}

func NewOfferHandler(offerUseCase usecaseInterface.OfferUseCase) interfaces.OfferHandler {
	return &offerHandler{
		offerUseCase: offerUseCase,
	}
}

// SaveOffer godoc
//
//	@Summary		Add offer (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to add an offer (Admin)
//	@Id				SaveOffer
//	@Tags			Admin Offers
//	@Param			input	body	requests.Offer{}	true	"input field"
//	@Router			/admin/offers [post]
//	@Success		200	{object}	responses.Response{}	"Successfully offer added"
//	@Failure		409	{object}	responses.Response{}	"Offer already exist"
//	@Failure		400	{object}	responses.Response{}	"Invalid inputs"
func (p *offerHandler) SaveOffer(ctx *gin.Context) {

	var body requests.Offer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.offerUseCase.SaveOffer(ctx, body)
	if err != nil {
		var statusCode int

		switch {
		case errors.Is(err, usecases.ErrOfferNameAlreadyExist):
			statusCode = http.StatusConflict
		case errors.Is(err, usecases.ErrInvalidOfferEndDate):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to add offer", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully offer added", nil)
}

// GetAllOffers godoc
//
//	@Summary		Get all offers (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all offers
//	@Id				GetAllOffers
//	@Tags			Admin Offers
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/offers [get]
//	@Success		200	{object}	responses.Response{}	""Successfully	found	all	offers"
//	@Failure		500	{object}	responses.Response{}	"Failed to get all offers"
func (c *offerHandler) GetAllOffers(ctx *gin.Context) {

	pagination := requests.GetPagination(ctx)

	offers, err := c.offerUseCase.FindAllOffers(ctx, pagination)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get all offers", err, nil)

		return
	}

	if offers == nil {
		responses.SuccessResponse(ctx, http.StatusOK, "No offer found", offers)

		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully found all offers", offers)
}

// RemoveOffer godoc
//
//	@summary		Remove offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to remove an offer
//	@Id				RemoveOffer
//	@Tags			Admin Offers
//	@Param			offer_id	path	int	true	"Offer ID"
//	@Router			/admin/offers/{offer_id} [delete]
//	@Success		200	{object}	responses.Response{}	"successfully offer added"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c *offerHandler) RemoveOffer(ctx *gin.Context) {

	offerID, err := requests.GetParamAsUint(ctx, "offer_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.offerUseCase.RemoveOffer(ctx, offerID)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer", err, nil)
		return
	}

	responses.SuccessResponse(ctx, 200, "successfully offer removed", nil)

}

// @Summary		Add category offer (Admin)
// @Security		BearerAuth
// @Description	API for admin to add an offer category
// @Id				SaveCategoryOffer
// @Tags			Admin Offers
// @Param			input	body	requests.OfferCategory{}	true	"input field"
// @Router			/admin/offers/category [post]
// @Success		200	{object}	responses.Response{}	"successfully offer added for category"
// @Failure		400	{object}	responses.Response{}	"invalid input"
func (c *offerHandler) SaveCategoryOffer(ctx *gin.Context) {

	var body requests.OfferCategory

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.offerUseCase.SaveCategoryOffer(ctx, body)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, usecases.ErrOfferAlreadyEnded):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrCategoryOfferAlreadyExist):
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to add offer", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully offer added for given category")
}

// GetAllCategoryOffers godoc
//
//	@Summary		Get all category offers (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all category offers
//	@Id				GetAllCategoryOffers
//	@Tags			Admin Offers
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/offers/category [get]
//	@Success		200	{object}	responses.Response{}	"successfully got all offer_category"
//	@Failure		500	{object}	responses.Response{}	"failed to get offers_category"
func (c *offerHandler) GetAllCategoryOffers(ctx *gin.Context) {

	pagination := requests.GetPagination(ctx)

	offerCategories, err := c.offerUseCase.FindAllCategoryOffers(ctx, pagination)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get offer categories", err, nil)
		return
	}

	if len(offerCategories) == 0 {
		responses.SuccessResponse(ctx, http.StatusOK, "No offer categories found", nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully found offers categories", offerCategories)
}

// RemoveCategoryOffer godoc
//
//	@Summary		Remove category offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to remove a offer from category
//	@Id				RemoveCategoryOffer
//	@Tags			Admin Offers
//	@Param			offer_category_id	path	int	true	"Offer Category ID"
//	@Router			/admin/offers/category/{offer_category_id} [delete]
//	@Success		200	{object}	responses.Response{}	"successfully offer added for category"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c *offerHandler) RemoveCategoryOffer(ctx *gin.Context) {

	offerCategoryID, err := requests.GetParamAsUint(ctx, "offer_category_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.offerUseCase.RemoveCategoryOffer(ctx, offerCategoryID)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer form category", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully offer removed from category")
}

// ChangeCategoryOffer godoc
//
//	@Summary		Change product offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to change product offer to another offer
//	@Id				ChangeCategoryOffer
//	@Tags			Admin Offers
//	@Param			input	body	requests.UpdateCategoryOffer{}	true	"input field"
//	@Router			/admin/offers/category [patch]
//	@Success		200	{object}	responses.Response{}	"successfully offer replaced for category"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c *offerHandler) ChangeCategoryOffer(ctx *gin.Context) {

	var body requests.UpdateCategoryOffer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.offerUseCase.ChangeCategoryOffer(ctx, body.CategoryOfferID, body.OfferID)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to change offer for given category offer", err, nil)
		return
	}

	responses.SuccessResponse(ctx, 200, "Successfully offer changed for given category offer")
}

// SaveProductOffer godoc
//
//	@Summary		Add product offer (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to add an offer for product
//	@Id				SaveProductOffer
//	@Tags			Admin Offers
//	@Param			input	body	requests.OfferProduct{}	true	"input field"
//	@Router			/admin/offers/products [post]
//	@Success		200	{object}	responses.Response{}	"successfully offer added for product"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c *offerHandler) SaveProductOffer(ctx *gin.Context) {

	var body requests.OfferProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offerProduct models.OfferProduct
	copier.Copy(&offerProduct, &body)

	err := c.offerUseCase.SaveProductOffer(ctx, offerProduct)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add offer for given product", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully offer added to given product")
}

// GetAllProductsOffers godoc
//
//	@Summary		Get all product offers (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all product offers
//	@Id				GetAllProductsOffers
//	@Tags			Admin Offers
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/offers/products [get]
//	@Success		200	{object}	responses.Response{}	"successfully got all offers_categories"
//	@Failure		500	{object}	responses.Response{}	"failed to get offer_products"
func (c *offerHandler) GetAllProductsOffers(ctx *gin.Context) {

	pagination := requests.GetPagination(ctx)

	offersOfCategories, err := c.offerUseCase.FindAllProductOffers(ctx, pagination)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get all offer products", err, nil)
		return
	}

	if offersOfCategories == nil {
		responses.SuccessResponse(ctx, http.StatusOK, "No offer products found", nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully found all offer products", offersOfCategories)
}

// RemoveProductOffer godoc
//
//	@Summary		Remove product offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to remove a offer from product
//	@Id				RemoveProductOffer
//	@Tags			Admin Offers
//	@param			offer_product_id	path	int	true	"offer_product_id"
//	@Router			/admin/offers/products/{offer_product_id} [delete]
//	@Success		200	{object}	responses.Response{}	"Successfully offer removed from product"
//	@Failure		400	{object}	responses.Response{}	"invalid input on params"
func (c *offerHandler) RemoveProductOffer(ctx *gin.Context) {

	offerProductID, err := requests.GetParamAsUint(ctx, "offer_product_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.offerUseCase.RemoveProductOffer(ctx, offerProductID)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer form product", err, nil)
		return
	}

	responses.SuccessResponse(ctx, 200, "Successfully offer removed from product")
}

// ChangeProductOffer godoc
//
//	@Summary		Change product offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to change product offer to another offer
//	@Id				ChangeProductOffer
//	@Tags			Admin Offers
//	@Param			input	body	requests.UpdateProductOffer{}	true	"input field"
//	@Router			/admin/offers/products [patch]
//	@Success		200	{object}	responses.Response{}	"Successfully offer changed for  given product offer"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
func (c *offerHandler) ChangeProductOffer(ctx *gin.Context) {

	var body requests.UpdateProductOffer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.offerUseCase.ChangeProductOffer(ctx, body.ProductOfferID, body.OfferID)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to change offer for given product offer", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully offer changed for  given product offer")
}
