package rest

import (
	"card-verification/api/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) any {
	return map[string]string{
		"message": "Welcome to the Card Verification Service!\n" +
			"Please provide a credit card number to verify.\n",
		"example": "http://localhost:8888/verify?card=4539%201488%200343%206467\n",
	}
}

type Response struct {
	Valid  bool   `json:"valid"`
	Vendor string `json:"vendor"`
}

func verifyPage(w http.ResponseWriter, r *http.Request) any {
	cardNumber := r.URL.Query().Get("card")
	cleanedCardNumber := services.CleanCardNumber(cardNumber)

	var response Response

	if len(cleanedCardNumber) < 13 || len(cleanedCardNumber) > 19 {
		response.Valid = false
		return response
	}

	if !services.LuhnCheck(cleanedCardNumber) {
		response.Valid = false
	} else {
		response.Valid = true
		response.Vendor = services.GetCardVendor(cleanedCardNumber)
	}

	return response
}

type jsonRequestFunc func(http.ResponseWriter, *http.Request) any

func jsonResponse(fn jsonRequestFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonResponse, err := json.Marshal(fn(w, r))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}

}

func HandleRequests() {
	serverUrl := ""
	serverPort := "8888"
	if serverUrl == "" {
		serverUrl = fmt.Sprintf("%s:%s", serverUrl, serverPort)
	}

	println("Server is started on " + serverUrl)
	http.HandleFunc("/", jsonResponse(homePage))
	http.HandleFunc("/verify", jsonResponse(verifyPage))
	log.Fatal(http.ListenAndServe(serverUrl, nil))
}
