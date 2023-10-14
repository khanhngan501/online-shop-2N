//go:build wireinject
// +build wireinject

package di

import (
	http "online-shop-2N/pkg/api"
	"online-shop-2N/pkg/api/handlers"
	"online-shop-2N/pkg/api/middlewares"
	"online-shop-2N/pkg/config"
	"online-shop-2N/pkg/database"
	"online-shop-2N/pkg/repositories"
	"online-shop-2N/pkg/services/cloud"
	"online-shop-2N/pkg/services/otp"
	"online-shop-2N/pkg/services/tokens"
	"online-shop-2N/pkg/usecases"

	"github.com/google/wire"
)

func InitializeApi(cfg config.Config) (*http.ServerHTTP, error) {

	wire.Build(database.ConnectToDB,
		//external
		tokens.NewTokenService,
		otp.NewOtpAuth,
		cloud.NewAWSCloudService,

		// repositories

		middlewares.NewMiddleware,
		repositories.NewAuthRepository,
		repositories.NewPaymentRepository,
		repositories.NewAdminRepository,
		repositories.NewUserRepository,
		repositories.NewCartRepository,
		repositories.NewProductRepository,
		repositories.NewCategoryRepository,
		repositories.NewOrderRepository,
		repositories.NewCouponRepository,
		repositories.NewOfferRepository,
		repositories.NewStockRepository,
		repositories.NewBrandDatabaseRepository,

		//usecases
		usecases.NewAuthUseCase,
		usecases.NewAdminUseCase,
		usecases.NewUserUseCase,
		usecases.NewCartUseCase,
		usecases.NewPaymentUseCase,
		usecases.NewProductUseCase,
		usecases.NewCategoryUseCase,
		usecases.NewOrderUseCase,
		usecases.NewCouponUseCase,
		usecases.NewOfferUseCase,
		usecases.NewStockUseCase,
		usecases.NewBrandUseCase,
		// handlers
		handlers.NewAuthHandler,
		handlers.NewAdminHandler,
		handlers.NewUserHandler,
		handlers.NewCartHandler,
		handlers.NewPaymentHandler,
		handlers.NewProductHandler,
		handlers.NewCategoryHandler,
		handlers.NewOrderHandler,
		handlers.NewCouponHandler,
		handlers.NewOfferHandler,
		handlers.NewStockHandler,
		handlers.NewBrandHandler,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
