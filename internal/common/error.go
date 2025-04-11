package common

import "errors"

var (

	// common errors
	ErrorNotFound      = errors.New("not found")
	ErrorAlreadyExists = errors.New("already exists")
	ErrorValidation    = errors.New("validation error")

	// auth-specific errors
	ErrorInvalidAuthheaderFormat = errors.New("invalid auth header format")
	ErrorInvalidToken            = errors.New("invalid token")
	ErrorNoUserID                = errors.New("no user id")
	ErrorLoginAlreadyExists      = errors.New("login already exists")
	ErrorInvalidLoginFormat      = errors.New("invalid login format")
	ErrorInvalidPasswordFormat   = errors.New("invalid password format")
	ErrorInvalidLoginPassword    = errors.New("invalid login/password")

	// order-specific errors
	ErrorNoOrderNumberSpecified   = errors.New("no order number specified")
	ErrorInvalidOrderNumberFormat = errors.New("invalid order number format")
	ErrorOrderDoesNotExist        = errors.New("order does not exist")
	ErrorOrderAlreadyExists       = errors.New("order already exists")

	// balance-specific errors
	ErrorInsufficientBalance = errors.New("insufficient balance")
)
