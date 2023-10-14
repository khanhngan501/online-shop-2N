package handlers

import (
	"errors"
	"fmt"
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

type ProductHandler struct {
	productUseCase usecaseInterface.ProductUseCase
}

func NewProductHandler(productUsecase usecaseInterface.ProductUseCase) interfaces.ProductHandler {
	return &ProductHandler{
		productUseCase: productUsecase,
	}
}

// SaveVariation godoc
//
//	@Summary		Add new variations (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to add new variations for a category
//	@Tags			Admin Category
//	@ID				SaveVariation
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path	int					true	"Category ID"
//	@Param			input		body	requests.Variation{}	true	"Variation details"
//	@Router			/admin/categories/{category_id}/variations [post]
//	@Success		201	{object}	responses.Response{}	"Successfully added variations"
//	@Failure		400	{object}	responses.Response{}	"Invalid input"
//	@Failure		500	{object}	responses.Response{}	"Failed to add variation"
func (p *ProductHandler) SaveVariation(ctx *gin.Context) {

	categoryID, err := requests.GetParamAsUint(ctx, "category_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	var body requests.Variation

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err = p.productUseCase.SaveVariation(ctx, categoryID, body.Names)

	if err != nil {
		var statusCode = http.StatusInternalServerError
		if errors.Is(err, usecases.ErrVariationAlreadyExist) {
			statusCode = http.StatusConflict
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to add variation", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusCreated, "Successfully added variations")
}

// SaveVariationOption godoc
//
//	@Summary		Add new variation options (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to add variation options for a variation
//	@Tags			Admin Category
//	@ID				SaveVariationOption
//	@Accept			json
//	@Produce		json
//	@Param			category_id		path	int							true	"Category ID"
//	@Param			variation_id	path	int							true	"Variation ID"
//	@Param			input			body	requests.VariationOption{}	true	"Variation option details"
//	@Router			/admin/categories/{category_id}/variations/{variation_id}/options [post]
//	@Success		201	{object}	responses.Response{}	"Successfully added variation options"
//	@Failure		400	{object}	responses.Response{}	"Invalid input"
//	@Failure		500	{object}	responses.Response{}	"Failed to add variation options"
func (p *ProductHandler) SaveVariationOption(ctx *gin.Context) {

	variationID, err := requests.GetParamAsUint(ctx, "variation_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	var body requests.VariationOption

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err = p.productUseCase.SaveVariationOption(ctx, variationID, body.Values)
	if err != nil {
		var statusCode = http.StatusInternalServerError
		if errors.Is(err, usecases.ErrVariationOptionAlreadyExist) {
			statusCode = http.StatusConflict
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to add variation options", err, nil)
		return
	}
	responses.SuccessResponse(ctx, http.StatusCreated, "Successfully added variation options")
}

// GetAllVariations godoc
//
//	@Summary		Get all variations (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all variation and its values of a category
//	@Tags			Admin Category
//	@ID				GetAllVariations
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path	int	true	"Category ID"
//	@Router			/admin/categories/{category_id}/variations [get]
//	@Success		200	{object}	responses.Response{}	"Successfully retrieved all variations and its values"
//	@Failure		400	{object}	responses.Response{}	"Invalid input"
//	@Failure		500	{object}	responses.Response{}	"Failed to Get variations and its values"
func (c *ProductHandler) GetAllVariations(ctx *gin.Context) {

	categoryID, err := requests.GetParamAsUint(ctx, "category_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	variations, err := c.productUseCase.FindAllVariationsAndItsValues(ctx, categoryID)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to Get variations and its values", err, nil)
		return
	}

	if len(variations) == 0 {
		responses.SuccessResponse(ctx, http.StatusOK, "No variations found")
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all variations and its values", variations)
}

// SaveProduct godoc
//
//	@Summary		Add a new product (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to add a new product
//	@ID				SaveProduct
//	@Tags			Admin Products
//	@Produce		json
//	@Param			name		formData	string				true	"Product Name"
//	@Param			description	formData	string				true	"Product Description"
//	@Param			category_id	formData	int					true	"Category Id"
//	@Param			brand_id	formData	int					true	"Brand Id"
//	@Param			price		formData	int					true	"Product Price"
//	@Param			image		formData	file				true	"Product Description"
//	@Success		200			{object}	responses.Response{}	"successfully product added"
//	@Router			/admin/products [post]
//	@Failure		400	{object}	responses.Response{}	"invalid input"
//	@Failure		409	{object}	responses.Response{}	"Product name already exist"
func (p *ProductHandler) SaveProduct(ctx *gin.Context) {

	name, err1 := requests.GetFormValuesAsString(ctx, "name")
	description, err2 := requests.GetFormValuesAsString(ctx, "description")
	categoryID, err3 := requests.GetFormValuesAsUint(ctx, "category_id")
	price, err4 := requests.GetFormValuesAsUint(ctx, "price")
	brandID, err5 := requests.GetFormValuesAsUint(ctx, "brand_id")

	fileHeader, err6 := ctx.FormFile("image")

	err := errors.Join(err1, err2, err3, err4, err5, err6)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	product := requests.Product{
		Name:            name,
		Description:     description,
		CategoryID:      categoryID,
		BrandID:         brandID,
		Price:           price,
		ImageFileHeader: fileHeader,
	}

	err = p.productUseCase.SaveProduct(ctx, product)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecases.ErrProductAlreadyExist) {
			statusCode = http.StatusConflict
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to add product", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusCreated, "Successfully product added")
}

// GetAllProductsAdmin godoc
//
//	@Summary		Get all products (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all products
//	@ID				GetAllProductsAdmin
//	@Tags			Admin Products
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/products [get]
//	@Success		200	{object}	responses.Response{}	"Successfully found all products"
//	@Failure		500	{object}	responses.Response{}	"Failed to Get all products"
func (p *ProductHandler) GetAllProductsAdmin() func(ctx *gin.Context) {
	return p.getAllProducts()
}

// GetAllProductsUser godoc
//
//	@Summary		Get all products (User)
//	@Security		BearerAuth
//	@Description	API for user to get all products
//	@ID				GetAllProductsUser
//	@Tags			User Products
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/products [get]
//	@Success		200	{object}	responses.Response{}	"Successfully found all products"
//	@Failure		500	{object}	responses.Response{}	"Failed to get all products"
func (p *ProductHandler) GetAllProductsUser() func(ctx *gin.Context) {
	return p.getAllProducts()
}

// Get products is common for user and admin so this function is to get the common Get all products func for them
func (p *ProductHandler) getAllProducts() func(ctx *gin.Context) {

	return func(ctx *gin.Context) {

		pagination := requests.GetPagination(ctx)

		products, err := p.productUseCase.FindAllProducts(ctx, pagination)

		if err != nil {
			responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to Get all products", err, nil)
			return
		}

		if len(products) == 0 {
			responses.SuccessResponse(ctx, http.StatusOK, "No products found", nil)
			return
		}

		responses.SuccessResponse(ctx, http.StatusOK, "Successfully found all products", products)
	}

}

// UpdateProduct godoc
//
//	@Summary		Update a product (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to update a product
//	@ID				UpdateProduct
//	@Tags			Admin Products
//	@Accept			json
//	@Produce		json
//	@Param			input	body	requests.UpdateProduct{}	true	"Product update input"
//	@Router			/admin/products [put]
//	@Success		200	{object}	responses.Response{}	"successfully product updated"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
//	@Failure		409	{object}	responses.Response{}	"Failed to update product"
//	@Failure		500	{object}	responses.Response{}	"Product name already exist for another product"
func (c *ProductHandler) UpdateProduct(ctx *gin.Context) {

	var body requests.UpdateProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var product models.Product
	copier.Copy(&product, &body)

	err := c.productUseCase.UpdateProduct(ctx, product)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecases.ErrProductAlreadyExist) {
			statusCode = http.StatusConflict
		}
		responses.ErrorResponse(ctx, statusCode, "Failed to update product", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully product updated", nil)
}

// SaveProductItem godoc
//
//	@Summary		Add a product item (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to add a product item for a specific product(should select at least one variation option from each variations)
//	@ID				SaveProductItem
//	@Tags			Admin Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id				path		int		true	"Product ID"
//	@Param			price					formData	int		true	"Price"
//	@Param			qty_in_stock			formData	int		true	"Quantity In Stock"
//	@Param			variation_option_ids	formData	[]int	true	"Variation Option IDs"
//	@Param			images					formData	file	true	"Images"
//	@Router			/admin/products/{product_id}/items [post]
//	@Success		200	{object}	responses.Response{}	"Successfully product item added"
//	@Failure		400	{object}	responses.Response{}	"invalid input"
//	@Failure		409	{object}	responses.Response{}	"Product have already this configured product items exist"
func (p *ProductHandler) SaveProductItem(ctx *gin.Context) {

	productID, err := requests.GetParamAsUint(ctx, "product_id")
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
	}

	price, err1 := requests.GetFormValuesAsUint(ctx, "price")
	qtyInStock, err2 := requests.GetFormValuesAsUint(ctx, "qty_in_stock")
	variationOptionIDS, err3 := requests.GetArrayFormValueAsUint(ctx, "variation_option_ids")
	imageFileHeaders, err4 := requests.GetArrayOfFromFiles(ctx, "images")

	err = errors.Join(err1, err2, err3, err4)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	productItem := requests.ProductItem{
		Price:              price,
		VariationOptionIDs: variationOptionIDS,
		QtyInStock:         qtyInStock,
		ImageFileHeaders:   imageFileHeaders,
	}

	fmt.Println(productItem, productID)

	err = p.productUseCase.SaveProductItem(ctx, productID, productItem)

	if err != nil {

		var statusCode int

		switch {
		case errors.Is(err, usecases.ErrProductItemAlreadyExist):
			statusCode = http.StatusConflict
		case errors.Is(err, usecases.ErrNotEnoughVariations):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}

		responses.ErrorResponse(ctx, statusCode, "Failed to add product item", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusCreated, "Successfully product item added", nil)
}

// GetAllProductItemsAdmin godoc
//
//	@Summary		Get all product items (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all product items for a specific product
//	@ID				GetAllProductItemsAdmin
//	@Tags			Admin Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path	int	true	"Product ID"
//	@Router			/admin/products/{product_id}/items [get]
//	@Success		200	{object}	responses.Response{}	"Successfully get all product items"
//	@Failure		400	{object}	responses.Response{}	"Invalid input"
//	@Failure		400	{object}	responses.Response{}	"Failed to get all product items"
func (p *ProductHandler) GetAllProductItemsAdmin() func(ctx *gin.Context) {
	return p.getAllProductItems()
}

// GetAllProductItemsUser godoc
//
//	@Summary		Get all product items (User)
//	@Security		BearerAuth
//	@Description	API for user to get all product items for a specific product
//	@ID				GetAllProductItemsUser
//	@Tags			User Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path	int	true	"Product ID"
//	@Router			/products/{product_id}/items [get]
//	@Success		200	{object}	responses.Response{}	"Successfully get all product items"
//	@Failure		400	{object}	responses.Response{}	"Invalid input"
//	@Failure		400	{object}	responses.Response{}	"Failed to get all product items"
func (p *ProductHandler) GetAllProductItemsUser() func(ctx *gin.Context) {
	return p.getAllProductItems()
}

// same functionality of get all product items for admin and user
func (p *ProductHandler) getAllProductItems() func(ctx *gin.Context) {

	return func(ctx *gin.Context) {

		productID, err := requests.GetParamAsUint(ctx, "product_id")
		if err != nil {
			responses.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		}

		productItems, err := p.productUseCase.FindAllProductItems(ctx, productID)

		if err != nil {
			responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get all product items", err, nil)
			return
		}

		// check the product have productItem exist or not
		if len(productItems) == 0 {
			responses.SuccessResponse(ctx, http.StatusOK, "No product items found")
			return
		}

		responses.SuccessResponse(ctx, http.StatusOK, "Successfully get all product items ", productItems)
	}
}
