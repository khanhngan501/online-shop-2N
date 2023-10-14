package database

import (
	"fmt"
	"time"

	commonConstant "online-shop-2N/pkg/common/constants"
	"online-shop-2N/pkg/models"
	"online-shop-2N/pkg/utils"

	"gorm.io/gorm"
)

// To save predefined order statuses on database if its not exist
func saveOrderStatuses(db *gorm.DB) error {

	statuses := []commonConstant.OrderStatusType{
		commonConstant.StatusPaymentPending,
		commonConstant.StatusOrderPlaced,
		commonConstant.StatusOrderCancelled,
		commonConstant.StatusOrderDelivered,
		commonConstant.StatusReturnRequested,
		commonConstant.StatusReturnApproved,
		commonConstant.StatusReturnCancelled,
		commonConstant.StatusOrderReturned,
	}

	var (
		searchQuery = `SELECT CASE WHEN id != 0 THEN 'T' ELSE 'F' END as exist 
		FROM order_statuses WHERE status = $1`
		insertQuery = `INSERT INTO order_statuses (status) VALUES ($1)`
		exist       bool
		err         error
	)

	for _, status := range statuses {

		err = db.Raw(searchQuery, status).Scan(&exist).Error
		if err != nil {
			return fmt.Errorf("Failed to check order status already exist err: %w", err)
		}

		if !exist {
			err = db.Exec(insertQuery, status).Error
			if err != nil {
				return fmt.Errorf("Failed to save status %w", err)
			}
		}
		exist = false
	}
	return nil
}

// To save predefined payment methods on database if its not exist
func savePaymentMethods(db *gorm.DB) error {
	paymentMethods := []models.PaymentMethod{
		{
			Name:          commonConstant.CodPayment,
			MaximumAmount: commonConstant.CodMaximumAmount,
		},
		{
			Name:          commonConstant.RazopayPayment,
			MaximumAmount: commonConstant.RazorPayMaximumAmount,
		},
		{
			Name:          commonConstant.StripePayment,
			MaximumAmount: commonConstant.StripeMaximumAmount,
		},
	}

	var (
		searchQuery = `SELECT CASE WHEN id != 0 THEN 'T' ELSE 'F' END as exist FROM payment_methods WHERE name = $1`
		insertQuery = `INSERT INTO payment_methods (name, maximum_amount) VALUES ($1, $2)`
		exist       bool
		err         error
	)

	for _, paymentMethod := range paymentMethods {

		err = db.Raw(searchQuery, paymentMethod.Name).Scan(&exist).Error
		if err != nil {
			return fmt.Errorf("failed to check payment methods already exist %w", err)
		}
		if !exist {
			err = db.Exec(insertQuery, paymentMethod.Name, paymentMethod.MaximumAmount).Error
			if err != nil {
				return fmt.Errorf("failed to save payment method %w", err)
			}
		}
		exist = false
	}
	return nil
}

func saveAdmin(db *gorm.DB, email, userName, password string) error {

	var (
		searchQuery = `SELECT CASE WHEN id != 0 THEN 'T' ELSE 'F' END as exist FROM admins WHERE email = $1`
		insertQuery = `INSERT INTO admins (email, user_name, password, created_at) VALUES ($1, $2, $3, $4)`
		exist       bool
		err         error
	)

	err = db.Raw(searchQuery, email).Scan(&exist).Error
	if err != nil {
		return fmt.Errorf("failed to check admin already exist err:%w", err)
	}

	if !exist {
		hashPass, err := utils.GetHashedPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password err: %w", err)
		}
		createdAt := time.Now()
		err = db.Exec(insertQuery, email, userName, hashPass, createdAt).Error
		if err != nil {
			return fmt.Errorf("failed to save admin details %w", err)
		}
	}
	return nil
}
