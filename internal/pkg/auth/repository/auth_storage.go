package repository

import (
	"context"
	"fmt"
	"github.com/IlyaChgn/merch_shop/internal/models"
	"github.com/IlyaChgn/merch_shop/internal/pkg/apperrors"
	pool "github.com/IlyaChgn/merch_shop/internal/pkg/server/repository"
	"github.com/IlyaChgn/merch_shop/internal/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

const tokenDuration = 72 * time.Hour

type UserClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"sub"`
	jwt.RegisteredClaims
}

type AuthStorage struct {
	pool      pool.PostgresPool
	secretKey []byte
}

func NewAuthStorage(pool pool.PostgresPool, secretKey string) *AuthStorage {
	return &AuthStorage{
		pool:      pool,
		secretKey: []byte(secretKey),
	}
}

func (storage *AuthStorage) Auth(ctx context.Context, username, password string) (*models.User, error) {
	passwordHash := utils.HashPassword(password)

	user, err := storage.getUserByUsername(ctx, username)
	if err != nil {
		user, err = storage.createUser(ctx, username, passwordHash)
		if err != nil {
			return nil, err
		}
	} else {
		if !utils.CheckPassword(password, user.PasswordHash) {
			return nil, apperrors.WrongPasswordError
		}
	}

	return user, nil
}

func (storage *AuthStorage) CreateToken(username string, id uint) (string, error) {
	payload := jwt.MapClaims{
		"sub":     username,
		"user_id": strconv.Itoa(int(id)),
		"exp":     time.Now().Add(tokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tok, err := token.SignedString(storage.secretKey)
	if err != nil {
		return "", err
	}

	return tok, nil
}

func (storage *AuthStorage) CheckAuth(tokenString string) (*models.User, bool) {
	if tokenString == "" {
		return nil, false
	}

	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return storage.secretKey, nil
	})

	userID, convErr := strconv.Atoi(claims.UserID)
	if convErr != nil {
		return nil, false
	}

	user := &models.User{
		Username: claims.Username,
		ID:       uint(userID),
	}

	return user, err == nil && token.Valid
}

func (storage *AuthStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return storage.getUserByUsername(ctx, username)
}

func (storage *AuthStorage) getUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	line := storage.pool.QueryRow(ctx, GetUserByUsernameQuery, username)
	if err := line.Scan(&user.ID, &user.Username, &user.PasswordHash); err != nil {
		return nil, err
	}

	return &user, nil
}

func (storage *AuthStorage) createUser(ctx context.Context, username, passwordHash string) (*models.User, error) {
	tx, err := storage.pool.Begin(ctx)
	if err != nil {
		return nil, apperrors.TxStartError
	}
	defer tx.Rollback(ctx)

	var user models.User

	line := tx.QueryRow(ctx, CreateUserQuery, username, passwordHash)
	if err := line.Scan(&user.ID, &user.Username, &user.PasswordHash); err != nil {
		return nil, apperrors.TxCreateUserError
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperrors.TxCommitError
	}

	return &user, nil
}
