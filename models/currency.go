package models

type Currency struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"fullName"`
	Sign     string `json:"sign"`
}

func GetAllCurrencies() ([]Currency, error) {
	rows, err := DB.Query("SELECT * FROM Currencies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []Currency

	for rows.Next() {
		var currency Currency
		if err = rows.Scan(&currency.ID, &currency.Code, &currency.FullName, &currency.Sign); err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func InsertCurrency(currency Currency) (Currency, error) {
	result, err := DB.Exec(
		"INSERT INTO Currencies (Code, FullName, Sign) VALUES (?, ?, ?)",
		currency.Code, currency.FullName, currency.Sign,
	)
	if err != nil {
		return currency, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return currency, err
	}
	
	currency.ID = int(id)
	return currency, nil
}

func GetCurrencyByCode(code string) (Currency, error) {
	var currency Currency

	query := "SELECT Id, Code, FullName, Sign FROM Currencies WHERE Code = ?"
	err := DB.QueryRow(query, code).Scan(&currency.ID, &currency.Code, &currency.FullName, &currency.Sign)
	if err != nil {
		return currency, err
	}
	return currency, nil
}