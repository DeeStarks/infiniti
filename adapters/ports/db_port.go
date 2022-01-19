package ports

import (
	"github.com/deestarks/infiniti/adapters/framework/db"
)

type DBPort interface {
	CloseDBConnection() error
	NewDBModel() *db.DBModel // For creation and deletion of tables
	
	// Add models here
	NewUserModel() *db.UserModel // CRUD operations on users
}