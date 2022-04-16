package services

type AppServicesPort interface {
	// Add services here
	NewUserService() 		*User
	NewStaffService()		*Staff
	NewAdminService() 		*Admin
	NewAccountService() 	*Account
	NewUserAuthService()	*UserAuth
}