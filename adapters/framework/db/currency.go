package db

import (
	"fmt"

	"github.com/deestarks/infiniti/utils"
	"github.com/lib/pq"
)

type (
	CurrencyModel struct {
		Id					int		`json:"id"`
		Name				string	`json:"name"`
		Symbol				string	`json:"symbol"`
		ConversionRateToUSD	float64	`json:"conversion_rate_to_usd"`
	}
	
	CurrencyAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewCurrencyAdapter() *CurrencyAdapter {
	return &CurrencyAdapter{
		adapter: adpt,
		tableName: "currencies",
	}
}

func (mAdapt *CurrencyAdapter) Create(data map[string]interface{}) (*CurrencyModel, error) {
	var currency CurrencyModel

	mToS := utils.MapToStructSlice(data)
	var (
		colStr	string
		valArr	[]interface{}
	)
	for i, s := range mToS {
		colStr += s.Key + ", "
		valArr = append(valArr, s.Value)
		if i == len(mToS)-1 {
			colStr = colStr[:len(colStr)-2] // remove the last ", "
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %s ( %s ) VALUES ( %s )
		RETURNING id, name, symbol, conversion_rate_to_usd
	`, mAdapt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&currency.Id, &currency.Name, &currency.Symbol, &currency.ConversionRateToUSD)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &currency, nil
}

func (mAdapt *CurrencyAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*CurrencyModel, error) {
	var (
		currency 	CurrencyModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, name, symbol, conversion_rate_to_usd
	`, mAdapt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1)

	valArr = append(valArr, colValue)

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&currency.Id, &currency.Name, &currency.Symbol, &currency.ConversionRateToUSD)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &currency, nil
}

func (mAdapt *CurrencyAdapter) Delete(colName string, value interface{}) (*CurrencyModel, error) {
	var (
		currency	CurrencyModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, name, symbol, conversion_rate_to_usd
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&currency.Id, &currency.Name, &currency.Symbol, &currency.ConversionRateToUSD)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &currency, nil
}

func (mAdapt *CurrencyAdapter) Get(colName string, value interface{}) (*CurrencyModel, error) {
	var (
		currency	CurrencyModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT id, name, symbol, conversion_rate_to_usd
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&currency.Id, &currency.Name, &currency.Symbol, &currency.ConversionRateToUSD)
    if err != nil {
		return nil, err
    }
	return &currency, nil
}

func (mAdapt *CurrencyAdapter) Filter(colName string, value interface{}) (*[]CurrencyModel, error) {
	var (
		currencies	[]CurrencyModel
		err			error
	)
	query := fmt.Sprintf(`
		SELECT id, name, symbol, conversion_rate_to_usd
		FROM %s WHERE %s = $1
	`, mAdapt.tableName, colName)
	rows, err := mAdapt.adapter.db.Query(query, value)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var currency CurrencyModel
		err := rows.Scan(&currency.Id, &currency.Name, &currency.Symbol, &currency.ConversionRateToUSD)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
		currencies = append(currencies, currency)
	}
	return &currencies, nil
}

func (mAdapt *CurrencyAdapter) List() (*[]CurrencyModel, error) {
	var (
		currencies	[]CurrencyModel
		err			error
	)
	query := fmt.Sprintf(`
		SELECT id, name, symbol, conversion_rate_to_usd
		FROM %s`, 
	mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var currency CurrencyModel
		err := rows.Scan(&currency.Id, &currency.Name, &currency.Symbol, &currency.ConversionRateToUSD)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
		currencies = append(currencies, currency)
	}
	return &currencies, nil
}