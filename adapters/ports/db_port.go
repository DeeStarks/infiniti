package ports

type DBPort interface {
	CloseDBConnection()
	
	// Add models here
	NewUserModel()
}