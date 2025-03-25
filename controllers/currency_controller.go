package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/habibulloxon/currency_exchange/models"
)

func CurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		currencies, err := models.GetAllCurrencies()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(currencies)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error parsing form data",
			})
			return
		}

		name := strings.TrimSpace(r.FormValue("fullName"))
		code := strings.TrimSpace(r.FormValue("code"))
		sign := strings.TrimSpace(r.FormValue("sign"))

		if name == "" || code == "" || sign == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Missing required fields. Provide FullName, Code, and Sign.",
			})
			return
		}

		_, err := models.GetCurrencyByCode(code)
		if err == nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Currency with that code already exists",
			})
			return
		} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		newCurrency := models.Currency{
			FullName: name,
			Code: code,
			Sign: sign,
		}

		insertedCurrency, err := models.InsertCurrency(newCurrency)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(insertedCurrency)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func CurrencyCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	codeStr := strings.TrimPrefix(r.URL.Path, "/currencies/")
	codeStr = strings.TrimSpace(codeStr)
	if codeStr == "" || codeStr == "/" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Currency code is missing in the URL",
		})
		return
	}

	currency, err := models.GetCurrencyByCode(codeStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Currency not found",
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
	json.NewEncoder(w).Encode(currency)
}