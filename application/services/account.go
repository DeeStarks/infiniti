package services

import (
	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
	"github.com/deestarks/infiniti/application/services/constants"
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
func (account *Account) GetAccount(key string, value interface{}, includes []string) (constants.ServiceStructReturnType, error) {
	accountAdapter := account.dbPort.NewAccountAdapter()
	acct, err := accountAdapter.Get(key, value)
	if err != nil {
		return nil, err
	}

	serializedAcct := constants.ServiceStructReturnType(utils.StructToMap(acct))
	return serializedAcct, nil
}

func (account *Account) ListAccounts() (constants.ServiceStructSliceReturnType, error) {
	accountAdapter := account.dbPort.NewAccountAdapter()
	accts, err := accountAdapter.List()
	if err != nil {
		return nil, err
	}

	serializedAccts := constants.ServiceStructSliceReturnType(utils.StructSliceToMap(accts))
	return serializedAccts, nil
}



// ACCOUNT TYPES
func (account *Account) GetAccountType(key string, value interface{}) (constants.ServiceStructReturnType, error) {
	accountTypeAdapter := account.dbPort.NewAccountTypeAdapter()
	acctType, err := accountTypeAdapter.Get(key, value)
	if err != nil {
		return nil, err
	}

	serializedAcctType := constants.ServiceStructReturnType(utils.StructToMap(acctType))
	return serializedAcctType, nil
}