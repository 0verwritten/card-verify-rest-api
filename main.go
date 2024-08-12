package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"
)

// luhnCheck validates a credit card number using the Luhn algorithm.
func luhnCheck(cardNumber string) bool {
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
func getCardVendor(cardNumber string) string {
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

// cleanCardNumber removes any spaces or non-numeric characters from the card number.
func cleanCardNumber(cardNumber string) string {
	var cleaned strings.Builder
	for _, r := range cardNumber {
		if unicode.IsDigit(r) {
			cleaned.WriteRune(r)
		}
	}
	return cleaned.String()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, this is the Card Verification Service!\n")
	fmt.Fprintf(w, "Please provide a credit card number to verify.\n")
	fmt.Fprintf(w, "Example: http://localhost:8888/verify?card=4539 1488 0343 6467\n")
}

func verifyPage(w http.ResponseWriter, r *http.Request) {
	cardNumber := r.URL.Query().Get("card")
	cleanedCardNumber := cleanCardNumber(cardNumber)

	type Response struct {
		Valid  bool   `json:"valid"`
		Vendor string `json:"vendor"`
	}

	var response Response

	if !luhnCheck(cleanedCardNumber) {
		response.Valid = false
	} else {
		response.Valid = true
		response.Vendor = getCardVendor(cleanedCardNumber)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func handleRequests() {
	serverUrl := ""
	serverPort := "8888"
	if serverUrl == "" {
		serverUrl = fmt.Sprintf("%s:%s", serverUrl, serverPort)
	}

	println("Server is started on " + serverUrl)
	http.HandleFunc("/", homePage)
	http.HandleFunc("/verify", verifyPage)
	log.Fatal(http.ListenAndServe(serverUrl, nil))
}

func main() {
	println("Starting Card Verification Service!")
	handleRequests()
}
