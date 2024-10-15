package validation

import (
	"errors"
	"unicode"
)

func ValidateUser(username, password string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if !containsLetterAndNumber(password) {
		return errors.New("password must contain at least one letter and one number")
	}
	return nil
}

func ValidateTransaction(amount float64) error {
	if amount <= 0 {
		return errors.New("transaction amount must be greater than zero")
	}
	return nil
}

func containsLetterAndNumber(s string) bool {
	var hasLetter, hasDigit bool
	for _, c := range s {
		if unicode.IsLetter(c) {
			hasLetter = true
		}
		if unicode.IsDigit(c) {
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}
