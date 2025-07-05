package validate

import (
	"fmt"
	"regexp"
)

// Returns true if the str has lower or equal number of chars than length.
func stringLength(str string, length int) bool {
	return len(str) <= length
}

func LengthListDescription(txt string) bool {
	return stringLength(txt, 100)
}

func LengthListTitle(txt string) bool {
	return stringLength(txt, 50)
}

func LengthTodoTitle(txt string) bool {
	return stringLength(txt, 50)
}

func LengthTodoDescription(txt string) bool {
	return stringLength(txt, 200)
}

func Password(password string) (bool, error) {
	if length := len(password); length < 8 || length > 32 {
		return false, nil
	}
	hasLetter, err := regexp.MatchString(`\p{L}`, password)
	if err != nil {
		return false, fmt.Errorf("failed to match letter regex: %w", err)
	}
	hasNumber, err := regexp.MatchString(`\d`, password)
	if err != nil {
		return false, fmt.Errorf("failed to match number regex: %w", err)
	}
	hasSpecialChar, err := regexp.MatchString(`\p{P}|\p{S}`, password)
	if err != nil {
		return false, fmt.Errorf("failed to match punctuation regex: %w", err)
	}

	return hasLetter && hasNumber && hasSpecialChar, nil
}

func Username(username string) (bool, error) {
	if length := len(username); length < 3 || length > 20 {
		return false, nil
	}
	hasDisallowedChars, err := regexp.MatchString(`[^\p{L}\p{N}\s_-]`, username)
	if err != nil {
		return false, err
	}

	return !hasDisallowedChars, nil
}
