package repositories

import (
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/models"
	"online-shop-2N/pkg/repositories/interfaces"

	"gorm.io/gorm"
)

type brandDatabase struct {
	DB *gorm.DB
}

func NewBrandDatabaseRepository(db *gorm.DB) interfaces.BrandRepository {
	return &brandDatabase{
		DB: db,
	}
}

func (c *brandDatabase) IsExist(brand models.Brand) (bool, error) {

	res := c.DB.Where("name = ?", brand.Name).Find(&brand)
	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected != 0, nil
}

func (c *brandDatabase) Save(brand models.Brand) (models.Brand, error) {

	err := c.DB.Create(&brand).Error

	return brand, err
}

func (c *brandDatabase) Update(brand models.Brand) error {

	return c.DB.Where("id = ?", brand.ID).Updates(&brand).Error
}

func (c *brandDatabase) FindAll(pagination requests.Pagination) (brands []models.Brand, err error) {

	err = c.DB.Limit(int(pagination.Count)).Offset(int(pagination.PageNumber) - 1).Find(&brands).Error

	return
}

func (c *brandDatabase) FindOne(brandID uint) (brand models.Brand, err error) {

	err = c.DB.Where("id = ?", brandID).First(&brand).Error

	return
}

func (c *brandDatabase) Delete(brandID uint) error {

	return c.DB.Where("id = ?", brandID).Delete(&models.Brand{}).Error
}
