package ports

import (
	"github.com/deestarks/infiniti/adapters/framework/db"
)

type DBPort interface {
	CloseDBConnection() error
	
	// Add models here
	NewUserModel() *db.UserModel // CRUD operations on users
	NewPermissionsModel() *db.PermissionsModel // CRUD operations on permissions
}