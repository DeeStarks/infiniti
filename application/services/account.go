package services

import (
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/application/services/constants"
	"github.com/deestarks/infiniti/utils"
)

type Account struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

func (service *Service) NewAccountService() *Account {
	return &Account{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

// ACCOUNT

// Key -> Column to get data from
// Value -> Value to match
// Includes -> Account related columns to include in the result
func (account *Account) GetAccount(key string, value interface{}, includes []string) (constants.ServiceStructReturnType, utils.AppError) {
	accountAdapter := account.dbPort.NewAccountAdapter()
	acct, err := accountAdapter.Get(key, value)
	if err != nil {
		return nil, utils.NewAppError("Account not found", http.StatusNotFound)
	}

	serializedAcct := constants.ServiceStructReturnType(utils.StructToMap(acct))
	return serializedAcct, nil
}

func (account *Account) ListAccounts() (constants.ServiceStructSliceReturnType, utils.AppError) {
	accountAdapter := account.dbPort.NewAccountAdapter()
	accts, err := accountAdapter.List()
	if err != nil {
		return nil, utils.NewAppError(err.Error(), http.StatusInternalServerError)
	}

	serializedAccts := constants.ServiceStructSliceReturnType(utils.StructSliceToMap(accts))
	return serializedAccts, nil
}



// ACCOUNT TYPES
func (account *Account) GetAccountType(key string, value interface{}) (constants.ServiceStructReturnType, utils.AppError) {
	accountTypeAdapter := account.dbPort.NewAccountTypeAdapter()
	acctType, err := accountTypeAdapter.Get(key, value)
	if err != nil {
		return nil, utils.NewAppError("Account type not found", http.StatusNotFound)
	}

	serializedAcctType := constants.ServiceStructReturnType(utils.StructToMap(acctType))
	return serializedAcctType, nil
}