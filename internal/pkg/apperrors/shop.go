package apperrors

import "errors"

var (
	WrongMerchTypeError       = errors.New("wrong merch type")
	LowBalanceError           = errors.New("low balance")
	SaveItemError             = errors.New("error occurred while saving item")
	UpdateBalanceError        = errors.New("error occurred while updating balance")
	WriteTransactionError     = errors.New("error occurred while writing transaction")
	GetBalanceError           = errors.New("error occurred while getting balance")
	GetTransactionsError      = errors.New("error occurred while getting transactions")
	ScanningTransactionsError = errors.New("error occurred while scanning transactions")
	NonPositiveAmountError    = errors.New("amount must be positive number")
	WrongReceiverError        = errors.New("receiver must be different from the sender")
)
