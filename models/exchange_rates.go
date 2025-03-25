package models

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
