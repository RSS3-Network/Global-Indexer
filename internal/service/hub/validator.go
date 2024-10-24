package hub

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Custom validation for no "://"
func validateNoScheme(fl validator.FieldLevel) bool {
	account := fl.Field().String()
	return !strings.Contains(account, "://")
}
