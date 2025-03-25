package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/habibulloxon/currency_exchange/models"
	"net/http"
	"strconv"
	"strings"
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
		} else if err != sql.ErrNoRows {
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

	switch r.Method {
	case http.MethodGet:
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

	case http.MethodPatch:
		pairStr := strings.TrimPrefix(r.URL.Path, "/exchangeRate/")
		pairStr = strings.TrimSpace(pairStr)
		if len(pairStr) != 6 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Currency pair is missing or invalid in the URL",
			})
			return
		}
		baseCurrencyCode := pairStr[:3]
		targetCurrencyCode := pairStr[3:]

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error parsing form data",
			})
			return
		}

		rateStr := strings.TrimSpace(r.FormValue("rate"))
		if rateStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Missing required field: rate",
			})
			return
		}

		newRate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid rate value",
			})
			return
		}

		updatedExchangeRate, err := models.UpdateExchangeRateByPair(baseCurrencyCode, targetCurrencyCode, newRate)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Exchange rate for the given currency pair not found",
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
		json.NewEncoder(w).Encode(updatedExchangeRate)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
