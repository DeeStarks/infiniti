package services

type AppServicesPort interface {
	// Add services here
	NewUserService() 		*User
	NewStaffService()		*Staff
	NewAdminService() 		*Admin
	NewAccountService() 	*Account
	NewAccountTypeService() *AccountType
	NewUserAuthService()	*UserAuth
	NewCurrencyService()	*Currency
}