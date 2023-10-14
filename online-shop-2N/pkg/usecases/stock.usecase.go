package usecases

import (
	"context"
	"log"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/repositories/interfaces"
	service "online-shop-2N/pkg/usecases/interfaces"
)

type stockUseCase struct {
	stockRepo interfaces.StockRepository
}

func NewStockUseCase(stockRepo interfaces.StockRepository) service.StockUseCase {

	return &stockUseCase{
		stockRepo: stockRepo,
	}
}

func (c *stockUseCase) GetAllStockDetails(ctx context.Context, pagination requests.Pagination) (stocks []responses.Stock, err error) {
	stocks, err = c.stockRepo.FindAll(ctx, pagination)

	if err != nil {
		return stocks, err
	}
	log.Printf("successfully got stock details")
	return stocks, nil
}

func (c *stockUseCase) UpdateStockBySKU(ctx context.Context, updateDetails requests.UpdateStock) error {

	err := c.stockRepo.Update(ctx, updateDetails)

	if err != nil {
		return err
	}

	log.Printf("successfully updated of stock details of stock with sku %v", updateDetails.SKU)
	return nil
}
