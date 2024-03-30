package validate

import (
	"regexp"
	"strings"
)

const passwordMinLen = 8

// validates user's password
func ValidatePassword(pwd string) []string {
	if pwd == "" {
		return []string{"password is required"}
	}

	errs := []string{}

	if len(pwd) < passwordMinLen {
		errs = append(errs, "password should contain at least 8 characters")
	}

	if pwd == strings.ToLower(pwd) {
		errs = append(errs, "password should contain at least one uppercase letter")
	}

	if pwd == strings.ToUpper(pwd) {
		errs = append(errs, "password should contain at least one lowercase letter")
	}

	numeric := regexp.MustCompile(`\d`).MatchString(pwd)
	if !numeric {
		errs = append(errs, "password should contain at least one digit")
	}

	alphabetic := regexp.MustCompile(`[A-Za-z_]`).MatchString(pwd)
	if !alphabetic {
		errs = append(errs, "password should contain at least one latin letter")
	}

	return errs
}
