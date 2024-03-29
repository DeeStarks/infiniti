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

func (auth *UserAuth) AuthenticateUser(data map[string]interface{}) (UserFKResource, error) {
	requires := []string{"email", "password"}
	var missing string

	for _, field := range requires {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return UserFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required field(s): '%s'", missing[2:]),
		}
	}

	// Get the user
	email := data["email"].(string)
	password := data["password"].(string)

	userAdapter := auth.dbPort.NewUserAdapter()
	user, err := userAdapter.Get("email", email)
	if err != nil {
		return UserFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("User not found"),
		}
	}

	if err = auth.corePort.ComparePassword(user.Password, password); err != nil {
		return UserFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("invalid password"),
		}
	}

	// Get the user's account
	accountAdapter := auth.dbPort.NewAccountAdapter()
	account, err := accountAdapter.Get("user_id", user.Id)
	if err != nil {
		return UserFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Get user's group
	userGroup, err := auth.dbPort.NewUserGroupAdapter().Get("user_id", user.Id)
	if err != nil {
		return UserFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}
	group, err := auth.dbPort.NewGroupAdapter().Get("id", userGroup.GroupId)
	if err != nil {
		return UserFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Serialization and return
	var (
		userRes 	UserFKResource
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

func (auth *UserAuth) AuthenticateStaff(data map[string]interface{}) (StaffFKResource, error) {
	requires := []string{"email", "password"}
	var missing string

	for _, field := range requires {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required field(s): '%s'", missing[2:]),
		}
	}

	// Get the staff
	email := data["email"].(string)
	password := data["password"].(string)

	userAdapter := auth.dbPort.NewUserAdapter()
	user, err := userAdapter.Get("email", email)
	if err != nil {
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("Staff not found"),
		}
	}

	if err = auth.corePort.ComparePassword(user.Password, password); err != nil {
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("invalid password"),
		}
	}

	// Get group
	groupObj, err := auth.dbPort.NewUserGroupAdapter().Get("user_id", user.Id)
	if err != nil {
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}
	group, err := auth.dbPort.NewGroupAdapter().Get("id", groupObj.GroupId)
	if err != nil {
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Confirm it's a staff
	if group.Name != "staff" {
		return StaffFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("User is not a staff"),
		}
	}

	// Serialization and return
	var (
		staffRes 	StaffFKResource
		groupRes 	GroupResource
	)
	// User
	userJson, _ := json.Marshal(user)
	json.Unmarshal(userJson, &staffRes)
	// Group
	groupJson, _ := json.Marshal(group)
	json.Unmarshal(groupJson, &groupRes)

	// Combine and return
	staffRes.Group = groupRes
	return staffRes, nil
}


func (auth *UserAuth) AuthenticateAdmin(data map[string]interface{}) (AdminFKResource, error) {
	requires := []string{"email", "password"}
	var missing string

	for _, field := range requires {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required field(s): '%s'", missing[2:]),
		}
	}

	// Get the admin
	email := data["email"].(string)
	password := data["password"].(string)

	userAdapter := auth.dbPort.NewUserAdapter()
	user, err := userAdapter.Get("email", email)
	if err != nil {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("Admin not found"),
		}
	}

	if err = auth.corePort.ComparePassword(user.Password, password); err != nil {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("invalid password"),
		}
	}

	// Get group
	groupObj, err := auth.dbPort.NewUserGroupAdapter().Get("user_id", user.Id)
	if err != nil {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}
	group, err := auth.dbPort.NewGroupAdapter().Get("id", groupObj.GroupId)
	if err != nil {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Confirm it's an admin
	if group.Name != "admin" {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("User is not an admin"),
		}
	}

	// Serialization and return
	var (
		adminRes 	AdminFKResource
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

// id: user id;
// data: require "current_password" and "new_password";
func (auth *UserAuth) UpdatePassword(id int, data map[string]interface{}) error {
	// First check that all required fields are present
	var missing string
	for _, field := range []string{"new_password", "current_password"} {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required field(s): %s", missing[2:]),
		}
	}

	// Confirm user's existence
	user, err := auth.dbPort.NewUserAdapter().Get("id", id)
	if err != nil {
		return &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("user not found"),
		}
	}

	current := data["current_password"].(string)
	new := data["new_password"].(string)

	// Check password
	if err := auth.corePort.ComparePassword(user.Password, current); err != nil {
		return &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("invalid password"),
		}
	}
	hashedPassword, err := auth.corePort.HashPassword(new)
	if err != nil {
		return &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	err,
		}
	}

	userAdapter := auth.dbPort.NewUserAdapter()
	_, err = userAdapter.Update("id", id, map[string]interface{}{
		"password": hashedPassword,
	})
	if err != nil {
		return &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	err,
		}
	}
	return nil
}