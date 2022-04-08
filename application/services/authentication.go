package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)

type UserAuth struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

func (service *Service) NewUserAuthService() *UserAuth {
	return &UserAuth{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

func (auth *UserAuth) AuthenticateUser(data map[string]interface{}) (UserResource, error) {
	requires := []string{"email", "password"}
	var missing string

	for _, field := range requires {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required fields: '%s'", missing[2:]),
		}
	}

	// Get the user
	email := data["email"].(string)
	password := data["password"].(string)

	userAdapter := auth.dbPort.NewUserAdapter()
	user, err := userAdapter.Get("email", email)
	if err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("User not found"),
		}
	}

	if err = auth.corePort.ComparePassword(user.Password, password); err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("invalid password"),
		}
	}

	// Get the user's account
	accountAdapter := auth.dbPort.NewAccountAdapter()
	account, err := accountAdapter.Get("user_id", user.Id)
	if err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Get user's group
	userGroup, err := auth.dbPort.NewUserGroupAdapter().Get("user_id", user.Id)
	if err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}
	group, err := auth.dbPort.NewGroupAdapter().Get("id", userGroup.GroupId)
	if err != nil {
		return UserResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Serialization and return
	var (
		userRes 	UserResource
		acctRes 	AccountResource
		groupRes 	GroupResource
	)
	// User
	userJson, _ := json.Marshal(user)
	json.Unmarshal(userJson, &userRes)
	// Account
	acctJson, _ := json.Marshal(account)
	json.Unmarshal(acctJson, &acctRes)
	// Group
	groupJson, _ := json.Marshal(group)
	json.Unmarshal(groupJson, &groupRes)

	// Combine and return
	userRes.Account = acctRes
	userRes.Group = groupRes
	return userRes, nil
}

func (auth *UserAuth) AuthenticateAdmin(data map[string]interface{}) (AdminResource, error) {
	requires := []string{"email", "password"}
	var missing string

	for _, field := range requires {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return AdminResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required fields: '%s'", missing[2:]),
		}
	}

	// Get the admin
	email := data["email"].(string)
	password := data["password"].(string)

	userAdapter := auth.dbPort.NewUserAdapter()
	user, err := userAdapter.Get("email", email)
	if err != nil {
		return AdminResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("Admin not found"),
		}
	}

	if err = auth.corePort.ComparePassword(user.Password, password); err != nil {
		return AdminResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("invalid password"),
		}
	}

	// Get group
	groupObj, err := auth.dbPort.NewUserGroupAdapter().Get("user_id", user.Id)
	if err != nil {
		return AdminResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}
	group, err := auth.dbPort.NewGroupAdapter().Get("id", groupObj.GroupId)
	if err != nil {
		return AdminResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Confirm it's an admin
	if group.Name != "admin" {
		return AdminResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("User is not an admin"),
		}
	}

	// Serialization and return
	var (
		adminRes 	AdminResource
		groupRes 	GroupResource
	)
	// User
	userJson, _ := json.Marshal(user)
	json.Unmarshal(userJson, &adminRes)
	// Group
	groupJson, _ := json.Marshal(group)
	json.Unmarshal(groupJson, &groupRes)

	// Combine and return
	adminRes.Group = groupRes
	return adminRes, nil
}
