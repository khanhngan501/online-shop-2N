package repositories

import (
	"context"
	"errors"
	"fmt"
	"online-shop-2N/pkg/api/handlers/responses"
	"online-shop-2N/pkg/models"
	"online-shop-2N/pkg/repositories/interfaces"
	"time"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB: DB}
}

func (c *userDatabase) FindUserByUserID(ctx context.Context, userID uint) (user models.User, err error) {

	query := `SELECT * FROM users WHERE id = $1`
	err = c.DB.Raw(query, userID).Scan(&user).Error

	return user, err
}

func (c *userDatabase) FindUserByEmail(ctx context.Context, email string) (user models.User, err error) {

	query := `SELECT * FROM users WHERE email = $1`
	err = c.DB.Raw(query, email).Scan(&user).Error

	return user, err
}

func (c *userDatabase) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (user models.User, err error) {

	query := `SELECT * FROM users WHERE phone = $1`
	err = c.DB.Raw(query, phoneNumber).Scan(&user).Error

	return user, err
}
func (c *userDatabase) FindUserByUserName(ctx context.Context, userName string) (user models.User, err error) {

	query := `SELECT * FROM users WHERE user_name = $1`
	err = c.DB.Raw(query, userName).Scan(&user).Error

	return user, err
}

func (c *userDatabase) FindUserByUserNameEmailOrPhoneNotID(ctx context.Context,
	userDetails models.User) (user models.User, err error) {

	query := `SELECT * FROM users WHERE (user_name = $1 OR email = $2 OR phone = $3) AND id != $4`
	err = c.DB.Raw(query, userDetails.UserName, userDetails.Email, userDetails.Phone, userDetails.ID).Scan(&user).Error

	return
}

func (c *userDatabase) SaveUser(ctx context.Context, user models.User) (userID uint, err error) {

	//save the user details
	query := `INSERT INTO users (user_name, first_name, 
		last_name, age, email, phone, password, google_image, created_at, verified) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ) RETURNING id`

	createdAt := time.Now()
	err = c.DB.Raw(query, user.UserName, user.FirstName, user.LastName,
		user.Age, user.Email, user.Phone, user.Password, user.GoogleImage, createdAt, true).Scan(&userID).Error

	return userID, err
}

func (c *userDatabase) UpdateVerified(ctx context.Context, userID uint) error {

	query := `UPDATE users SET verified = 'T' WHERE id = $1`
	err := c.DB.Exec(query, userID).Error

	return err
}

func (c *userDatabase) UpdateUser(ctx context.Context, user models.User) (err error) {

	updatedAt := time.Now()
	// check password need to update or not
	if user.Password != "" {
		query := `UPDATE users SET user_name = $1, first_name = $2, last_name = $3,age = $4, 
		email = $5, phone = $6, password = $7, updated_at = $8 WHERE id = $9`
		err = c.DB.Exec(query, user.UserName, user.FirstName, user.LastName, user.Age, user.Email,
			user.Phone, user.Password, updatedAt, user.ID).Error
	} else {
		query := `UPDATE users SET user_name = $1, first_name = $2, last_name = $3,age = $4, 
		email = $5, phone = $6, updated_at = $7 WHERE id = $8`
		err = c.DB.Exec(query, user.UserName, user.FirstName, user.LastName, user.Age, user.Email,
			user.Phone, updatedAt, user.ID).Error
	}

	if err != nil {
		return fmt.Errorf("filed to update user detail of user with user_id %d", user.ID)
	}
	return nil
}

func (c *userDatabase) UpdateBlockStatus(ctx context.Context, userID uint, blockStatus bool) error {

	query := `UPDATE users SET block_status = $1 WHERE id = $2`
	err := c.DB.Exec(query, blockStatus, userID).Error

	return err
}

func (c *userDatabase) IsAddressIDExist(ctx context.Context, addressID uint) (exist bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM addresses WHERE id = $1) AS exist FROM addresses`
	err = c.DB.Raw(query, addressID).Scan(&exist).Error

	return
}
func (c *userDatabase) FindAddressByID(ctx context.Context, addressID uint) (address responses.Address, err error) {

	query := `SELECT adrs.id, adrs.detail_address, adrs.name, adrs.phone_number, adrs.commune, adrs.district, adrs.province,
	adrs.pincode, country_id, country_name FROM addresses adrs 
	INNER JOIN countries c ON c.id = adrs.country_id  
	INNER JOIN user_addresses uadrs ON uadrs.address_id = adrs.id 
	WHERE adrs.id = $1 `
	err = c.DB.Raw(query, addressID).Scan(&address).Error

	return
}

func (c *userDatabase) IsAddressAlreadyExistForUser(ctx context.Context, address models.Address, userID uint) (exist bool, err error) {
	query := `SELECT DISTINCT CASE  WHEN adrs.id != 0 THEN 'T' ELSE 'F' END AS exist 
	FROM addresses adrs 
	INNER JOIN user_addresses urs ON adrs.id = urs.address_id 
	WHERE adrs.name = $1 AND adrs.detail_address = $2 AND adrs.commune = $3 AND adrs.district = $4 AND adrs.province = $5 
	AND adrs.pincode = $6 AND adrs.country_id = $7  AND urs.user_id = $8`
	err = c.DB.Raw(query, address.Name, address.DetailAddress, address.Commune, address.District, address.Province, address.Pincode, address.CountryID, userID).Scan(&exist).Error

	return
}

func (c *userDatabase) FindAllAddressByUserID(ctx context.Context, userID uint) (addresses []responses.Address, err error) {

	query := `SELECT a.id, a.detail_address, a.name, a.phone_number, a.commune, a.district,a.province, 
	a.pincode, a.country_id, c.country_name, ua.is_default
	FROM user_addresses ua JOIN addresses a ON ua.address_id=a.id 
	INNER JOIN countries c ON a.country_id=c.id AND ua.user_id = $1`

	err = c.DB.Raw(query, userID).Scan(&addresses).Error

	return addresses, err
}

func (c *userDatabase) FindCountryByID(ctx context.Context, countryID uint) (models.Country, error) {

	var country models.Country

	if c.DB.Raw("SELECT * FROM countries WHERE id = ?", countryID).Scan(&country).Error != nil {
		return country, errors.New("filed to find the country")
	}

	return country, nil
}

// save address
func (c *userDatabase) SaveAddress(ctx context.Context, address models.Address) (addressID uint, err error) {
	query := `INSERT INTO addresses (name, phone_number, detail_address,commune, district, province, pincode, country_id, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	createdAt := time.Now()

	if c.DB.Raw(query, address.Name, address.PhoneNumber,
		address.DetailAddress, address.Commune, address.District, address.Province,
		address.Pincode, address.CountryID, createdAt,
	).Scan(&address).Error != nil {
		return addressID, errors.New("filed to insert address on database")
	}
	return address.ID, nil
}

