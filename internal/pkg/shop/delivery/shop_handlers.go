package delivery

import (
	"encoding/json"
	"errors"
	"github.com/IlyaChgn/merch_shop/internal/models"
	"github.com/IlyaChgn/merch_shop/internal/pkg/apperrors"
	authusecases "github.com/IlyaChgn/merch_shop/internal/pkg/auth/usecases"
	"github.com/IlyaChgn/merch_shop/internal/pkg/server/delivery/responses"
	"github.com/IlyaChgn/merch_shop/internal/pkg/shop/usecases"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ShopHandler struct {
	storage     usecases.ShopStorageInterface
	authStorage authusecases.AuthStorageInterface
}

func NewShopHandler(storage usecases.ShopStorageInterface, authStorage authusecases.AuthStorageInterface) *ShopHandler {
	return &ShopHandler{
		storage:     storage,
		authStorage: authStorage,
	}
}

func (handler *ShopHandler) BuyMerch(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	vars := mux.Vars(req)
	item := vars["item"]

	user, _ := handler.authStorage.CheckAuth(req.Header.Get("Authorization"))

	err := handler.storage.BuyItem(ctx, item, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.LowBalanceError):
			responses.SendErrResponse(writer, responses.StatusBadRequest, responses.ErrLowBalance)
		case errors.Is(err, apperrors.WrongMerchTypeError):
			responses.SendErrResponse(writer, responses.StatusBadRequest, responses.ErrWrongMerchType)
		default:
			log.Println(err)

			responses.SendErrResponse(writer, responses.StatusInternalServerError, responses.ErrInternalServer)
		}

		return
	}
}

func (handler *ShopHandler) SendCoins(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var reqData models.SendCoinRequest

	err := json.NewDecoder(req.Body).Decode(&reqData)
	if err != nil {
		responses.SendErrResponse(writer, responses.StatusBadRequest, responses.ErrWrongJSONFormat)

		return
	}

	sender, _ := handler.authStorage.CheckAuth(req.Header.Get("Authorization"))

	receiver, err := handler.authStorage.GetUserByUsername(ctx, reqData.ToUser)
	if err != nil {
		responses.SendErrResponse(writer, responses.StatusBadRequest, responses.ErrWrongUsername)

		return
	}

	err = handler.storage.SendCoins(ctx, receiver.ID, sender.ID, reqData.Amount)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.LowBalanceError):
			responses.SendErrResponse(writer, responses.StatusBadRequest, responses.ErrLowBalance)
		case errors.Is(err, apperrors.WrongReceiverError):
			responses.SendErrResponse(writer, responses.StatusBadRequest, responses.ErrWrongReceiver)
		case errors.Is(err, apperrors.NonPositiveAmountError):
			responses.SendErrResponse(writer, responses.StatusBadRequest, responses.ErrWrongAmount)
		default:
			log.Println(err)

			responses.SendErrResponse(writer, responses.StatusInternalServerError, responses.ErrInternalServer)
		}

		return
	}
}

func (handler *ShopHandler) GetInfo(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	user, _ := handler.authStorage.CheckAuth(req.Header.Get("Authorization"))

	info, err := handler.storage.GetInfo(ctx, user.ID)
	if err != nil {
		log.Println(err)

		responses.SendErrResponse(writer, responses.StatusInternalServerError, responses.ErrInternalServer)

		return
	}

	responses.SendOkResponse(writer, info)
}
