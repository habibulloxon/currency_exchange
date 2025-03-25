package models

import (
	"database/sql"
	"errors"
)

type ExchangeRate struct {
	ID             int      `json:"id"`
	BaseCurrency   Currency `json:"baseCurrencyId"`
	TargetCurrency Currency `json:"targetCurrencyId"`
	Rate           float64  `json:"rate"`
}

func GetAllExchangeRates() ([]ExchangeRate, error) {
	query := `
	SELECT 
		er.ID,
		bc.ID, bc.FullName, bc.Code, bc.Sign,
		tc.ID, tc.FullName, tc.Code, tc.Sign,
		er.Rate
	FROM ExchangeRates er
	JOIN Currencies bc ON er.BaseCurrencyId = bc.ID
	JOIN Currencies tc ON er.TargetCurrencyId = tc.ID;
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []ExchangeRate

	for rows.Next() {
		var rate ExchangeRate
		err = rows.Scan(
			&rate.ID,
			&rate.BaseCurrency.ID, &rate.BaseCurrency.FullName, &rate.BaseCurrency.Code, &rate.BaseCurrency.Sign,
			&rate.TargetCurrency.ID, &rate.TargetCurrency.FullName, &rate.TargetCurrency.Code, &rate.TargetCurrency.Sign,
			&rate.Rate,
		)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}
	return rates, nil
}

func InsertExchangeRate(base Currency, target Currency, rate float64) (ExchangeRate, error) {
	var newExchangeRate ExchangeRate

	result, err := DB.Exec(
		"INSERT INTO ExchangeRates (BaseCurrencyId, TargetCurrencyId, Rate) VALUES (?, ?, ?)",
		base.ID, target.ID, rate,
	)
	if err != nil {
		return newExchangeRate, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return newExchangeRate, err
	}

	newExchangeRate.ID = int(id)
	newExchangeRate.BaseCurrency = base
	newExchangeRate.TargetCurrency = target
	newExchangeRate.Rate = rate

	return newExchangeRate, nil
}

func GetExchangeRateByPair(baseCode, targetCode string) (ExchangeRate, error) {
	var exchangeRate ExchangeRate

	query := `
		SELECT 
			ExchangeRateTable.ID,
			BaseCurrencyTable.ID, BaseCurrencyTable.FullName, BaseCurrencyTable.Code, BaseCurrencyTable.Sign,
			TargetCurrencyTable.ID, TargetCurrencyTable.FullName, TargetCurrencyTable.Code, TargetCurrencyTable.Sign,
			ExchangeRateTable.Rate
		FROM ExchangeRates AS ExchangeRateTable
		JOIN Currencies AS BaseCurrencyTable ON ExchangeRateTable.BaseCurrencyId = BaseCurrencyTable.ID
		JOIN Currencies AS TargetCurrencyTable ON ExchangeRateTable.TargetCurrencyId = TargetCurrencyTable.ID
		WHERE BaseCurrencyTable.Code = ? AND TargetCurrencyTable.Code = ?;
	`

	err := DB.QueryRow(query, baseCode, targetCode).Scan(
		&exchangeRate.ID,
		&exchangeRate.BaseCurrency.ID,
		&exchangeRate.BaseCurrency.FullName,
		&exchangeRate.BaseCurrency.Code,
		&exchangeRate.BaseCurrency.Sign,
		&exchangeRate.TargetCurrency.ID,
		&exchangeRate.TargetCurrency.FullName,
		&exchangeRate.TargetCurrency.Code,
		&exchangeRate.TargetCurrency.Sign,
		&exchangeRate.Rate,
	)
	if err != nil {
		return exchangeRate, err
	}

	return exchangeRate, nil
}

func UpdateExchangeRateByPair(baseCode, targetCode string, newRate float64) (ExchangeRate, error) {
	var updatedRate ExchangeRate

	baseCurrency, err := GetCurrencyByCode(baseCode)
	if err != nil {
		return updatedRate, err
	}

	targetCurrency, err := GetCurrencyByCode(targetCode)
	if err != nil {
		return updatedRate, err
	}

	queryGet := "SELECT ID FROM ExchangeRates WHERE BaseCurrencyId = ? AND TargetCurrencyId = ?"
	var exchangeRateID int
	err = DB.QueryRow(queryGet, baseCurrency.ID, targetCurrency.ID).Scan(&exchangeRateID)
	if err != nil {
		return updatedRate, err
	}

	queryUpdate := "UPDATE ExchangeRates SET Rate = ? WHERE ID = ?"
	result, err := DB.Exec(queryUpdate, newRate, exchangeRateID)
	if err != nil {
		return updatedRate, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return updatedRate, sql.ErrNoRows
	}

	updatedRate, err = GetExchangeRateByPair(baseCode, targetCode)
	if err != nil {
		return updatedRate, err
	}
	return updatedRate, nil
}

func GetExchangeRateForPair(fromCode, toCode string) (ExchangeRate, error) {
	direct, err := GetExchangeRateByPair(fromCode, toCode)
	if err == nil {
		return direct, nil
	}

	reversed, err2 := GetExchangeRateByPair(toCode, fromCode)
	if err2 == nil {
		if reversed.Rate == 0 {
			return ExchangeRate{}, errors.New("incorrect coefficient of exchange")
		}
		var result ExchangeRate
		result.ID = reversed.ID
		result.BaseCurrency = reversed.TargetCurrency
		result.TargetCurrency = reversed.BaseCurrency
		result.Rate = 1 / reversed.Rate
		return result, nil
	}

	usdToFrom, err3 := GetExchangeRateByPair("USD", fromCode)
	usdToTo, err4 := GetExchangeRateByPair("USD", toCode)
	if err3 == nil && err4 == nil {
		if usdToFrom.Rate == 0 {
			return ExchangeRate{}, errors.New("incorrect coefficient of exchange")
		}
		var result ExchangeRate
		result.Rate = usdToTo.Rate / usdToFrom.Rate
		result.BaseCurrency = usdToFrom.TargetCurrency
		result.TargetCurrency = usdToTo.TargetCurrency
		result.ID = 0
		return result, nil
	}

	return ExchangeRate{}, sql.ErrNoRows
}