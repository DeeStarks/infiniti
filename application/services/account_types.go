package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)

type AccountType struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

type AccountTypeResource struct {
	Id		int		`json:"id"`
	Name	string	`json:"name"`
}

func (service *Service) NewAccountTypeService() *AccountType {
	return &AccountType{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

func (acctT *AccountType) GetAccountType(key string, value interface{}) (AccountTypeResource, error) {
	dbAdapter := acctT.dbPort.NewAccountTypeAdapter()
	acctTRes, err := dbAdapter.Get(key, value)
	fmt.Println(acctTRes)
	if err != nil {
		return AccountTypeResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf("account type does not exist"),
		}
	}

	// Serialization and return
	var res AccountTypeResource
	jsonAcctT, _ := json.Marshal(acctTRes)
	json.Unmarshal(jsonAcctT, &res)
	return res, nil
}

func (acctT *AccountType) CreateAccountType(data map[string]interface{}) (AccountTypeResource, error) {
	rsrts := []string{"id"} // Restricted columns to update
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range rsrts {
		delete(data, field)
	}
	
	dbAdapter := acctT.dbPort.NewAccountTypeAdapter()
	acctTRes, err := dbAdapter.Create(data)
	if err != nil {
		return AccountTypeResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf(err.Error()),
		}
	}

	// Serialization and return
	var res AccountTypeResource
	jsonAcctType, _ := json.Marshal(acctTRes)
	json.Unmarshal(jsonAcctType, &res)
	return res, nil
}

func (acctT *AccountType) UpdateAccountType(key string, value interface{}, data map[string]interface{}) (AccountTypeResource, error) {
	rsrts := []string{"id"} // Restricted columns to update
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range rsrts {
		delete(data, field)
	}
	
	dbAdapter := acctT.dbPort.NewAccountTypeAdapter()
	acctTRes, err := dbAdapter.Update(key, value, data)
	if err != nil {
		return AccountTypeResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf(err.Error()),
		}
	}

	// Serialization and return
	var res AccountTypeResource
	jsonAcctType, _ := json.Marshal(acctTRes)
	json.Unmarshal(jsonAcctType, &res)
	return res, nil
}

func (acctT *AccountType) ListAccountTypes() ([]AccountTypeResource, error) {
	dbAdapter := acctT.dbPort.NewAccountTypeAdapter()
	acctTs, err := dbAdapter.List()
	if err != nil {
		return []AccountTypeResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	fmt.Errorf("error fetching account types"),
		}
	}

	// Serialization and return
	var res []AccountTypeResource
	jsonAcctTypes, _ := json.Marshal(acctTs)
	json.Unmarshal(jsonAcctTypes, &res)
	return res, nil
}

func (acct *AccountType) DeleteAccountType(key string, value interface{}) error {
	adapter := acct.dbPort.NewAccountTypeAdapter()
	// Check if exists
	_, err := adapter.Get(key, value)
	if err != nil {
		return &utils.RequestError{
			Code: http.StatusNotFound,
			Err: fmt.Errorf("account type does not exist"),
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