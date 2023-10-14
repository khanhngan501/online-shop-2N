package handlers

import (
	"net/http"
	"online-shop-2N/pkg/api/handlers/interfaces"

	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	usecaseInterface "online-shop-2N/pkg/usecases/interfaces"

	"github.com/gin-gonic/gin"
)

type stockHandler struct {
	stockUseCase usecaseInterface.StockUseCase
}

func NewStockHandler(stockUseCase usecaseInterface.StockUseCase) interfaces.StockHandler {
	return &stockHandler{
		stockUseCase: stockUseCase,
	}
}

// GetAllStocks godoc
//
//	@Summary		Get all stocks (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all stocks
//	@Id				GetAllStocks
//	@Tags			Admin Stock
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/stocks [get]
//	@Success		200	{object}	responses.Response{}	"Successfully found all stocks"
//	@Success		204	{object}	responses.Response{}	"No stocks found"
//	@Failure		500	{object}	responses.Response{}	"Failed to Get all stocks"
func (c *stockHandler) GetAllStocks(ctx *gin.Context) {

	pagination := requests.GetPagination(ctx)

	stocks, err := c.stockUseCase.GetAllStockDetails(ctx, pagination)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to Get all stocks", err, nil)
		return
	}

	if len(stocks) == 0 {
		responses.SuccessResponse(ctx, http.StatusNoContent, "No stocks found", nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully found all stocks", stocks)
}

// UpdateStock godoc
//
//	@Summary		Update stocks (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to update stock details
//	@Id				UpdateStock
//	@Tags			Admin Stock
//	@Param			input	body	requests.UpdateStock{}	true	"Update stock details"
//	@Router			/admin/stocks [patch]
//	@Success		200	{object}	responses.Response{}	"Successfully updated sock"
//	@Failure		400	{object}	responses.Response{}	"Failed to bind input"
//	@Failure		500	{object}	responses.Response{}	"Failed to update stock"
func (c *stockHandler) UpdateStock(ctx *gin.Context) {

	var body requests.UpdateStock

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		responses.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	err = c.stockUseCase.UpdateStockBySKU(ctx, body)

	if err != nil {
		responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update stock", err, nil)
		return
	}

	responses.SuccessResponse(ctx, http.StatusOK, "Successfully updated sock", nil)
}
