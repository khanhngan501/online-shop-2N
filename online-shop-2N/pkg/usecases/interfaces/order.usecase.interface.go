package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type OrderUseCase interface {

	//
	SaveOrder(ctx context.Context, userID, addressID uint) (shopOrderID uint, err error)

	// Find order and order items
	FindAllShopOrders(ctx context.Context, pagination requests.Pagination) (shopOrders []responses.ShopOrder, err error)
	FindUserShopOrder(ctx context.Context, userID uint, pagination requests.Pagination) ([]responses.ShopOrder, error)
	FindOrderItems(ctx context.Context, shopOrderID uint, pagination requests.Pagination) ([]responses.OrderItem, error)

	// cancel order and change order status
	FindAllOrderStatuses(ctx context.Context) (orderStatuses []models.OrderStatus, err error)
	UpdateOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error
	CancelOrder(ctx context.Context, shopOrderID uint) error

	// return and update
	SubmitReturnRequest(ctx context.Context, returnDetails requests.Return) error
	FindAllPendingOrderReturns(ctx context.Context, pagination requests.Pagination) ([]responses.OrderReturn, error)
	FindAllOrderReturns(ctx context.Context, pagination requests.Pagination) ([]responses.OrderReturn, error)
	UpdateReturnDetails(ctx context.Context, updateDetails requests.UpdateOrderReturn) error

	// wallet
	FindUserWallet(ctx context.Context, userID uint) (wallet models.Wallet, err error)
	FindUserWalletTransactions(ctx context.Context, userID uint, pagination requests.Pagination) (transactions []models.Transaction, err error)
}
