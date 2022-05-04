package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)

type TransactionTypes struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

type TransactionTypesResource struct {
	Id		int		`json:"id"`
	Name	string	`json:"name"`
}

func (service *Service) NewTransactionTypesService() *TransactionTypes {
	return &TransactionTypes{
		dbPort: 	service.dbPort,
		corePort:	service.corePort,
	}
}

func (tt *TransactionTypes) GetTransactionType(colName string, colValue interface{}) (TransactionTypesResource, error) {
	dbAdapter := tt.dbPort.NewTransactionTypeAdapter()
	ttRes, err := dbAdapter.Get(colName, colValue)
	if err != nil {
		return TransactionTypesResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf("requested transaction type not found"),
		}
	}

	// Serialization and return
	var res TransactionTypesResource
	jsonRes, _ := json.Marshal(ttRes)
	json.Unmarshal(jsonRes, &res)
	return res, nil
}

func (tt *TransactionTypes) CreateTransactionType(data map[string]interface{}) (TransactionTypesResource, error) {
	rsrts := []string{"id"} // Restricted columns to create
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range rsrts {
		delete(data, field)
	}

	requires := []string{"name"} // Required fields
	var notFound string
	for _, field := range requires {
		if _, ok := data[field]; !ok {
			notFound += " ,"+field
		}
	}
	if notFound != "" {
		return TransactionTypesResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("missing fields: %s", notFound[2:]),
		}
	}
	
	dbAdapter := tt.dbPort.NewTransactionTypeAdapter()
	ttRes, err := dbAdapter.Create(data)
	if err != nil {
		return TransactionTypesResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf(err.Error()),
		}
	}

	// Serialization and return
	var res TransactionTypesResource
	jsonRes, _ := json.Marshal(ttRes)
	json.Unmarshal(jsonRes, &res)
	return res, nil
}

func (tt *TransactionTypes) UpdateTransactionType(colName string, colValue interface{}, data map[string]interface{}) (TransactionTypesResource, error) {
	rsrts := []string{"id"} // Restricted columns to update
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range rsrts {
		delete(data, field)
	}
	
	dbAdapter := tt.dbPort.NewTransactionTypeAdapter()
	// Check if exists
	_, err := dbAdapter.Get(colName, colValue)
	if err != nil {
		return TransactionTypesResource{}, &utils.RequestError{
			Code: http.StatusNotFound,
			Err: fmt.Errorf("could not update - transaction type does not exist"),
		}
	}

	ttRes, err := dbAdapter.Update(colName, colValue, data)
	if err != nil {
		return TransactionTypesResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf(err.Error()),
		}
	}

	// Serialization and return
	var res TransactionTypesResource
	jsonRes, _ := json.Marshal(ttRes)
	json.Unmarshal(jsonRes, &res)
	return res, nil
}

func (tt *TransactionTypes) ListTransactionTypes() ([]TransactionTypesResource, error) {
	dbAdapter := tt.dbPort.NewTransactionTypeAdapter()
	tts, err := dbAdapter.List()
	if err != nil {
		return []TransactionTypesResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	fmt.Errorf("error retrieving transactions"),
		}
	}

	// Serialization and return
	var res []TransactionTypesResource
	jsonRes, _ := json.Marshal(tts)
	json.Unmarshal(jsonRes, &res)
	return res, nil
}

func (tt *TransactionTypes) DeleteTransactionType(colName string, colValue interface{}) error {
	adapter := tt.dbPort.NewTransactionTypeAdapter()
	// Check if exists
	transType, err := adapter.Get(colName, colValue)
	if err != nil {
		return &utils.RequestError{
			Code: http.StatusNotFound,
			Err: fmt.Errorf("could not delete - transaction type does not exist"),
		}
	}

	// Don't allow deletion of "deposit", "withdrawal" and "transfer"
	undeletables := map[string]bool{
		"deposit": true,
		"withdrawal": true,
		"transfer": true,
	}
	if undeletables[transType.Name] {
		return &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("transaction type of '%s' cannot be deleted", transType.Name),
		}
	}
	
	_, err = adapter.Delete(colName, colValue)
	if err != nil {
		return &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: err,
		}
	}
	return nil
}