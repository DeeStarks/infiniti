package services

type AppServicesPort interface {
	NewAccountService() 	*Account
}