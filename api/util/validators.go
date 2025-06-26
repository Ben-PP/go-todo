package util

import "regexp"

func ValidatePassword(password string) (bool, error) {
	if length := len(password); length < 8 && length > 32 {
		return false, nil
	}
	hasLetter, err := regexp.MatchString(`\p{L}`, password)
	if err != nil {
		// TODO Log
		return false, err
	}
	hasNumber, err := regexp.MatchString(`\d`, password)
	if err != nil {
		// TODO Log
		return false, err
	}
	hasSpecialChar, err := regexp.MatchString(`\p{P}|\p{S}`, password)
	if err != nil {
		// TODO Log
		return false, err
	}

	return hasLetter && hasNumber && hasSpecialChar, nil
}