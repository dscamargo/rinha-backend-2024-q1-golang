package apperrors

import "errors"

var (
	ErrBalanceIsNull = errors.New("client balance is null")
	ErrInvalidBody   = errors.New("invalid body")
)
