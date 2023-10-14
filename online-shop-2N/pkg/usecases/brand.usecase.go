package usecases

import (
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/models"
	repoInterface "online-shop-2N/pkg/repositories/interfaces"
	"online-shop-2N/pkg/usecases/interfaces"
	"online-shop-2N/pkg/utils"
)

type brandUseCase struct {
	brandRepo repoInterface.BrandRepository
}

func NewBrandUseCase(brandRepo repoInterface.BrandRepository) interfaces.BrandUseCase {
	return &brandUseCase{
		brandRepo: brandRepo,
	}
}

func (b *brandUseCase) Save(brand models.Brand) (models.Brand, error) {

	alreadyExist, err := b.brandRepo.IsExist(brand)
	if err != nil {
		return models.Brand{}, utils.PrependMessageToError(err, "failed to check brand name already exist")
	}

	if alreadyExist {
		return models.Brand{}, ErrBrandAlreadyExist
	}

	brand, err = b.brandRepo.Save(brand)
	if err != nil {
		return models.Brand{}, utils.PrependMessageToError(err, "failed to save brand on db")
	}

	return brand, nil
}

func (b *brandUseCase) Update(brand models.Brand) error {

	err := b.brandRepo.Update(brand)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to update brand on db")
	}

	return nil
}

func (b *brandUseCase) FindAll(pagination requests.Pagination) ([]models.Brand, error) {

	brands, err := b.brandRepo.FindAll(pagination)

	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find all brands from db")
	}

	return brands, nil
}

func (b *brandUseCase) FindOne(brandID uint) (models.Brand, error) {

	brand, err := b.brandRepo.FindOne(brandID)
	if err != nil {
		return models.Brand{}, utils.PrependMessageToError(err, "failed to find brand from db")
	}

	return brand, nil
}

func (b *brandUseCase) Delete(brandID uint) error {

	err := b.brandRepo.Delete(brandID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to delete brands from db")
	}

	return nil
}
