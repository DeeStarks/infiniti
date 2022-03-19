package ports

import (
	"github.com/deestarks/infiniti/adapters/framework/db"
)

type DBPort interface {
	CloseDBConnection() 			error
	
	// Add models here
	NewUserAdapter() 				*db.UserAdapter
	NewPermissionsAdapter() 		*db.PermissionsAdapter // permissions
	NewAccountTypeAdapter() 		*db.AccountTypeAdapter // account types
	NewAccountAdapter() 			*db.AccountAdapter // account
	NewCurrencyAdapter() 			*db.CurrencyAdapter // currencies
	NewGroupPermissionsAdapter() 	*db.GroupPermissionsAdapter // group permissions
	NewGroupAdapter() 				*db.GroupAdapter // groups
	NewTransactionTypeAdapter() 	*db.TransactionTypeAdapter // transaction types
	NewTransactionAdapter() 		*db.TransactionAdapter // transactions
	NewUserGroupAdapter() 			*db.UserGroupAdapter // user groups
	NewUserPermissionsAdapter() 	*db.UserPermissionsAdapter // user permissions
}