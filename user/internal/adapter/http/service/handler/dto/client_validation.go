package dto

import (
	"regexp"
)

const (
	passwordRegex = `^[A-Za-z\d!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*[A-Z][A-Za-z\d!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?][A-Za-z\d!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*$`
)

func validateClientCreateRequest(req ClientCreateRequest) error {
	validations := []func() error{
		func() error { return validatePassword(req.Password) },
	}

	for _, v := range validations {
		err := v()
		if err != nil {
			return err
		}
	}

	return nil
}

func validatePassword(password string) error {
	re := regexp.MustCompile(passwordRegex)
	if !re.MatchString(password) {
		return ErrInvalidPassword
	}

	return nil
}
