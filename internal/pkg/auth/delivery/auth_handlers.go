package delivery

import (
	"encoding/json"
	"errors"
	"github.com/IlyaChgn/merch_shop/internal/models"
	"github.com/IlyaChgn/merch_shop/internal/pkg/apperrors"
	"github.com/IlyaChgn/merch_shop/internal/pkg/auth/usecases"
	"github.com/IlyaChgn/merch_shop/internal/pkg/server/delivery/responses"
	"log"
	"net/http"
)

type AuthHandler struct {
	storage usecases.AuthStorageInterface
}

func NewAuthHandler(storage usecases.AuthStorageInterface) *AuthHandler {
	return &AuthHandler{
		storage: storage,
	}
}

func (handler *AuthHandler) Auth(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var authData models.AuthRequest

	err := json.NewDecoder(req.Body).Decode(&authData)
	if err != nil {
		responses.SendErrResponse(writer, responses.StatusBadRequest, responses.ErrWrongJSONFormat)

		return
	}

	user, err := handler.storage.Auth(ctx, authData.Username, authData.Password)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.WrongPasswordError):
			responses.SendErrResponse(writer, responses.StatusUnauthorized, responses.ErrWrongCredentials)
		default:
			log.Println(err)

			responses.SendErrResponse(writer, responses.StatusInternalServerError, responses.ErrInternalServer)
		}

		return
	}

	token, err := handler.storage.CreateToken(user.Username, user.ID)
	if err != nil {
		responses.SendErrResponse(writer, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	responses.SendOkResponse(writer, &models.AuthResponse{Token: token})
}
