package validate

import (
	"regexp"
	"strings"
)

const passwordMinLen = 8

const (
	Weak    = "weak"
	Good    = "good"
	Perfect = "perfect"
)

// validates user's password
// returns "weak" if password length < passwordMinLen or password doesn't contain digits
// returns "perfect" if password length >= passwordMinLen and  password contains uppercase and lowercase latin letters
// and contains special symbols
// otherwise returns "good"
func ValidatePassword(pwd string) string {

	if len(pwd) < passwordMinLen {
		return Weak
	}

	numeric := regexp.MustCompile(`\d`).MatchString(pwd)
	if !numeric {
		return Weak
	}

	alphabetic := regexp.MustCompile(`[A-Za-z_]`).MatchString(pwd)

	specsymb := regexp.MustCompile(`\^|\$|\.|\||\?|\*|\+|\-|\[|\]|\{|\}|\(|\)|\\|\!`).MatchString(pwd)

	if specsymb && alphabetic && pwd != strings.ToUpper(pwd) && pwd != strings.ToLower(pwd) {
		return Perfect
	}

	return Good
}
