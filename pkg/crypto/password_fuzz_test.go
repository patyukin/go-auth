package crypto

import (
	"testing"
)

func ParseComplex(data []byte) bool {
	if len(data) == 5 {
		if data[0] == 'F' &&
			data[1] == 'U' &&
			data[2] == 'Z' &&
			data[3] == 'Z' &&
			data[4] == 'I' &&
			data[5] == 'T' {
			return true
		}
	}
	return false
}

func FuzzSanitizeAndCheckUsername(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		ParseComplex(data)
	})
}

// func isValidPassword(password string) bool {
// 	if len(password) < 8 {
// 		return false
// 	}
// 	var hasUpper, hasLower, hasNumber, hasSpecial bool
// 	for _, r := range password {
// 		switch {
// 		case unicode.IsUpper(r):
// 			hasUpper = true
// 		case unicode.IsLower(r):
// 			hasLower = true
// 		case unicode.IsNumber(r):
// 			hasNumber = true
// 		case unicode.IsPunct(r) || unicode.IsSymbol(r):
// 			hasSpecial = true
// 		}
// 	}
// 	return hasUpper && hasLower && hasNumber && hasSpecial
// }

// func FuzzIsValidPassword(f *testing.F) {
// 	seedCorpus := []string{
// 		"Password1!",
// 		"short",
// 		"nouppercase1!",
// 		"NOLOWERCASE1!",
// 		"NoNumber!",
// 		"NoSpecialChar1",
// 		"",
// 	}

// 	for _, seed := range seedCorpus {
// 		f.Add(seed)
// 	}

// 	f.Fuzz(func(t *testing.T, password string) {
// 		isValidPassword(password)
// 	})
// }

// // func sanitizeAndCheckUsername(username string) bool {
// // 	username = strings.TrimSpace(username)

// // 	for _, char := range username {
// // 		if !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9') {
// // 			return false
// // 		}
// // 	}

// // 	if len(username) < 1 || len(username) > 20 {
// // 		return false
// // 	}

// // 	return true
// // }

// // func FuzzSanitizeAndCheckUsername(f *testing.F) {
// // 	testCases := []string{
// // 		"validUsername",
// // 		"    leadingSpace",
// // 		"trailingSpace    ",
// // 		" middle Space ",
// // 		"Valid123",
// // 		"user@name!",
// // 		"user\nname",
// // 		"_____",
// // 		"",
// // 		"abcdefghijklmnopqrstuv", // 21 characters
// // 		"你好",
// // 		"😊",
// // 		"👩‍💻",
// // 		string([]byte{0xC2, 0xA0}), // Non-breaking space
// // 		"\u200B",                   // Zero-width space
// // 	}

// // 	for _, tc := range testCases {
// // 		f.Add(tc)
// // 	}

// // 	f.Fuzz(func(t *testing.T, username string) {
// // 		result := sanitizeAndCheckUsername(username)
// // 		if len(username) == 0 && result {
// // 			t.Errorf("sanitizeAndCheckUsername(%q) = %v; want false", username, result)
// // 		}
// // 	})
// // }

// // func CheckPasswordStrength(password string) bool {
// // 	if len(password) < 8 {
// // 		return false
// // 	}

// // 	hasUpper := false
// // 	hasLower := false
// // 	hasNumber := false
// // 	hasSpecial := false

// // 	for _, char := range password {
// // 		switch {
// // 		case unicode.IsUpper(char):
// // 			hasUpper = true
// // 		case unicode.IsLower(char):
// // 			hasLower = true
// // 		case unicode.IsNumber(char):
// // 			hasNumber = true
// // 		case unicode.IsPunct(char) || unicode.IsSymbol(char):
// // 			hasSpecial = true
// // 		}
// // 		// if hasUpper && hasLower && hasNumber && hasSpecial {
// // 		// 	return true
// // 		// }
// // 	}

// // 	return hasUpper && hasLower && hasNumber && hasSpecial
// // }

// // // Fuzz function for password strength checker
// // func FuzzCheckPasswordStrength(f *testing.F) {
// // 	// Adding basic seed cases
// // 	// seeds := []string{"password123!", "Admin123@", "abcDEF123!", "weakPass", "12345678", "!@#AnB1"}

// // 	// for _, seed := range seeds {
// // 	// 	f.Add(seed)
// // 	// }

// // 	f.Fuzz(func(t *testing.T, password string) {
// // 		strength := CheckPasswordStrength(password)
// // 		if len(password) > 8 && !strength {
// // 			t.Errorf("Expected strong password but got weak: %v", password)
// // 		}
// // 		if len(password) < 8 && strength {
// // 			t.Errorf("Expected weak password but got strong: %v", password)
// // 		}
// // 	})
// // }

// // func RegisterUser(username, password string) error {
// // 	if strings.TrimSpace(username) == "" {
// // 		return errors.New("username is required")
// // 	}
// // 	if strings.TrimSpace(password) == "" {
// // 		return errors.New("password is required")
// // 	}
// // 	if len(password) < 8 {
// // 		return errors.New("password must be at least 8 characters long")
// // 	}

// // 	isAlphaNum := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
// // 	if !isAlphaNum(username) {
// // 		return errors.New("username must be alphanumeric")
// // 	}

// // 	commonPasswords := []string{"password", "12345678", "qwerty", "abc123"}
// // 	for _, commonPassword := range commonPasswords {
// // 		if password == commonPassword {
// // 			return errors.New("password is too common")
// // 		}
// // 	}

// // 	return nil

// // }

// // // Fuzz function for user registration
// // func FuzzRegisterUser(f *testing.F) {
// // 	// Adding basic seed cases
// // 	seeds := []struct {
// // 		username string
// // 		password string
// // 	}{
// // 		{"user1", "password123"},
// // 		{"", "password123"},
// // 		{"user2", ""},
// // 		{"user3", "short"},
// // 		{"user4", "sp@ci@lCh@r"},
// // 	}

// // 	for _, seed := range seeds {
// // 		f.Add(seed.username, seed.password)
// // 	}

// // 	f.Fuzz(func(t *testing.T, username, password string) {
// // 		err := RegisterUser(username, password)

// // 		if strings.TrimSpace(username) == "" && err == nil {
// // 			t.Errorf("Expected an error for empty username, got nil")
// // 		}

// // 		if strings.TrimSpace(password) == "" && err == nil {
// // 			t.Errorf("Expected an error for empty password, got nil")
// // 		}

// // 		if len(password) < 8 && err == nil {
// // 			t.Errorf("Expected an error for a short password, got nil")
// // 		}
// // 	})
// // }
