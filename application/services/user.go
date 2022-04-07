package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)

type User struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

type UserResource struct {
	Id 			int 				`json:"id"`
	FirstName 	string 				`json:"first_name"`
	LastName 	string 				`json:"last_name"`
	Email 		string 				`json:"email"`
	Group 		GroupResource 		`json:"group"`
	Account 	AccountResource		`json:"account"`
	CreatedAt 	string 				`json:"created_at"`
}

func (service *Service) NewUserService() *User {
	return &User{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

func (u *User) CreateUser(data map[string]interface{}) (UserResource, error) {
	userAdapter := u.dbPort.NewUserAdapter()
	// First check that all required fields are present
	required := []string{"email", "password", "first_name", "last_name", "account_type_id", "currency_id"}
	var missing string

	for _, field := range required {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required fields: %s", missing),
		}
	}

	// Hash the password
	encyptPass, err := u.corePort.HashPassword(data["password"].(string))
	if err != nil {
		return UserResource{}, &utils.RequestError{
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
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	err,
		}
	}

	// Create the user's profile - account
	// First generate user's account number
	// This is done by replacing the last characters of a main number with the user's id
	mainNo := 1220000000
	accountNo := strconv.Itoa(mainNo+returnedUser.Id) // Type of "account_number" in the database is  a varchar

	// Creating the user's profile
	returnedAcct, err := u.dbPort.NewAccountAdapter().Create(map[string]interface{}{
		"user_id": returnedUser.Id,
		"account_type_id": data["account_type_id"],
		"currency_id": data["currency_id"],
		"account_number": accountNo,
	})
	if err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Create user's group
	// Get the group id
	group, err := u.dbPort.NewGroupAdapter().Get("name", "user")
	if err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Create the user's group
	_, err = u.dbPort.NewUserGroupAdapter().Create(map[string]interface{}{
		"user_id": returnedUser.Id,
		"group_id": group.Id,
	})
	if err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Combine and return
	var (
		userRes 	UserResource
		groupRes 	GroupResource
		acctRes 	AccountResource
	)
	// User
	userJson, _ := json.Marshal(returnedUser)
	json.Unmarshal(userJson, &userRes)
	// Group
	groupJson, _ := json.Marshal(group)
	json.Unmarshal(groupJson, &groupRes)
	// Account
	acctJson, _ := json.Marshal(returnedAcct)
	json.Unmarshal(acctJson, &acctRes)

	// Combine and return
	userRes.Group = groupRes
	userRes.Account = acctRes
	return userRes, nil
}