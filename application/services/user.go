package services

import (
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/application/services/constants"
	"github.com/deestarks/infiniti/utils"
)

type User struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

func (service *Service) NewUserService() *User {
	return &User{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

func (user *User) CreateUser(data map[string]interface{}) (constants.ServiceStructReturnType, error) {
	userAdapter := user.dbPort.NewUserAdapter()
	// First check that all required fields are present
	required := []string{"email", "password", "first_name", "last_name", "account_type_id", "currency_id"}
	var missing string

	for _, field := range required {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return nil, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required fields: %s", missing),
		}
	}

	// Hash the password
	encyptPass, err := user.corePort.HashPassword(data["password"].(string))
	if err != nil {
		return nil, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Prepare user data
	userData := map[string]interface{}{
		"email": data["email"],
		"password": encyptPass,
		"first_name": data["first_name"],
		"last_name": data["last_name"],
	}
	
	// Create the user
	returnedUser, err := userAdapter.Create(userData)
	if err != nil {
		return nil, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	err,
		}
	}

	// Create the user's profile - account
	// First generate user's account number
	// This is done by replacing the last characters of a main number with the user's id
	mainNo := []byte("1220000000")
	userId := []byte(fmt.Sprintf("%d", returnedUser.Id))
	accountNo := fmt.Sprintf("%s%s", mainNo[:len(mainNo)-len(userId)], userId)

	// Prepare user's profile data
	accountData := map[string]interface{}{
		"user_id": returnedUser.Id,
		"account_type_id": data["account_type_id"],
		"currency_id": data["currency_id"],
		"account_number": accountNo,
	}

	// Creating the user's profile
	returnedAcct, err := user.dbPort.NewAccountAdapter().Create(accountData)
	if err != nil {
		return nil, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Add the user's profile to the returned user
	// Serialization and return
	serializedUser := constants.ServiceStructReturnType(utils.StructToMap(returnedUser))
	serializedAcct := constants.ServiceStructReturnType(utils.StructToMap(returnedAcct))

	// Add the account to the user
	serializedUser["account"] = serializedAcct
	
	// Return the user
	return serializedUser, nil
}