package validator

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidateCNPJ validates a Brazilian CNPJ number.
// Accepts both formatted (XX.XXX.XXX/XXXX-XX) and unformatted (14 digits).
func ValidateCNPJ(cnpj string) error {
	// Remove formatting characters.
	cnpj = strings.ReplaceAll(cnpj, ".", "")
	cnpj = strings.ReplaceAll(cnpj, "/", "")
	cnpj = strings.ReplaceAll(cnpj, "-", "")
	cnpj = strings.TrimSpace(cnpj)

	if len(cnpj) != 14 {
		return fmt.Errorf("CNPJ must have 14 digits, got %d", len(cnpj))
	}

	// Check if all digits are the same (e.g., 00000000000000).
	allSame := true
	for i := 1; i < len(cnpj); i++ {
		if cnpj[i] != cnpj[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return fmt.Errorf("invalid CNPJ: all digits are the same")
	}

	// Validate check digits.
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	digits := make([]int, 14)
	for i, c := range cnpj {
		d, err := strconv.Atoi(string(c))
		if err != nil {
			return fmt.Errorf("CNPJ contains non-numeric character at position %d", i)
		}
		digits[i] = d
	}

	// First check digit.
	sum := 0
	for i, w := range weights1 {
		sum += digits[i] * w
	}
	remainder := sum % 11
	firstCheck := 0
	if remainder >= 2 {
		firstCheck = 11 - remainder
	}
	if digits[12] != firstCheck {
		return fmt.Errorf("invalid CNPJ: first check digit mismatch")
	}

	// Second check digit.
	sum = 0
	for i, w := range weights2 {
		sum += digits[i] * w
	}
	remainder = sum % 11
	secondCheck := 0
	if remainder >= 2 {
		secondCheck = 11 - remainder
	}
	if digits[13] != secondCheck {
		return fmt.Errorf("invalid CNPJ: second check digit mismatch")
	}

	return nil
}

// SanitizeCNPJ removes formatting and returns only digits.
func SanitizeCNPJ(cnpj string) string {
	cnpj = strings.ReplaceAll(cnpj, ".", "")
	cnpj = strings.ReplaceAll(cnpj, "/", "")
	cnpj = strings.ReplaceAll(cnpj, "-", "")
	return strings.TrimSpace(cnpj)
}

// FormatCNPJ formats a 14-digit CNPJ string as XX.XXX.XXX/XXXX-XX.
func FormatCNPJ(cnpj string) string {
	cnpj = SanitizeCNPJ(cnpj)
	if len(cnpj) != 14 {
		return cnpj
	}
	return fmt.Sprintf("%s.%s.%s/%s-%s", cnpj[0:2], cnpj[2:5], cnpj[5:8], cnpj[8:12], cnpj[12:14])
}
