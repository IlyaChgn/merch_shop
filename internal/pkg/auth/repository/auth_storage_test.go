package repository_test

import (
	"context"
	"github.com/IlyaChgn/merch_shop/internal/models"
	"github.com/IlyaChgn/merch_shop/internal/pkg/auth/repository"
	pgxmocks "github.com/IlyaChgn/merch_shop/internal/pkg/server/mocks"
	"github.com/IlyaChgn/merch_shop/internal/pkg/utils"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthStorage_Auth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxmocks.NewMockPostgresPool(ctrl)

	var (
		username = "username"
		password = "password"
		hash     = utils.HashPassword(password)
		id       = 1
	)

	tests := []struct {
		name        string
		username    string
		password    string
		expected    *models.User
		expectedErr bool
		setup       func()
	}{
		{
			name:     "existing user with correct password",
			username: username,
			password: password,
			expected: &models.User{
				ID:           uint(id),
				Username:     username,
				PasswordHash: hash,
			},
			expectedErr: false,
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(uint(id), username, hash)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByUsernameQuery, username).
					Return(pgxRows)
			},
		},
		{
			name:        "existing user with incorrect password",
			username:    username,
			password:    "fakePassword",
			expected:    nil,
			expectedErr: true,
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(uint(id), username, hash)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByUsernameQuery, username).
					Return(pgxRows)
			},
		},
		{
			name:     "non-existing user",
			username: username,
			password: password,
			expected: &models.User{
				ID:           uint(id),
				Username:     username,
				PasswordHash: hash,
			},
			expectedErr: false,
			setup: func() {
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByUsernameQuery, username).
					Return(models.EmptyRow{})

				mockTx := pgxmocks.NewMockTx(ctrl)
				mockPool.EXPECT().
					Begin(context.Background()).
					Return(mockTx, nil)

				pgxRows := pgxpoolmock.NewRow(uint(id), username, hash)
				mockTx.EXPECT().
					QueryRow(context.Background(), repository.CreateUserQuery, username, gomock.Any()).
					Return(pgxRows)
				mockTx.EXPECT().
					Commit(context.Background()).
					Return(nil)
				mockTx.EXPECT().
					Rollback(context.Background()).
					Return(nil)
			},
		},
		{
			name:        "tx error",
			username:    username,
			password:    password,
			expected:    nil,
			expectedErr: true,
			setup: func() {
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByUsernameQuery, username).
					Return(models.EmptyRow{})

				mockTx := pgxmocks.NewMockTx(ctrl)
				mockPool.EXPECT().
					Begin(context.Background()).
					Return(mockTx, nil)

				mockTx.EXPECT().
					QueryRow(context.Background(), repository.CreateUserQuery, username, gomock.Any()).
					Return(models.EmptyRow{})
				mockTx.EXPECT().
					Rollback(context.Background()).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewAuthStorage(mockPool, "secret")

			got, err := storage.Auth(context.Background(), tt.username, tt.password)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestAuthStorage_GetUserByUsername(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxmocks.NewMockPostgresPool(ctrl)

	var (
		username = "username"
		password = "password"
		hash     = utils.HashPassword(password)
		id       = 1
	)

	tests := []struct {
		name        string
		username    string
		expected    *models.User
		expectedErr bool
		setup       func()
	}{
		{
			name:     "existing user",
			username: username,
			expected: &models.User{
				ID:           uint(id),
				Username:     username,
				PasswordHash: hash,
			},
			expectedErr: false,
			setup: func() {
				pgxRows := pgxpoolmock.NewRow(uint(id), username, hash)
				mockPool.EXPECT().
					QueryRow(context.Background(), repository.GetUserByUsernameQuery, username).
					Return(pgxRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			storage := repository.NewAuthStorage(mockPool, "secret")

			got, err := storage.GetUserByUsername(context.Background(), tt.username)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