// update address
func (c *userDatabase) UpdateAddress(ctx context.Context, address models.Address) error {

	address.CountryID = 1 // hardcoded !!!! should change
	query := `UPDATE addresses SET name=$1, phone_number=$2, detail_address=$3, commune=$4, district=$5, 
	province=$6, pincode=$7,country_id=$8, updated_at = $9 WHERE id=$10`

	updatedAt := time.Now()
	if c.DB.Raw(query, address.Name, address.PhoneNumber, address.DetailAddress,
		address.Commune, address.District, address.Province, address.Pincode,
		address.CountryID, updatedAt, address.ID).Scan(&address).Error != nil {
		return errors.New("filed to update the address for edit address")
	}
	return nil
}

func (c *userDatabase) SaveUserAddress(ctx context.Context, userAddress models.UserAddress) error {
	// first check user's first address is this or not
	var userID uint
	query := `SELECT address_id FROM user_addresses WHERE user_id = $1`
	err := c.DB.Raw(query, userAddress.UserID).Scan(&userID).Error
	if err != nil {
		return fmt.Errorf("failed to check user have already address exit or not with user_id %v", userAddress.UserID)
	}

	// if the given address is need to set default  then remove all other from default
	if userID == 0 { // it means user have no other addresses
		userAddress.IsDefault = true
	} else if userAddress.IsDefault {
		query := `UPDATE user_addresses SET is_default = 'f' WHERE user_id = ?`
		if c.DB.Raw(query, userAddress.UserID).Scan(&userAddress).Error != nil {
			return errors.New("filed to remove default status of address")
		}
	}

	query = `INSERT INTO user_addresses (user_id,address_id,is_default) VALUES ($1, $2, $3)`
	err = c.DB.Exec(query, userAddress.UserID, userAddress.AddressID, userAddress.IsDefault).Error
	if err != nil {
		return errors.New("filed to insert userAddress on database")
	}
	return nil
}

func (c *userDatabase) UpdateUserAddress(ctx context.Context, userAddress models.UserAddress) error {
	// if it need to set default the change the old default
	if userAddress.IsDefault {
		query := `UPDATE user_addresses SET is_default = 'f' WHERE user_id = ?`
		if c.DB.Raw(query, userAddress.UserID).Scan(&userAddress).Error != nil {
			return errors.New("filed to remove default status of address")
		}
	}

	// update the user address
	query := `UPDATE user_addresses SET is_default = ? WHERE address_id=? AND user_id=?`
	if c.DB.Raw(query, userAddress.IsDefault, userAddress.AddressID, userAddress.UserID).Scan(&userAddress).Error != nil {
		return errors.New("filed to update user address")
	}
	return nil
}

// wish list

func (c *userDatabase) FindWishListItem(ctx context.Context, productID, userID uint) (models.WishList, error) {

	var wishList models.WishList
	query := `SELECT * FROM wish_lists WHERE user_id=? AND product_item_id=?`
	if c.DB.Raw(query, userID, productID).Scan(&wishList).Error != nil {
		return wishList, errors.New("filed to find wishlist item")
	}
	return wishList, nil
}

func (c *userDatabase) FindAllWishListItemsByUserID(ctx context.Context, userID uint) (productItems []responses.WishListItem, err error) {

	query := `SELECT p.name, wl.id, pi.id AS product_item_id, pi.product_id, pi.price, pi.discount_price, 
	pi.qty_in_stock, sku FROM wish_lists wl 
	INNER JOIN product_items pi ON wl.product_item_id = pi.id 
	INNER JOIN products p ON pi.product_id = p.id 
	AND wl.user_id = $1`
	err = c.DB.Raw(query, userID).Scan(&productItems).Error

	return
}

func (c *userDatabase) SaveWishListItem(ctx context.Context, wishList models.WishList) error {

	query := `INSERT INTO wish_lists (user_id,product_item_id) VALUES ($1,$2) RETURNING *`

	if c.DB.Raw(query, wishList.UserID, wishList.ProductItemID).Scan(&wishList).Error != nil {
		return errors.New("filed to insert new wishlist on database")
	}
	return nil
}

func (c *userDatabase) RemoveWishListItem(ctx context.Context, userID, productItemID uint) error {

	query := `DELETE FROM wish_lists WHERE product_item_id = $1 AND user_id = $2`
	err := c.DB.Exec(query, productItemID, userID).Error

	return err
}
