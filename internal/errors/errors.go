package errors

import "errors"

var (
	ErrWrongMethod              = errors.New("no methods except POST allowed")
	ErrContentLengthNotProvided = errors.New("content-length is required")
	ErrBodyTooLarge             = errors.New("body is too large")
	ErrWrongBodySyntax          = errors.New("body syntax is invalid")
	ErrInternalServer 			= errors.New("server error")

	ErrNotNumber 				= errors.New("element is not a number")
)