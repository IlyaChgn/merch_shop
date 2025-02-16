package apperrors

import "errors"

var (
	TxStartError      = errors.New("failed to start transaction")
	TxCreateUserError = errors.New("error creating user in transaction")
	TxCommitError     = errors.New("failed to commit transaction")

	WrongPasswordError = errors.New("wrong password")
)
