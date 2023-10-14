package database

import (
	"fmt"
	"online-shop-2N/pkg/config"
	"online-shop-2N/pkg/models"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB(config config.Config) (*gorm.DB, error) {
	url := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Error("Failed to open database connection. Due to error: ", err)
		return nil, err
	}
	// Migrate database models to database
	err = db.AutoMigrate(
		//auth
		models.RefreshSession{},
		models.OtpSession{},
		//user
		models.User{},
		models.Country{},
		models.Address{},
		models.UserAddress{},

		//admin
		models.Admin{},

		//product
		models.Category{},
		models.Product{},
		models.ProductItem{},
		models.ProductImage{},

		// wish list
		models.WishList{},

		// cart
		models.Cart{},
		models.CartItem{},

		// order
		models.OrderStatus{},
		models.ShopOrder{},
		models.OrderLine{},
		models.OrderReturn{},

		//offer
		models.Offer{},
		models.OfferCategory{},
		models.OfferProduct{},

		// coupon
		models.Coupon{},
		models.CouponUses{},

		//wallet
		models.Wallet{},
		models.Transaction{},
	)
	if err != nil {
		log.Error("Failed to migrate database tables. Due to error: ", err)
	}
	// setup the triggers
	if err := SetUpDBTriggers(db); err != nil {
		log.Error("Failed to setup database triggers. Due to error: ", err)
		return nil, err
	}
	if err := saveAdmin(db, config.AdminEmail, config.AdminUserName, config.AdminPassword); err != nil {
		return nil, err
	}

	if err := saveOrderStatuses(db); err != nil {
		return nil, err
	}
	if err := savePaymentMethods(db); err != nil {
		return nil, err
	}
	return db, nil
}
