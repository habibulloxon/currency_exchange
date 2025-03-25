package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/habibulloxon/currency_exchange/models"
	"net/http"
	"strings"
	"strconv"
)

func ExchangeRatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		rates, err := models.GetAllExchangeRates()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rates)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error parsing form data",
			})
			return
		}

		baseCurrencyCode := strings.TrimSpace(r.FormValue("baseCurrencyCode"))
		targetCurrencyCode := strings.TrimSpace(r.FormValue("targetCurrencyCode"))
		rateStr := strings.TrimSpace(r.FormValue("rate"))

		if baseCurrencyCode == "" || targetCurrencyCode == "" || rateStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Missing required fields. Provide baseCurrencyCode, targetCurrencyCode, and rate.",
			})
			return
		}

		rateValue, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid rate value",
			})
			return
		}

		_, err = models.GetExchangeRateByPair(baseCurrencyCode, targetCurrencyCode)
		if err == nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Exchange rate for the given currency pair already exists",
			})
			return
		} else if err != nil && err != sql.ErrNoRows {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		baseCurrency, err := models.GetCurrencyByCode(baseCurrencyCode)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Base currency not found",
				})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
			}
			return
		}

		targetCurrency, err := models.GetCurrencyByCode(targetCurrencyCode)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Target currency not found",
				})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
			}
			return
		}

		newExchangeRate, err := models.InsertExchangeRate(baseCurrency, targetCurrency, rateValue)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newExchangeRate)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func ExchangeRatePairHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pairStr := strings.TrimPrefix(r.URL.Path, "/exchangeRate/")
	pairStr = strings.TrimSpace(pairStr)

	if len(pairStr) != 6 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Currency pair is missing or invalid in the URL",
		})
		return
	}

	baseCode := pairStr[:3]
	targetCode := pairStr[3:]

	exRate, err := models.GetExchangeRateByPair(baseCode, targetCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Exchange rate for the given pair not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exRate)
}
