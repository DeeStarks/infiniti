package ports

import (
	"github.com/deestarks/infiniti/adapters/framework/db"
)

type DBPort interface {
	CloseDBConnection() error
	
	// Add models here
	NewUserAdapter() *db.UserAdapter // CRUD operations on users
	NewPermissionsAdapter() *db.PermissionsAdapter // CRUD operations on permissions
}