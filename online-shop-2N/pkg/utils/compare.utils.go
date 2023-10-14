package utils

import (
	"errors"
	"online-shop-2N/pkg/models"
)

// To compare and return error for the same fields
func CompareUserExistingDetails(user1, user2 models.User) error {
	var err error
	if user1.Email == user2.Email {
		err = AppendMessageToError(err, "user already exist with this email")
	}
	if user1.UserName == user2.UserName {
		err = AppendMessageToError(err, "user already exist with this user name")
	}
	if user1.Phone == user2.Phone {
		err = AppendMessageToError(err, "user already exist with this phone")
	}

	if err == nil {
		return errors.New("failed to find existing details")
	}

	return err
}
