package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)

type Account struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

type AccountResource struct {
	Id 				int 		`json:"id"`
	UserId 			int 		`json:"user_id"`
	AccountTypeId 	int 		`json:"account_type_id"`
	AccountNumber 	string 		`json:"account_number"`
	Balance 		float64 	`json:"balance"`
	CurrencyId 		int 		`json:"currency_id"`
}

func (service *Service) NewAccountService() *Account {
	return &Account{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

// key -> Column to get data from
// value -> Value to match
// includes -> Account related columns to include in the result
func (account *Account) GetAccount(key string, value interface{}, includes []string) (AccountResource, error) {
	accountAdapter := account.dbPort.NewAccountAdapter()
	acct, err := accountAdapter.Get(key, value)
	if err != nil {
		return AccountResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf("Account not found"),
		}
	}

	// Serialization and return
	var res AccountResource
	jsonAcct, _ := json.Marshal(acct)
	json.Unmarshal(jsonAcct, &res)
	return res, nil
}

// Update an account
// value -> value to match
// key -> column to match
// data -> update data
func (account *Account) UpdateAccount(key string, value interface{}, data map[string]interface{}) (AccountResource, error) {
	rsrts := []string{"id", "user_id", "account_number", "balance"} // Restricted columns to update
	var foundRsrts string // Restrictions found
	// Check if the data contains restricted columns
	for _, field := range rsrts {
		if _, ok := data[field]; ok {
			foundRsrts = fmt.Sprintf("%s, %s", foundRsrts, field)
		}
	}
	if len(foundRsrts) > 0 {
		return AccountResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("cannot update restricted fields: '%s'", foundRsrts),
		}
	}

	accountAdapter := account.dbPort.NewAccountAdapter()
	acct, err := accountAdapter.Update(key, value, data)
	if err != nil {
		return AccountResource{}, &utils.RequestError{
			Code:	http.StatusNotFound,
			Err: 	fmt.Errorf(err.Error()),
		}
	}

	// Serialization and return
	var res AccountResource
	jsonAcct, _ := json.Marshal(acct)
	json.Unmarshal(jsonAcct, &res)
	return res, nil
}

func (account *Account) ListAccounts() ([]AccountResource, error) {
	accountAdapter := account.dbPort.NewAccountAdapter()
	accts, err := accountAdapter.List()
	if err != nil {
		return []AccountResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	fmt.Errorf("error listing accounts"),
		}
	}

	// Serialization and return
	var res []AccountResource
	jsonAccts, _ := json.Marshal(accts)
	json.Unmarshal(jsonAccts, &res)
	return res, nil
}