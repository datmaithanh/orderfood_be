package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(value string, minLen int, maxLen int) error {
	n := len(value)
	if n < minLen || n > maxLen {
		return fmt.Errorf("string length must be between %d and %d", minLen, maxLen)
	}
	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(username) {
		return fmt.Errorf("username can only contain alphanumeric characters and underscores")
	}
	return nil
}

func ValidateFullName(fullName string) error {
	if err := ValidateString(fullName, 3, 100); err != nil {
		return err
	}
	if !isValidFullname(fullName) {
		return fmt.Errorf("full name can only contain alphabetic characters and spaces")
	}
	return nil
}

func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

func ValidateEmail (email string) error {
	if err := ValidateString(email, 3, 200); err != nil {
		return err
	}

	_ ,err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
