package services

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// luhnCheck validates a credit card number using the Luhn algorithm.
func LuhnCheck(cardNumber string) bool {
	sum := 0
	alternate := false

	// Iterate over the card number from right to left
	for i := len(cardNumber) - 1; i >= 0; i-- {
		n := int(cardNumber[i] - '0')

		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}

// getCardVendor identifies the card vendor based on the IIN/BIN.
func GetCardVendor(cardNumber string) string {
	if len(cardNumber) < 4 {
		return "Unknown"
	}

	switch {
	case strings.HasPrefix(cardNumber, "4"):
		return "Visa"
	case strings.HasPrefix(cardNumber, "51"), strings.HasPrefix(cardNumber, "52"),
		strings.HasPrefix(cardNumber, "53"), strings.HasPrefix(cardNumber, "54"),
		strings.HasPrefix(cardNumber, "55"):
		return "MasterCard"
	case strings.HasPrefix(cardNumber, "34"), strings.HasPrefix(cardNumber, "37"):
		return "American Express"
	case strings.HasPrefix(cardNumber, "6011"), strings.HasPrefix(cardNumber, "65"),
		strings.HasPrefix(cardNumber, "644"), strings.HasPrefix(cardNumber, "645"),
		strings.HasPrefix(cardNumber, "646"), strings.HasPrefix(cardNumber, "647"),
		strings.HasPrefix(cardNumber, "648"), strings.HasPrefix(cardNumber, "649"):
		return "Discover"
	case strings.HasPrefix(cardNumber, "36"), strings.HasPrefix(cardNumber, "38"),
		strings.HasPrefix(cardNumber, "300"), strings.HasPrefix(cardNumber, "301"),
		strings.HasPrefix(cardNumber, "302"), strings.HasPrefix(cardNumber, "303"),
		strings.HasPrefix(cardNumber, "304"), strings.HasPrefix(cardNumber, "305"):
		return "Diners Club"
	case strings.HasPrefix(cardNumber, "3528"), strings.HasPrefix(cardNumber, "3589"):
		return "JCB"
	case strings.HasPrefix(cardNumber, "50"), strings.HasPrefix(cardNumber, "56"),
		strings.HasPrefix(cardNumber, "57"), strings.HasPrefix(cardNumber, "58"),
		strings.HasPrefix(cardNumber, "59"), strings.HasPrefix(cardNumber, "60"),
		strings.HasPrefix(cardNumber, "62"), strings.HasPrefix(cardNumber, "63"),
		strings.HasPrefix(cardNumber, "64"), strings.HasPrefix(cardNumber, "67"):
		return "Maestro"
	default:
		return "Unknown"
	}
}

// generateCardNumber generates a random valid credit card number.
func GenerateCardNumber() string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random 15-digit number
	num := rand.Intn(900000000000000) + 100000000000000

	// Convert the number to a string
	cardNumber := strconv.Itoa(num)

	// Add a random check digit using the Luhn algorithm
	checkDigit := GenerateCheckDigit(cardNumber)
	cardNumber += strconv.Itoa(checkDigit)

	// Format the card number with spaces
	formattedCardNumber := FormatCardNumber(cardNumber)

	return formattedCardNumber
}

// generateCheckDigit generates a random check digit using the Luhn algorithm.
func GenerateCheckDigit(cardNumber string) int {
	sum := 0
	alternate := false

	// Iterate over the card number from right to left
	for i := len(cardNumber) - 1; i >= 0; i-- {
		n := int(cardNumber[i] - '0')

		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alternate = !alternate
	}

	checkDigit := (sum * 9) % 10

	return checkDigit
}

// formatCardNumber formats the card number with spaces every 4 digits.
func FormatCardNumber(cardNumber string) string {
	var formatted strings.Builder
	for i, r := range cardNumber {
		formatted.WriteRune(r)
		if (i+1)%4 == 0 && i != len(cardNumber)-1 {
			formatted.WriteRune(' ')
		}
	}
	return formatted.String()
}

// cleanCardNumber removes any spaces or non-numeric characters from the card number.
func CleanCardNumber(cardNumber string) string {
	var cleaned strings.Builder
	for _, r := range cardNumber {
		if unicode.IsDigit(r) {
			cleaned.WriteRune(r)
		}
	}
	return cleaned.String()
}
