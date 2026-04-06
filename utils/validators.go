package utils

import (
	"regexp"
	"strings"
	"unicode"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateRegister(name, username, email, password string) []ValidationError {
	var errors []ValidationError

	if !validateName(name) {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must contain only letters (no numbers, spaces, or special characters)",
		})
	}

	if !validateUsername(username) {
		errors = append(errors, ValidationError{
			Field:   "username",
			Message: "Username must start with a letter, can contain letters and numbers (no special characters)",
		})
	}

	if !validateEmail(email) {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Please enter a valid email address",
		})
	}

	if !validatePassword(password, name, username, email) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must be 8+ chars with uppercase, lowercase, number, special char. Must not contain your name, username, or email",
		})
	}

	return errors
}

func validateName(name string) bool {
	if len(name) < 1 {
		return false
	}
	if name[0] == ' ' {
		return false
	}
	nameOnlyLetters := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameOnlyLetters.MatchString(name) {
		return false
	}
	hasLetter := false
	for _, char := range name {
		if unicode.IsLetter(char) {
			hasLetter = true
		} else if unicode.IsDigit(char) {
			return false
		}
	}
	return hasLetter
}

func validateUsername(username string) bool {
	if len(username) < 1 {
		return false
	}
	if !unicode.IsLetter(rune(username[0])) {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`)
	return usernameRegex.MatchString(username)
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}

	blockedDomains := []string{"tempmail.com", "throwaway.com", "10minutemail.com", "guerrillamail.com", "temp-mail.org", "mailinator.com"}
	emailLower := strings.ToLower(email)
	for _, domain := range blockedDomains {
		if strings.HasSuffix(emailLower, "@"+domain) {
			return false
		}
	}

	return true
}

func validatePassword(password, name, username, email string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper, hasLower, hasDigit, hasSpecial := false, false, false, false

	nameParts := strings.Fields(strings.ToLower(name))
	emailPart := strings.Split(strings.ToLower(email), "@")[0]
	usernameLower := strings.ToLower(username)
	nameOrEmailInPassword := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	lowerPassword := strings.ToLower(password)
	for _, part := range nameParts {
		if len(part) >= 3 && strings.Contains(lowerPassword, part) {
			nameOrEmailInPassword = true
		}
	}
	if len(usernameLower) >= 3 && strings.Contains(lowerPassword, usernameLower) {
		nameOrEmailInPassword = true
	}
	if len(emailPart) >= 3 && strings.Contains(lowerPassword, emailPart) {
		nameOrEmailInPassword = true
	}

	weakPattern := regexp.MustCompile(`^[a-zA-Z]{3,}[@#$%^&*]?[0-9]{1,3}$`)
	if weakPattern.MatchString(password) {
		return false
	}

	return hasUpper && hasLower && hasDigit && hasSpecial && !nameOrEmailInPassword
}
