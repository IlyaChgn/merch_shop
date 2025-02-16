package responses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IlyaChgn/merch_shop/internal/models"
)

const (
	StatusOk = 200

	StatusBadRequest   = 400
	StatusUnauthorized = 401

	StatusInternalServerError = 500
)

const (
	ErrWrongCredentials = "Wrong password"
	ErrWrongUsername    = "Wrong username"
	ErrWrongJSONFormat  = "Wrong JSON format"
	ErrNotAuthorized    = "Not authorized"

	ErrLowBalance     = "Low balance"
	ErrWrongMerchType = "Wrong item type"
	ErrWrongAmount    = "Amount must be positive number"
	ErrWrongReceiver  = "Receiver must be different from sender"

	ErrInternalServer = "Server error"
)

func newErrResponse(errors string) *models.ErrResponse {
	return &models.ErrResponse{
		Errors: errors,
	}
}

func sendResponse(writer http.ResponseWriter, response any) {
	serverResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Something went wrong while marshalling JSON", err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}

	_, err = writer.Write(serverResponse)
	if err != nil {
		log.Println("Something went wrong while senddng response", err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}
}

func SendOkResponse(writer http.ResponseWriter, body any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(StatusOk)

	sendResponse(writer, body)
}

func SendErrResponse(writer http.ResponseWriter, code int, errors string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	response := newErrResponse(errors)

	sendResponse(writer, response)
}
