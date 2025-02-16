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
	shopdel "github.com/IlyaChgn/merch_shop/internal/pkg/shop/delivery"
	shoprepo "github.com/IlyaChgn/merch_shop/internal/pkg/shop/repository"
	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSendCoin(t *testing.T) {
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
	shopStorage := shoprepo.NewShopStorage(postgresPool)

	shopHandler := shopdel.NewShopHandler(shopStorage, authStorage)
	authHandler := authdel.NewAuthHandler(authStorage)

	authReqBody1 := models.AuthRequest{
		Username: "user1",
		Password: "password",
	}

	authReqJSON1, err := json.Marshal(authReqBody1)
	if err != nil {
		t.Fatalf("Failed to marshal JSON %v", err)
	}

	authReq1, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authReqJSON1))
	if err != nil {
		t.Fatalf("failed to create auth request %v", err)
	}
	authReq1.Header.Add("Content-Type", "application/json")

	authWriter1 := httptest.NewRecorder()
	authHandle1 := http.HandlerFunc(authHandler.Auth)

	authHandle1.ServeHTTP(authWriter1, authReq1)

	if authWriter1.Code != responses.StatusOk {
		t.Errorf("Expected status OK, got %v", authWriter1.Code)
	}

	authReqBody2 := models.AuthRequest{
		Username: "user2",
		Password: "password",
	}

	authReqJSON2, err := json.Marshal(authReqBody2)
	if err != nil {
		t.Fatalf("Failed to marshal JSON %v", err)
	}

	authReq2, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authReqJSON2))
	if err != nil {
		t.Fatalf("failed to create auth request %v", err)
	}
	authReq2.Header.Add("Content-Type", "application/json")

	authWriter2 := httptest.NewRecorder()
	authHandle2 := http.HandlerFunc(authHandler.Auth)

	authHandle2.ServeHTTP(authWriter2, authReq2)

	var tokenStruct *models.AuthResponse
	err = json.NewDecoder(authWriter2.Body).Decode(&tokenStruct)
	if err != nil {
		t.Fatal(err)
	}

	if authWriter2.Code != responses.StatusOk {
		t.Errorf("Expected status OK, got %v", authWriter2.Code)
	}
	shopReqBody := models.SendCoinRequest{
		ToUser: "user1",
		Amount: 100,
	}

	shopReqJSON, err := json.Marshal(shopReqBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON %v", err)
	}

	shopReq, err := http.NewRequest("GET", "/api/sendCoin", bytes.NewBuffer(shopReqJSON))
	if err != nil {
		t.Fatalf("failed to create send coin request %v", err)
	}
	shopReq.Header.Set("Authorization", tokenStruct.Token)

	shopWriter := httptest.NewRecorder()
	shopHandle := http.HandlerFunc(shopHandler.SendCoins)

	shopHandle.ServeHTTP(shopWriter, shopReq)

	if shopWriter.Code != responses.StatusOk {
		t.Errorf("Expected status OK, got %v", shopWriter.Code)
	}
}
