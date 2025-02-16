package repository_test

import (
	"context"
	"errors"
	"github.com/IlyaChgn/merch_shop/internal/models"
	pgxmocks "github.com/IlyaChgn/merch_shop/internal/pkg/server/mocks"
	"github.com/IlyaChgn/merch_shop/internal/pkg/shop/repository"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShopStorage_SendCoins(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxmocks.NewMockPostgresPool(ctrl)

	var (
		receiverID = uint(1)
		senderID   = uint(2)
		amount     = 100
	)

	tests := []struct {
		name        string
		receiverID  uint
		senderID    uint
		amount      int
		expectedErr bool
		setup       func()
	}{
		{
			name:        "negative coins amount",
			receiverID:  receiverID,
			senderID:    senderID,
			amount:      -100,
			expectedErr: true,
			setup:       func() {},
		},
		{
			name:        "sender same as receiver",
			receiverID:  receiverID,
			senderID:    receiverID,
			amount:      amount,
			expectedErr: true,
			setup:       func() {},
		},
		{
			name:        "successful case",
			receiverID:  receiverID,
			senderID:    senderID,
			amount:      amount,
			expectedErr: false,
			setup: func() {
				mockTx := pgxmocks.NewMockTx(ctrl)
				mockPool.EXPECT().
					Begin(context.Background()).
					Return(mockTx, nil)

				mockTx.EXPECT().
					Exec(context.Background(), repository.DecreaseBalanceQuery, senderID, uint(amount)).
					Return(pgconn.CommandTag{}, nil)
				mockTx.EXPECT().
					Exec(context.Background(), repository.IncreaseBalanceQuery, receiverID, uint(amount)).
					Return(pgconn.CommandTag{}, nil)
				mockTx.EXPECT().
					Exec(context.Background(), repository.AddTransactionQuery, senderID, receiverID, uint(amount)).
					Return(pgconn.CommandTag{}, nil)

				mockTx.EXPECT().
					Commit(context.Background()).
					Return(nil)
				mockTx.EXPECT().
					Rollback(context.Background()).
					Return(nil)
			},
		},
		{
			name:        "low balance",
			receiverID:  receiverID,
			senderID:    senderID,
			amount:      amount,
			expectedErr: true,
			setup: func() {
				mockTx := pgxmocks.NewMockTx(ctrl)
				mockPool.EXPECT().
					Begin(context.Background()).
					Return(mockTx, nil)

				mockTx.EXPECT().
					Exec(context.Background(), repository.DecreaseBalanceQuery, senderID, uint(amount)).
					Return(pgconn.CommandTag{}, errors.New(""))

				mockTx.EXPECT().
					Rollback(context.Background()).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewShopStorage(mockPool)

			err := storage.SendCoins(context.Background(), tt.receiverID, tt.senderID, tt.amount)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestShopStorage_BuyItem(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxmocks.NewMockPostgresPool(ctrl)

	var (
		userID = uint(1)
		item   = "item"
		price  = uint(100)
		itemID = uint(1)
	)

	tests := []struct {
		name        string
		item        string
		userID      uint
		expectedErr bool
		setup       func()
	}{
		{
			name:        "wrong item type",
			item:        item,
			userID:      userID,
			expectedErr: true,
			setup: func() {
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetMerchItemQuery, item).
					Return(models.EmptyRow{})
			},
		},
		{
			name:        "low balance",
			item:        item,
			userID:      userID,
			expectedErr: true,
			setup: func() {
				pgxRow := pgxpoolmock.NewRow(itemID, item, price)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetMerchItemQuery, item).
					Return(pgxRow)

				mockTx := pgxmocks.NewMockTx(ctrl)
				mockPool.EXPECT().
					Begin(context.Background()).
					Return(mockTx, nil)

				mockTx.EXPECT().
					Exec(context.Background(), repository.SaveItemQuery, itemID, userID).
					Return(pgconn.CommandTag{}, nil)
				mockTx.EXPECT().
					Exec(context.Background(), repository.DecreaseBalanceQuery, userID, uint(price)).
					Return(pgconn.CommandTag{}, errors.New(""))

				mockTx.EXPECT().
					Rollback(context.Background()).
					Return(nil)
			},
		},
		{
			name:        "successful case",
			item:        item,
			userID:      userID,
			expectedErr: false,
			setup: func() {
				pgxRow := pgxpoolmock.NewRow(itemID, item, price)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetMerchItemQuery, item).
					Return(pgxRow)

				mockTx := pgxmocks.NewMockTx(ctrl)
				mockPool.EXPECT().
					Begin(context.Background()).
					Return(mockTx, nil)

				mockTx.EXPECT().
					Exec(context.Background(), repository.SaveItemQuery, itemID, userID).
					Return(pgconn.CommandTag{}, nil)
				mockTx.EXPECT().
					Exec(context.Background(), repository.DecreaseBalanceQuery, userID, uint(price)).
					Return(pgconn.CommandTag{}, nil)

				mockTx.EXPECT().
					Commit(context.Background()).
					Return(nil)
				mockTx.EXPECT().
					Rollback(context.Background()).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewShopStorage(mockPool)

			err := storage.BuyItem(context.Background(), tt.item, tt.userID)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
