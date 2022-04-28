package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)

type Currency struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

type CurrencyResource struct {
	Id		int		`json:"id"`
	Name	string	`json:"name"`
	Symbol 	string	`json:"symbol"`
}

func (service *Service) NewCurrencyService() *Currency {
	return &Currency{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

func (curr *Currency) GetCurrency(key string, value interface{}) (CurrencyResource, error) {
	dbAdapter := curr.dbPort.NewCurrencyAdapter()
	currRes, err := dbAdapter.Get(key, value)
	if err != nil {
		return CurrencyResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf("currency does not exist"),
		}
	}

	// Serialization and return
	var res CurrencyResource
	jsonCurr, _ := json.Marshal(currRes)
	json.Unmarshal(jsonCurr, &res)
	return res, nil
}

func (curr *Currency) CreateCurrency(data map[string]interface{}) (CurrencyResource, error) {
	rsrts := []string{"id"} // Restricted columns to create
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range rsrts {
		delete(data, field)
	}

	requires := []string{"name", "symbol"} // Required fields
	var notFound string
	for _, field := range requires {
		if _, ok := data[field]; !ok {
			notFound += " ,"+field
		}
	}
	if notFound != "" {
		return CurrencyResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("missing fields: %s", notFound[2:]),
		}
	}
	
	dbAdapter := curr.dbPort.NewCurrencyAdapter()
	currRes, err := dbAdapter.Create(data)
	if err != nil {
		return CurrencyResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf(err.Error()),
		}
	}

	// Serialization and return
	var res CurrencyResource
	jsonCurr, _ := json.Marshal(currRes)
	json.Unmarshal(jsonCurr, &res)
	return res, nil
}

func (curr *Currency) UpdateCurrency(key string, value interface{}, data map[string]interface{}) (CurrencyResource, error) {
	rsrts := []string{"id"} // Restricted columns to update
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range rsrts {
		delete(data, field)
	}
	
	dbAdapter := curr.dbPort.NewCurrencyAdapter()
	// Check if exists
	_, err := dbAdapter.Get(key, value)
	if err != nil {
		return CurrencyResource{}, &utils.RequestError{
			Code: http.StatusNotFound,
			Err: fmt.Errorf("could not update - currency does not exist"),
		}
	}

	currRes, err := dbAdapter.Update(key, value, data)
	if err != nil {
		return CurrencyResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf(err.Error()),
		}
	}

	// Serialization and return
	var res CurrencyResource
	jsonCurr, _ := json.Marshal(currRes)
	json.Unmarshal(jsonCurr, &res)
	return res, nil
}

func (curr *Currency) ListCurrencies() ([]CurrencyResource, error) {
	dbAdapter := curr.dbPort.NewCurrencyAdapter()
	currs, err := dbAdapter.List()
	if err != nil {
		return []CurrencyResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	fmt.Errorf("error retrieving currencies"),
		}
	}

	// Serialization and return
	var res []CurrencyResource
	jsonCurrs, _ := json.Marshal(currs)
	json.Unmarshal(jsonCurrs, &res)
	return res, nil
}

func (curr *Currency) DeleteCurrency(key string, value interface{}) error {
	adapter := curr.dbPort.NewCurrencyAdapter()
	// Check if exists
	_, err := adapter.Get(key, value)
	if err != nil {
		return &utils.RequestError{
			Code: http.StatusNotFound,
			Err: fmt.Errorf("could not delete - currency does not exist"),
		}
	}
	
	_, err = adapter.Delete(key, value)
	if err != nil {
		return &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: err,
		}
	}
	return nil
}