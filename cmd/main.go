package main

import (
	"coinService/internal/api"
	"coinService/internal/auth"
	"coinService/internal/database"
	"log"
	"net/http"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	authService := auth.NewAuthService(db)
	merchService := api.NewMerchService(db)

	http.HandleFunc("/api/auth", authService.HandleAuth)
	http.HandleFunc("/api/info", merchService.HandleInfo)
	http.HandleFunc("/api/sendCoin", merchService.HandleSendCoin)
	http.HandleFunc("/api/buy/", merchService.HandleBuy)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
