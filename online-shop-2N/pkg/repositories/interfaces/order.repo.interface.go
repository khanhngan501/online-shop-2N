package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	commonConstant "online-shop-2N/pkg/common/constants"
	"online-shop-2N/pkg/models"
)

type OrderRepository interface {
	Transaction(callBack func(transactionRepo OrderRepository) error) error

	SaveOrderLine(ctx context.Context, orderLine models.OrderLine) error

	UpdateShopOrderOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error
	UpdateShopOrderStatusAndSavePaymentMethod(ctx context.Context, shopOrderID, orderStatusID, paymentID uint) error

	// shop order order
	SaveShopOrder(ctx context.Context, shopOrder models.ShopOrder) (shopOrderID uint, err error)
	FindShopOrderByShopOrderID(ctx context.Context, shopOrderID uint) (models.ShopOrder, error)
	FindAllShopOrders(ctx context.Context, pagination requests.Pagination) (shopOrders []responses.ShopOrder, err error)
	FindAllShopOrdersByUserID(ctx context.Context, userID uint, pagination requests.Pagination) ([]responses.ShopOrder, error)

	// find shop order items
	FindAllOrdersItemsByShopOrderID(ctx context.Context,
		shopOrderID uint, pagination requests.Pagination) (orderItems []responses.OrderItem, err error)

	// order status
	FindOrderStatusByShopOrderID(ctx context.Context, shopOrderID uint) (models.OrderStatus, error)
	FindOrderStatusByID(ctx context.Context, orderStatusID uint) (models.OrderStatus, error)
	FindOrderStatusByStatus(ctx context.Context, orderStatus commonConstant.OrderStatusType) (models.OrderStatus, error)
	FindAllOrderStatuses(ctx context.Context) ([]models.OrderStatus, error)

	//order return
	FindOrderReturnByReturnID(ctx context.Context, orderReturnID uint) (models.OrderReturn, error)
	FindOrderReturnByShopOrderID(ctx context.Context, shopOrderID uint) (orderReturn models.OrderReturn, err error)
	FindAllOrderReturns(ctx context.Context, pagination requests.Pagination) ([]responses.OrderReturn, error)
	FindAllPendingOrderReturns(ctx context.Context, pagination requests.Pagination) ([]responses.OrderReturn, error)
	SaveOrderReturn(ctx context.Context, orderReturn models.OrderReturn) error
	UpdateOrderReturn(ctx context.Context, orderReturn models.OrderReturn) error

	// wallet
	FindWalletByUserID(ctx context.Context, userID uint) (wallet models.Wallet, err error)
	SaveWallet(ctx context.Context, userID uint) (walletID uint, err error)
	UpdateWallet(ctx context.Context, walletID, updateTotalAmount uint) error
	SaveWalletTransaction(ctx context.Context, walletTrx models.Transaction) error

	FindWalletTransactions(ctx context.Context, walletID uint,
		pagination requests.Pagination) (transaction []models.Transaction, err error)
}
