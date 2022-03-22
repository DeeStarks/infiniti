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

// Key -> Column to get data from
// Value -> Value to match
// Includes -> Account related columns to include in the result
func (account *Account) Get(key string, value string, includes []string) (constants.ServiceReturnType, error) {
	accountAdapter := account.dbPort.NewAccountAdapter()
	acct, err := accountAdapter.Get(key, value)
	if err != nil {
		return nil, err
	}

	serializedAcct := constants.ServiceReturnType(utils.StructToMap(acct))
	return serializedAcct, nil
}