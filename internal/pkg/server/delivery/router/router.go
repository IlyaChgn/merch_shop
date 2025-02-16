package router

import (
	authdelivery "github.com/IlyaChgn/merch_shop/internal/pkg/auth/delivery"
	authusecases "github.com/IlyaChgn/merch_shop/internal/pkg/auth/usecases"
	"github.com/IlyaChgn/merch_shop/internal/pkg/middlewares/auth"
	myrecovery "github.com/IlyaChgn/merch_shop/internal/pkg/middlewares/recover"
	shopdelivery "github.com/IlyaChgn/merch_shop/internal/pkg/shop/delivery"
	shopusecases "github.com/IlyaChgn/merch_shop/internal/pkg/shop/usecases"
	"github.com/gorilla/mux"
)

func NewRouter(authStorage authusecases.AuthStorageInterface,
	shopStorage shopusecases.ShopStorageInterface) *mux.Router {
	router := mux.NewRouter()

	router.Use(myrecovery.RecoveryMiddleware)

	authHandler := authdelivery.NewAuthHandler(authStorage)
	shopHandler := shopdelivery.NewShopHandler(shopStorage, authStorage)

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/auth", authHandler.Auth).Methods("POST")

	authMiddleware := auth.AuthMiddleware(authStorage)

	subrouterAuth := apiRouter.PathPrefix("").Subrouter()
	subrouterAuth.Use(authMiddleware)

	subrouterAuth.HandleFunc("/info", shopHandler.GetInfo).Methods("GET")
	subrouterAuth.HandleFunc("/buy/{item}", shopHandler.BuyMerch).Methods("GET")
	subrouterAuth.HandleFunc("/sendCoin", shopHandler.SendCoins).Methods("POST")

	return router
}
