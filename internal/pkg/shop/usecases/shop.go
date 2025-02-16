//go:generate mockgen -source=shop.go -destination=../mocks/shop_mock.go -package=mocks

package usecases

import (
	"context"
	"github.com/IlyaChgn/merch_shop/internal/models"
)

type ShopStorageInterface interface {
	BuyItem(ctx context.Context, item string, userID uint) error
	SendCoins(ctx context.Context, receiverID, senderID uint, amount int) error
	GetInfo(ctx context.Context, userID uint) (*models.InfoResponse, error)
}
