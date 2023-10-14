package interfaces

import (
	"context"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
)

type UserRepository interface {
	FindUserByUserID(ctx context.Context, userID uint) (user models.User, err error)
	FindUserByEmail(ctx context.Context, email string) (user models.User, err error)
	FindUserByUserName(ctx context.Context, userName string) (user models.User, err error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (user models.User, err error)
	FindUserByUserNameEmailOrPhoneNotID(ctx context.Context, user models.User) (models.User, error)

	SaveUser(ctx context.Context, user models.User) (userID uint, err error)
	UpdateVerified(ctx context.Context, userID uint) error
	UpdateUser(ctx context.Context, user models.User) (err error)
	UpdateBlockStatus(ctx context.Context, userID uint, blockStatus bool) error

	//address
	// FindCountryByID(ctx context.Context, countryID uint) (models.Country, error)
	FindAddressByID(ctx context.Context, addressID uint) (responses.Address, error)
	IsAddressIDExist(ctx context.Context, addressID uint) (exist bool, err error)
	IsAddressAlreadyExistForUser(ctx context.Context, address models.Address, userID uint) (bool, error)
	FindAllAddressByUserID(ctx context.Context, userID uint) ([]responses.Address, error)
	SaveAddress(ctx context.Context, address models.Address) (addressID uint, err error)
	UpdateAddress(ctx context.Context, address models.Address) error
	// address join table
	SaveUserAddress(ctx context.Context, userAdress models.UserAddress) error
	UpdateUserAddress(ctx context.Context, userAddress models.UserAddress) error

	//wishlist
	FindWishListItem(ctx context.Context, productID, userID uint) (models.WishList, error)
	FindAllWishListItemsByUserID(ctx context.Context, userID uint) ([]responses.WishListItem, error)
	SaveWishListItem(ctx context.Context, wishList models.WishList) error
	RemoveWishListItem(ctx context.Context, userID, productItemID uint) error
}
