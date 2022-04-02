package services

type AppServicesPort interface {
	// Add services here
	NewUserService() 		*User
	NewAccountService() 	*Account
	NewUserAuthService()	*UserAuth
}