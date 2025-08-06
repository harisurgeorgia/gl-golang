package validation

import (
	"errors"
	"fmt"
	"unicode"
)

func IsValidPassword(pw string) error {
	if len(pw) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	var hasUpper, hasLower, hasSpecial bool
	for _, ch := range pw {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsSymbol(ch), unicode.IsPunct(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil // valid password
}

func CheckPasswordMatch(password, confirmPassword string) error {
	if (password != confirmPassword) || (password == "") {
		return fmt.Errorf("password does not match or empty")
		//return utils.Render(c, 400, views.Layout(data, views.UserForm(data.Header, "Passwords do not match or are empty.", user)))

	}
	return nil
}
