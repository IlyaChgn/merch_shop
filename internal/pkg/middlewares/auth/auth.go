package auth

import (
	authusecases "github.com/IlyaChgn/merch_shop/internal/pkg/auth/usecases"
	"github.com/IlyaChgn/merch_shop/internal/pkg/server/delivery/responses"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthMiddleware(authStorage authusecases.AuthStorageInterface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			tokenString := req.Header.Get("Authorization")

			_, valid := authStorage.CheckAuth(tokenString)
			if !valid {
				responses.SendErrResponse(writer, responses.StatusUnauthorized, responses.ErrNotAuthorized)

				return
			}

			next.ServeHTTP(writer, req)
		})
	}
}
