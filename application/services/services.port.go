package services

type AppServicesPort interface {
	// Add services here
	NewUserService() 		*User
	NewAdminService() 		*Admin
	NewAccountService() 	*Account
	NewUserAuthService()	*UserAuth
}