package interfaces

import (
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/models"
)

type BrandUseCase interface {
	Save(brand models.Brand) (models.Brand, error)
	Update(brand models.Brand) error
	FindAll(pagination requests.Pagination) ([]models.Brand, error)
	FindOne(brandID uint) (models.Brand, error)
	Delete(brandID uint) error
}
