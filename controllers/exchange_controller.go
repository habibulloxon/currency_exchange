package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/habibulloxon/currency_exchange/models"
)

func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	if from == "" || to == "" || amountStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "provide required fields"})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid data for amount"})
		return
	}

	exchangeRate, err := models.GetExchangeRateForPair(from, to)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "currency not found"})
		return
	}

	convertedAmount := amount * exchangeRate.Rate

	response := struct {
		BaseCurrency    models.Currency `json:"baseCurrency"`
		TargetCurrency  models.Currency `json:"targetCurrency"`
		Rate            float64         `json:"rate"`
		Amount          float64         `json:"amount"`
		ConvertedAmount float64         `json:"convertedAmount"`
	}{
		BaseCurrency:    exchangeRate.BaseCurrency,
		TargetCurrency:  exchangeRate.TargetCurrency,
		Rate:            exchangeRate.Rate,
		Amount:          amount,
		ConvertedAmount: convertedAmount,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
