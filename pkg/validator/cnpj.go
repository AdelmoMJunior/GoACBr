package validator

import (
	"regexp"
	"strconv"
)

// IsValidCNPJ validates a Brazilian CNPJ using the Modulo 11 algorithm.
func IsValidCNPJ(cnpj string) bool {
	// Remove non-digit characters
	re := regexp.MustCompile(`\D`)
	cnpj = re.ReplaceAllString(cnpj, "")

	// CNPJ must have exactly 14 digits
	if len(cnpj) != 14 {
		return false
	}

	// Check for known invalid CNPJs (all digits same)
	invalid := true
	for i := 1; i < 14; i++ {
		if cnpj[i] != cnpj[0] {
			invalid = false
			break
		}
	}
	if invalid {
		return false
	}

	// Helper function to calculate check digit
	calcCheckDigit := func(cnpj string, weights []int) int {
		sum := 0
		for i, weight := range weights {
			digit, _ := strconv.Atoi(string(cnpj[i]))
			sum += digit * weight
		}
		rem := sum % 11
		if rem < 2 {
			return 0
		}
		return 11 - rem
	}

	// Calculate first check digit
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	digit1 := calcCheckDigit(cnpj[:12], weights1)

	// Calculate second check digit
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	digit2 := calcCheckDigit(cnpj[:13], weights2)

	// Verify digits
	expected1, _ := strconv.Atoi(string(cnpj[12]))
	expected2, _ := strconv.Atoi(string(cnpj[13]))

	return digit1 == expected1 && digit2 == expected2
}
