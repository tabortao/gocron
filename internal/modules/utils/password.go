package utils

import (
	"regexp"
	"unicode"
)

// PasswordMinLength 密码最小长度
const PasswordMinLength = 8

// ValidatePassword 验证密码复杂度
// 要求：至少8位，包含字母和数字
func ValidatePassword(password string) (bool, string) {
	if len(password) < PasswordMinLength {
		return false, "password_min_length_8"
	}

	hasLetter := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return false, "password_must_contain_letter_and_digit"
	}

	return true, ""
}

// ValidatePasswordStrong 验证强密码
// 要求：至少8位，包含大小写字母、数字和特殊字符
func ValidatePasswordStrong(password string) (bool, string) {
	if len(password) < PasswordMinLength {
		return false, "password_min_length_8"
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	specialChars := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if specialChars.MatchString(password) {
		hasSpecial = true
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return false, "password_must_contain_upper_lower_digit_special"
	}

	return true, ""
}
