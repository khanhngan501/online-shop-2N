package usecases

import (
	"context"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/repositories/interfaces"
	service "online-shop-2N/pkg/usecases/interfaces"
	"online-shop-2N/pkg/utils"
)

type categoryUseCase struct {
	categoryRepo interfaces.CategoryRepository
}

// to get a new instance of categoryUseCase
func NewCategoryUseCase(categoryRepo interfaces.CategoryRepository) service.CategoryUseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (c *categoryUseCase) FindAllCategories(ctx context.Context, pagination requests.Pagination) ([]responses.Category, error) {

	categories, err := c.categoryRepo.FindAllMainCategories(ctx, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed find all main categories")
	}

	for i, category := range categories {

		subCategory, err := c.categoryRepo.FindAllSubCategories(ctx, category.ID)
		if err != nil {
			return nil, utils.PrependMessageToError(err, "failed to find sub categories")
		}
		categories[i].SubCategory = subCategory
	}

	return categories, nil
}

// Save category
func (c *categoryUseCase) SaveCategory(ctx context.Context, categoryName string) error {

	categoryExist, err := c.categoryRepo.IsCategoryNameExist(ctx, categoryName)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check category already exist")
	}
	if categoryExist {
		return ErrCategoryAlreadyExist
	}

	err = c.categoryRepo.SaveCategory(ctx, categoryName)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save category")
	}

	return nil
}

// Save Sub category
func (c *categoryUseCase) SaveSubCategory(ctx context.Context, subCategory requests.SubCategory) error {

	subCatExist, err := c.categoryRepo.IsSubCategoryNameExist(ctx, subCategory.Name, subCategory.CategoryID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check sub category already exist")
	}
	if subCatExist {
		return ErrCategoryAlreadyExist
	}

	err = c.categoryRepo.SaveSubCategory(ctx, subCategory.CategoryID, subCategory.Name)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save sub category")
	}

	return nil
}
