package main

import (
	"log"
	"net/http"

	"github.com/habibulloxon/currency_exchange/controllers"
	"github.com/habibulloxon/currency_exchange/models"
)

func main() {
	if err := models.OpenExistingDB(); err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer models.DB.Close()

	http.HandleFunc("/currencies", controllers.CurrenciesHandler)
	http.HandleFunc("/currencies/", controllers.CurrencyCodeHandler)

	http.HandleFunc("/exchange", controllers.ExchangeHandler)
	http.HandleFunc("/exchangeRates", controllers.ExchangeRatesHandler)
	http.HandleFunc("/exchangeRate/", controllers.ExchangeRatePairHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("salam aleykum"))
	})

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
