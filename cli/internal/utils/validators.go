package utils

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode"
)

// ValidateDomain validates a domain name format
func ValidateDomain(val interface{}) error {
	domain, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid domain type")
	}

	// Basic domain validation
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain format (e.g., example.com)")
	}

	return nil
}

// ValidateSiteID validates a site ID (alphanumeric, 3-16 chars)
func ValidateSiteID(val interface{}) error {
	name, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid site ID type")
	}

	// Alphanumeric only, 3-16 characters
	if len(name) < 3 || len(name) > 16 {
		return fmt.Errorf("site ID must be 3-16 characters")
	}

	alphanumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumRegex.MatchString(name) {
		return fmt.Errorf("site ID must be alphanumeric only")
	}

	return nil
}

// ValidateEmail validates an email address format
func ValidateEmail(val interface{}) error {
	email, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid email type")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidatePasswordStrength validates password complexity
// Requires: minimum 12 characters, at least one uppercase, one lowercase,
// one number, and one special character
func ValidatePasswordStrength(val interface{}) error {
	password, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid password type")
	}

	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return fmt.Errorf("password must contain at least one %s", strings.Join(missing, ", "))
	}

	return nil
}

// ValidateIP validates an IP address format
func ValidateIP(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid type")
	}

	if net.ParseIP(str) == nil {
		return fmt.Errorf("invalid IP address format")
	}

	return nil
}

// ValidatePort validates a port number (1-65535)
func ValidatePort(val interface{}) error {
	port, ok := val.(int)
	if !ok {
		return fmt.Errorf("invalid port type")
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	return nil
}
