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
	CreatedAt 	string 				`json:"created_at"`
}

type UserFKResource struct { // User with foreign keys
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

func (u *User) CreateUser(data map[string]interface{}) (UserFKResource, error) {
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
		return UserFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required field(s): %s", missing[2:]),
		}
	}

	// Hash the password
	encyptPass, err := u.corePort.HashPassword(data["password"].(string))
	if err != nil {
		return UserFKResource{}, &utils.RequestError{
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
		return UserFKResource{}, &utils.RequestError{
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
		return UserFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Create user's group
	// Get the group id
	group, err := u.dbPort.NewGroupAdapter().Get("name", "user")
	if err != nil {
		return UserFKResource{}, &utils.RequestError{
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
		return UserFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Combine and return
	var (
		userRes 	UserFKResource
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

func (u *User) GetUser(key string, value interface{}) (UserResource, error) {
	userAdapter := u.dbPort.NewUserAdapter()
	user, err := userAdapter.Get(key, value)
	if err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("user not found"),
		}
	}
	var userRes UserResource
	userJson, _ := json.Marshal(user)
	json.Unmarshal(userJson, &userRes)
	return userRes, nil
}

func (u *User) ListUsers() ([]UserFKResource, error) {
	userAdapter := u.dbPort.NewUserAdapter()
	selector := userAdapter.NewUserCustomSelector("groups.name", "user", "users.id", true).
		Join("user_groups", "user_id", "users", "id", []string{"user_id", "group_id"}).
		Join("groups", "id", "user_groups", "group_id", []string{"id", "name"}).
		Join("user_accounts", "user_id", "users", "id", []string{"id", "user_id", "account_type_id", "account_number", "balance", "currency_id"})
	data := selector.Query()

	var res []UserFKResource
	for _, user := range data {
		// Prepare user data
		userData := map[string]interface{}{
			"id": user["users__id"],
			"email": user["users__email"],
			"password": user["users__password"],
			"first_name": user["users__first_name"],
			"last_name": user["users__last_name"],
			"created_at": user["users__created_at"],
		}
		
		// Prepare account data
		accountData := map[string]interface{}{
			"id": user["user_accounts__id"],
			"user_id": user["user_accounts__user_id"],
			"account_type_id": user["user_accounts__account_type_id"],
			"account_number": user["user_accounts__account_number"],
			"balance": user["user_accounts__balance"],
			"currency_id": user["user_accounts__currency_id"],
		}

		// Prepare group data
		groupData := map[string]interface{}{
			"id": user["groups__id"],
			"name": user["groups__name"],
		}

		// Combine and return
		var userRes UserFKResource
		userJson, _ := json.Marshal(userData)
		json.Unmarshal(userJson, &userRes)

		var accountRes AccountResource
		accountJson, _ := json.Marshal(accountData)
		json.Unmarshal(accountJson, &accountRes)

		var groupRes GroupResource
		groupJson, _ := json.Marshal(groupData)
		json.Unmarshal(groupJson, &groupRes)

		userRes.Account = accountRes
		userRes.Group = groupRes
		res = append(res, userRes)
	}
	return res, nil
}

func (u *User) UpdateUser(key string, value interface{}, data map[string]interface{}) (UserResource, error) {
	for _, v := range []string{"id", "created_at"} { // These fields cannot be updated
		delete(data, v)
	}
	userAdapter := u.dbPort.NewUserAdapter()
	user, err := userAdapter.Update(key, value, data)
	if err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	err,
		}
	}
	var userRes UserResource
	userJson, _ := json.Marshal(user)
	json.Unmarshal(userJson, &userRes)
	return userRes, nil
}