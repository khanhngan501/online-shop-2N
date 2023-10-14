package http

import (
	"net/http"
	handlerInterface "online-shop-2N/pkg/api/handlers/interfaces"
	"online-shop-2N/pkg/api/middlewares"
	"online-shop-2N/pkg/api/routes"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	Engine *gin.Engine
}

func NewServerHTTP(authHandler handlerInterface.AuthHandler, middlewares middlewares.Middleware,
	adminHandler handlerInterface.AdminHandler, userHandler handlerInterface.UserHandler,
	cartHandler handlerInterface.CartHandler, paymentHandler handlerInterface.PaymentHandler,
	productHandler handlerInterface.ProductHandler, categoryHandler handlerInterface.CategoryHandler,
	orderHandler handlerInterface.OrderHandler,
	couponHandler handlerInterface.CouponHandler, offerHandler handlerInterface.OfferHandler,
	stockHandler handlerInterface.StockHandler, branHandler handlerInterface.BrandHandler,
) *ServerHTTP {
	engine := gin.New()

	// engine.LoadHTMLGlob("views/*.html")

	engine.Use(gin.Logger())

	// Set up routers and handlers
	routes.UserRoutes(engine.Group("/api"), authHandler, middlewares, userHandler, cartHandler,
		productHandler, paymentHandler, orderHandler, couponHandler)
	routes.AdminRoutes(engine.Group("/api/admin"), authHandler, middlewares, adminHandler,
		productHandler, categoryHandler, paymentHandler, orderHandler, couponHandler, offerHandler, stockHandler, branHandler)

	// No hanldlers
	engine.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Invalid URL provided",
		})
	})
	return &ServerHTTP{Engine: engine}
}

func (s *ServerHTTP) Start() error {
	return s.Engine.Run(":8000")
}
