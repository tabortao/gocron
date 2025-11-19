package utils

import "testing"

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
		errKey   string
	}{
		{"abc123", false, "password_min_length_8"},
		{"abcdefgh", false, "password_must_contain_letter_and_digit"},
		{"12345678", false, "password_must_contain_letter_and_digit"},
		{"abc12345", true, ""},
		{"Test1234", true, ""},
		{"password123", true, ""},
	}

	for _, tt := range tests {
		valid, errKey := ValidatePassword(tt.password)
		if valid != tt.valid {
			t.Errorf("ValidatePassword(%q) valid = %v, want %v", tt.password, valid, tt.valid)
		}
		if errKey != tt.errKey {
			t.Errorf("ValidatePassword(%q) errKey = %q, want %q", tt.password, errKey, tt.errKey)
		}
	}
}

func TestValidatePasswordStrong(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
		errKey   string
	}{
		{"abc123", false, "password_min_length_8"},
		{"Abcd1234", false, "password_must_contain_upper_lower_digit_special"},
		{"Abcd123!", true, ""},
		{"Test@123", true, ""},
		{"P@ssw0rd", true, ""},
	}

	for _, tt := range tests {
		valid, errKey := ValidatePasswordStrong(tt.password)
		if valid != tt.valid {
			t.Errorf("ValidatePasswordStrong(%q) valid = %v, want %v", tt.password, valid, tt.valid)
		}
		if errKey != tt.errKey {
			t.Errorf("ValidatePasswordStrong(%q) errKey = %q, want %q", tt.password, errKey, tt.errKey)
		}
	}
}
