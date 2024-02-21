package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendJSONError(c echo.Context, code int) error {
	// we do not throw the detail error code to prevent unwanted leaks
	msg := http.StatusText(code)

	if msg == "" { // In case a non-standard code is used
		msg = http.StatusText(http.StatusInternalServerError)
	}

	return c.JSON(code, echo.Map{"msg": msg})
}
