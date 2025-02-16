package repository

import (
	"context"
	"github.com/IlyaChgn/merch_shop/internal/models"
	"github.com/IlyaChgn/merch_shop/internal/pkg/apperrors"
	"github.com/jackc/pgx/v5"
	"log"
)

func (storage *ShopStorage) SendCoins(ctx context.Context, receiverID, senderID uint, amount int) error {
	if amount <= 0 {
		return apperrors.NonPositiveAmountError
	}

	if senderID == receiverID {
		return apperrors.WrongReceiverError
	}

	return storage.sendCoins(ctx, receiverID, senderID, uint(amount))
}

func (storage *ShopStorage) sendCoins(ctx context.Context, receiverID, senderID, amount uint) error {
	tx, err := storage.pool.Begin(ctx)
	if err != nil {
		return apperrors.TxStartError
	}
	defer tx.Rollback(ctx)

	err = storage.decreaseBalance(ctx, tx, amount, senderID)
	if err != nil {
		return err
	}

	err = storage.increaseBalance(ctx, tx, amount, receiverID)
	if err != nil {
		return err
	}

	err = storage.writeTransaction(ctx, tx, amount, receiverID, senderID)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return apperrors.TxCommitError
	}

	return nil
}

func (storage *ShopStorage) getBalance(ctx context.Context, userID uint) (uint, error) {
	var balance uint

	line := storage.pool.QueryRow(ctx, GetBalanceQuery, userID)

	err := line.Scan(&balance)
	if err != nil {
		return 0, apperrors.GetBalanceError
	}

	return balance, nil
}

func (storage *ShopStorage) decreaseBalance(ctx context.Context, tx pgx.Tx, amount, userID uint) error {
	_, err := tx.Exec(ctx, DecreaseBalanceQuery, userID, amount)
	log.Println(err)
	if err != nil {
		return apperrors.LowBalanceError
	}

	return nil
}

func (storage *ShopStorage) increaseBalance(ctx context.Context, tx pgx.Tx, amount, userID uint) error {
	_, err := tx.Exec(ctx, IncreaseBalanceQuery, userID, amount)
	if err != nil {
		return apperrors.UpdateBalanceError
	}

	return nil
}

func (storage *ShopStorage) writeTransaction(ctx context.Context, tx pgx.Tx, amount, receiverID, senderID uint) error {
	_, err := tx.Exec(ctx, AddTransactionQuery, senderID, receiverID, amount)
	if err != nil {
		return apperrors.WriteTransactionError
	}

	return nil
}

func (storage *ShopStorage) getIncomingTransactions(ctx context.Context, userID uint) (*[]models.ReceivedCoinsItem, error) {
	var list []models.ReceivedCoinsItem

	rows, err := storage.pool.Query(ctx, GetIncomingTransactionsQuery, userID)
	if err != nil {
		return nil, apperrors.GetTransactionsError
	}
	defer rows.Close()

	for rows.Next() {
		var item models.ReceivedCoinsItem

		if err := rows.Scan(&item.Amount, &item.FromUser); err != nil {
			return nil, apperrors.ScanningTransactionsError
		}

		list = append(list, item)
	}

	return &list, nil
}

func (storage *ShopStorage) getOutgoingTransactions(ctx context.Context, userID uint) (*[]models.SentCoinsItem, error) {
	var list []models.SentCoinsItem

	rows, err := storage.pool.Query(ctx, GetOutgoingTransactionsQuery, userID)
	if err != nil {
		return nil, apperrors.GetTransactionsError
	}
	defer rows.Close()

	for rows.Next() {
		var item models.SentCoinsItem

		if err := rows.Scan(&item.Amount, &item.ToUser); err != nil {
			return nil, apperrors.ScanningTransactionsError
		}

		list = append(list, item)
	}

	return &list, nil
}
