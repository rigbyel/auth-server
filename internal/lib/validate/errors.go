package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationErrors(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "email":
			errMsgs = append(errMsgs, "invalid email")
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))

		}
	}

	return strings.Join(errMsgs, ", ")
}
