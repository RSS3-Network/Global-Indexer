package errorx

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:generate go run --mod=mod github.com/dmarkham/enumer --type=ErrorCode --transform=snake --values --trimprefix=ErrorCode --json --output error_code.go
type ErrorCode int

const (
	ErrorCodeBadRequest ErrorCode = iota + 1
	ErrorCodeValidationFailed
	ErrorCodeBadParams
	ErrorCodeInternalError
	ErrorCodeServiceUnavailable
)

var ErrNoNodesAvailable = errors.New("no Nodes are available to process this request")

type ErrorResponse struct {
	Error     string    `json:"error"`
	ErrorCode ErrorCode `json:"error_code"`
	Details   string    `json:"details,omitempty"`
}

func BadRequestError(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, &ErrorResponse{
		ErrorCode: ErrorCodeBadRequest,
		Error:     "Invalid request. Please check your input and try again.",
		Details:   fmt.Sprintf("%v", err),
	})
}

func ValidationFailedError(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, &ErrorResponse{
		ErrorCode: ErrorCodeValidationFailed,
		Error:     "Validation failed. Ensure all fields meet the required criteria and try again.",
		Details:   fmt.Sprintf("%v", err),
	})
}

func BadParamsError(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, &ErrorResponse{
		ErrorCode: ErrorCodeBadParams,
		Error:     "Invalid parameter combination. Verify the combination and try again.",
		Details:   fmt.Sprintf("%v", err),
	})
}

func ServiceUnavailableError(c echo.Context, err error) error {
	return c.JSON(http.StatusServiceUnavailable, &ErrorResponse{
		ErrorCode: ErrorCodeServiceUnavailable,
		Error:     "The requested service is temporarily unavailable, please try again later.",
		Details:   fmt.Sprintf("%v", err),
	})
}

func InternalError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, &ErrorResponse{
		ErrorCode: ErrorCodeInternalError,
		Error:     "An internal error has occurred, please try again later.",
	})
}
