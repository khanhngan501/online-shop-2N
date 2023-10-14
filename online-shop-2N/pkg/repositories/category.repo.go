package repositories

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/repositories/interfaces"

	"gorm.io/gorm"
)

type categoryDatabase struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepository {
	return &categoryDatabase{
		DB: db,
	}
}

func (c *categoryDatabase) Transactions(ctx context.Context, trxFn func(repo interfaces.CategoryRepository) error) error {

	trx := c.DB.Begin()

	repo := NewCategoryRepository(trx)

	if err := trxFn(repo); err != nil {
		trx.Rollback()
		return err
	}

	if err := trx.Commit().Error; err != nil {
		trx.Rollback()
		return err
	}
	return nil
}

// To check the category name exist
func (c *categoryDatabase) IsCategoryNameExist(ctx context.Context, name string) (exist bool, err error) {

	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1 AND category_id IS NULL)`
	err = c.DB.Raw(query, name).Scan(&exist).Error

	return
}

// Save Category
func (c *categoryDatabase) SaveCategory(ctx context.Context, categoryName string) (err error) {

	query := `INSERT INTO categories (name) VALUES ($1)`
	err = c.DB.Exec(query, categoryName).Error

	return err
}

// To check the sub category name already exist for the category
func (c *categoryDatabase) IsSubCategoryNameExist(ctx context.Context, name string, categoryID uint) (exist bool, err error) {

	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1 AND category_id = $2)`
	err = c.DB.Raw(query, name, categoryID).Scan(&exist).Error

	return
}

// Save Category as sub category
func (c *categoryDatabase) SaveSubCategory(ctx context.Context, categoryID uint, categoryName string) (err error) {

	query := `INSERT INTO categories (category_id, name) VALUES ($1, $2)`
	err = c.DB.Exec(query, categoryID, categoryName).Error

	return err
}

// Find all main category(its not have a category_id)
func (c *categoryDatabase) FindAllMainCategories(ctx context.Context,
	pagination requests.Pagination) (categories []responses.Category, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT id, name FROM categories WHERE category_id IS NULL 
	LIMIT $1 OFFSET $2`
	err = c.DB.Raw(query, limit, offset).Scan(&categories).Error

	return
}

// Find all sub categories of a category
func (c *categoryDatabase) FindAllSubCategories(ctx context.Context,
	categoryID uint) (subCategories []responses.SubCategory, err error) {

	query := `SELECT id, name FROM categories WHERE category_id = $1`
	err = c.DB.Raw(query, categoryID).Scan(&subCategories).Error

	return
}
