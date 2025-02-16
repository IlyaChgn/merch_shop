//go:generate mockgen -source=auth.go -destination=../mocks/auth_mock.go -package=mocks

package usecases

import (
	"context"
	"github.com/IlyaChgn/merch_shop/internal/models"
)

type AuthStorageInterface interface {
	Auth(ctx context.Context, username, password string) (*models.User, error)
	CreateToken(username string, id uint) (string, error)
	CheckAuth(tokenString string) (*models.User, bool)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}
