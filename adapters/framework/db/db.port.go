package db

type DBPort interface {
	CloseDBConnection() 			error
	
	// DB Adapters
	NewUserAdapter() 				*UserAdapter
	NewPermissionsAdapter() 		*PermissionsAdapter // permissions
	NewAccountTypeAdapter() 		*AccountTypeAdapter // account types
	NewAccountAdapter() 			*AccountAdapter // account
	NewCurrencyAdapter() 			*CurrencyAdapter // currencies
	NewGroupPermissionsAdapter() 	*GroupPermissionsAdapter // group permissions
	NewGroupAdapter() 				*GroupAdapter // groups
	NewTransactionTypeAdapter() 	*TransactionTypeAdapter // transaction types
	NewTransactionAdapter() 		*TransactionAdapter // transactions
	NewUserGroupAdapter() 			*UserGroupAdapter // user groups
	NewUserPermissionsAdapter() 	*UserPermissionsAdapter // user permissions
}