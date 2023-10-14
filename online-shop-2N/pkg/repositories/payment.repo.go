package repositories

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/models"
	"online-shop-2N/pkg/repositories/interfaces"

	"gorm.io/gorm"

	commonConstant "online-shop-2N/pkg/common/constants"
)

type paymentDatabase struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) interfaces.PaymentRepository {
	return &paymentDatabase{
		db: db,
	}
}

func (c *paymentDatabase) FindPaymentMethodByID(ctx context.Context, paymentMethodID uint) (paymentMethods models.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods WHERE id = $1`

	err = c.db.Raw(query, paymentMethodID).Scan(&paymentMethods).Error

	return paymentMethods, err
}

// find payment_method by name
func (c *paymentDatabase) FindPaymentMethodByType(ctx context.Context,
	paymentType commonConstant.PaymentType) (paymentMethod models.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods WHERE name = $1`
	err = c.db.Raw(query, paymentType).Scan(&paymentMethod).Error

	return paymentMethod, err
}

func (c *paymentDatabase) FindAllPaymentMethods(ctx context.Context) (paymentMethods []models.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods`
	err = c.db.Raw(query).Scan(&paymentMethods).Error

	return paymentMethods, err
}

func (c *paymentDatabase) SavePaymentMethod(ctx context.Context, paymentMethod models.PaymentMethod) (paymentMethodID uint, err error) {

	query := `INSERT INTO payment_methods (name, block_status, maximum_amount) VALUES ($1, $2, $3)`
	err = c.db.Raw(query, paymentMethod.Name, paymentMethod.BlockStatus, paymentMethod.MaximumAmount).Scan(&paymentMethod).Error

	return paymentMethod.ID, err
}
func (c *paymentDatabase) UpdatePaymentMethod(ctx context.Context, paymentMethod requests.PaymentMethodUpdate) error {

	query := `UPDATE payment_methods SET  block_status = $1, maximum_amount = $2 WHERE id = $3`

	err := c.db.Exec(query, paymentMethod.BlockStatus, paymentMethod.MaximumAmount, paymentMethod.ID).Error

	return err
}
