package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/IlyaChgn/merch_shop/internal/models"
	authdel "github.com/IlyaChgn/merch_shop/internal/pkg/auth/delivery"
	authrepo "github.com/IlyaChgn/merch_shop/internal/pkg/auth/repository"
	"github.com/IlyaChgn/merch_shop/internal/pkg/server/delivery/responses"
	pool "github.com/IlyaChgn/merch_shop/internal/pkg/server/repository"
	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Fatalf("Error loading env file %v", err)
	}

	postgresURL := pool.NewConnectionString(os.Getenv("TEST_POSTGRES_USER"),
		os.Getenv("TEST_POSTGRES_PASSWORD"), os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))

	postgresPool, err := pool.NewPostgresPool(postgresURL)
	if err != nil {
		t.Fatalf("Something went wrong while creating postgres pool %v", err)
	}

	err = postgresPool.Ping(context.Background())
	if err != nil {
		t.Fatalf("Cannot connect to postgres database %v", err)
	}

	_, err = postgresPool.Exec(context.Background(), "TRUNCATE public.user, public.balance, public.inventory_item,"+
		"public.transaction RESTART IDENTITY")
	if err != nil {
		t.Fatalf("Cannot truncate DB %v", err)
	}

	authStorage := authrepo.NewAuthStorage(postgresPool, os.Getenv("TEST_SECRET_KEY"))

	authHandler := authdel.NewAuthHandler(authStorage)

	authReqBody := models.AuthRequest{
		Username: "user1",
		Password: "password",
	}

	authReqJSON, err := json.Marshal(authReqBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON %v", err)
	}

	authReq, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authReqJSON))
	if err != nil {
		t.Fatalf("failed to create auth request %v", err)
	}
	authReq.Header.Add("Content-Type", "application/json")

	authWriter := httptest.NewRecorder()
	authHandle := http.HandlerFunc(authHandler.Auth)

	authHandle.ServeHTTP(authWriter, authReq)

	if authWriter.Code != responses.StatusOk {
		t.Errorf("Expected status OK, got %v", authWriter.Code)
	}
}
