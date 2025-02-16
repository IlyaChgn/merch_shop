package repository

import (
	"context"
	"github.com/IlyaChgn/merch_shop/internal/models"
	"github.com/IlyaChgn/merch_shop/internal/pkg/apperrors"
	pool "github.com/IlyaChgn/merch_shop/internal/pkg/server/repository"
	"github.com/jackc/pgx/v5"
)

type ShopStorage struct {
	pool pool.PostgresPool
}

func NewShopStorage(pool pool.PostgresPool) *ShopStorage {
	return &ShopStorage{
		pool: pool,
	}
}

func (storage *ShopStorage) BuyItem(ctx context.Context, item string, userID uint) error {
	merchItem, err := storage.getMerchItem(ctx, item)
	if err != nil {
		return err
	}

	return storage.buyItem(ctx, merchItem, userID)
}

func (storage *ShopStorage) GetInfo(ctx context.Context, userID uint) (*models.InfoResponse, error) {
	balance, err := storage.getBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	inventory, err := storage.getInventory(ctx, userID)
	if err != nil {
		return nil, err
	}

	sentCoins, err := storage.getOutgoingTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	receivedCoins, err := storage.getIncomingTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.InfoResponse{
		Coins:     int(balance),
		Inventory: *inventory,
		CoinHistory: models.CoinHistory{
			Sent:     *sentCoins,
			Received: *receivedCoins,
		},
	}, nil
}

func (storage *ShopStorage) buyItem(ctx context.Context, merchItem *models.MerchItem, userID uint) error {
	tx, err := storage.pool.Begin(ctx)
	if err != nil {
		return apperrors.TxStartError
	}
	defer tx.Rollback(ctx)

	err = storage.saveItem(ctx, tx, merchItem.ID, userID)
	if err != nil {
		return err
	}

	err = storage.decreaseBalance(ctx, tx, merchItem.Price, userID)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return apperrors.TxCommitError
	}

	return nil
}

func (storage *ShopStorage) getMerchItem(ctx context.Context, item string) (*models.MerchItem, error) {
	var merchItem models.MerchItem

	line := storage.pool.QueryRow(ctx, GetMerchItemQuery, item)
	if err := line.Scan(&merchItem.ID, &merchItem.Name, &merchItem.Price); err != nil {
		return nil, apperrors.WrongMerchTypeError
	}

	return &merchItem, nil
}

func (storage *ShopStorage) saveItem(ctx context.Context, tx pgx.Tx, itemID, userID uint) error {
	_, err := tx.Exec(ctx, SaveItemQuery, userID, itemID)
	if err != nil {
		return apperrors.SaveItemError
	}

	return nil
}

func (storage *ShopStorage) getInventory(ctx context.Context, userID uint) (*[]models.InventoryItem, error) {
	var list []models.InventoryItem

	rows, err := storage.pool.Query(ctx, GetInventoryQuery, userID)
	if err != nil {
		return nil, apperrors.GetTransactionsError
	}
	defer rows.Close()

	for rows.Next() {
		var item models.InventoryItem

		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, apperrors.ScanningTransactionsError
		}

		list = append(list, item)
	}

	return &list, nil
}
