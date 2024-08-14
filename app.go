package main

import "card-verification/api/internal/transport/rest"

func main() {
	println("Starting Card Verification Service!")
	rest.HandleRequests()
}
