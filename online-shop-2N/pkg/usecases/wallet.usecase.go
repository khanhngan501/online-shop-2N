package usecases

import (
	"context"
	"fmt"
	"log"
	"online-shop-2N/pkg/api/handlers/requests"
	"online-shop-2N/pkg/models"
)

// FindUserWallet implements interfaces.OrderUseCase.
func (c *OrderUseCase) FindUserWallet(ctx context.Context, userID uint) (wallet models.Wallet, err error) {
	// first find the user wallet
	wallet, err = c.orderRepo.FindWalletByUserID(ctx, userID)
	if err != nil {
		return wallet, err
	} else if wallet.ID == 0 { // if user have no wallet then create a wallet for user
		wallet.ID, err = c.orderRepo.SaveWallet(ctx, userID)
		if err != nil {
			return wallet, err
		}
	}

	log.Printf("successfully got user wallet with wallet_id %v for user user_id %v", wallet.ID, userID)
	return wallet, nil
}

// FindUserWalletTransactions implements interfaces.OrderUseCase.
func (c *OrderUseCase) FindUserWalletTransactions(ctx context.Context, userID uint, pagination requests.Pagination) (transactions []models.Transaction, err error) {
	// first find the user wallet
	wallet, err := c.orderRepo.FindWalletByUserID(ctx, userID)
	if err != nil {
		return transactions, err
	} else if wallet.ID == 0 {
		return transactions, fmt.Errorf("there is no wallet for user with user_id %v for showing transaction", userID)
	}

	// then find the transactions by wallet_id
	transactions, err = c.orderRepo.FindWalletTransactions(ctx, wallet.ID, pagination)

	if err != nil {
		return transactions, err
	}

	log.Printf("successfully got user transactions for user with user_id %v and wallet_id %v", userID, wallet.ID)

	return transactions, nil
}
