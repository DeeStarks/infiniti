package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)

type Admin struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

type AdminResource struct {
	Id 			int 	`json:"id"`
	FirstName 	string 	`json:"first_name"`
	LastName 	string 	`json:"last_name"`
	Email 		string 	`json:"email"`
	CreatedAt 	string 	`json:"created_at"`
}

type AdminFKResource struct { // Foreign key resource for admin
	Id 			int 			`json:"id"`
	FirstName 	string 			`json:"first_name"`
	LastName 	string 			`json:"last_name"`
	Email 		string 			`json:"email"`
	Group 		GroupResource 	`json:"group"`
	CreatedAt 	string 			`json:"created_at"`
}

// Initialize the admin service
func (service *Service) NewAdminService() *Admin {
	return &Admin{
		dbPort: 	service.dbPort,
		corePort: 	service.corePort,
	}
}

func (u *Admin) CreateAdmin(data map[string]interface{}) (AdminFKResource, error) {
	adapter := u.dbPort.NewUserAdapter()
	// First check that all required fields are present
	required := []string{"email", "password", "first_name", "last_name"}
	var missing string

	for _, field := range required {
		if _, ok := data[field]; !ok {
			missing = fmt.Sprintf("%s, %s", missing, field)
		}
	}
	if len(missing) > 0 {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("missing required field(s): %s", missing[2:]),
		}
	}

	// Hash the password
	encyptPass, err := u.corePort.HashPassword(data["password"].(string))
	if err != nil {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Prepare admin data
	adminData := map[string]interface{}{
		"email": data["email"],
		"password": encyptPass,
		"first_name": data["first_name"],
		"last_name": data["last_name"],
	}
	
	// Create the admin
	admin, err := adapter.Create(adminData)
	if err != nil {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	err,
		}
	}

	// Add the admin to the group
	group, err := u.dbPort.NewGroupAdapter().Get("name", "admin")
	if err != nil {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	_, err = u.dbPort.NewUserGroupAdapter().Create(map[string]interface{}{
		"user_id": admin.Id,
		"group_id": group.Id,
	})
	if err != nil {
		return AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	// Combine and return
	var (
		userRes 	AdminFKResource
		groupRes 	GroupResource
	)
	// User
	userJson, _ := json.Marshal(admin)
	json.Unmarshal(userJson, &userRes)
	// Group
	groupJson, _ := json.Marshal(group)
	json.Unmarshal(groupJson, &groupRes)

	// Combine and return
	userRes.Group = groupRes
	return userRes, nil
}

func (u *Admin) GetAdmin(key string, value interface{}) (AdminResource, error) {
	userAdapter := u.dbPort.NewUserAdapter()
	admin, err := userAdapter.Get(key, value)
	if err != nil {
		return AdminResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("admin not found"),
		}
	}
	var adminRes AdminResource
	adminJson, _ := json.Marshal(admin)
	json.Unmarshal(adminJson, &adminRes)
	return adminRes, nil
}

func (u *Admin) ListAdmins() ([]AdminFKResource, error) {
	userAdapter := u.dbPort.NewUserAdapter()
	conditions := map[string]interface{}{
		"groups.name": "admin",
	}
	selector := userAdapter.NewUserCustomSelector(conditions, "users.id", true).
		Join("user_groups", "user_id", "users", "id", []string{"user_id", "group_id"}).
		Join("groups", "id", "user_groups", "group_id", []string{"id", "name"})
	data, err := selector.Query()
	if err != nil {
		return []AdminFKResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err: 	err,
		}
	}

	var res []AdminFKResource
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

		// Prepare group data
		groupData := map[string]interface{}{
			"id": user["groups__id"],
			"name": user["groups__name"],
		}

		// Combine and return
		var adminRes AdminFKResource
		userJson, _ := json.Marshal(userData)
		json.Unmarshal(userJson, &adminRes)

		var groupRes GroupResource
		groupJson, _ := json.Marshal(groupData)
		json.Unmarshal(groupJson, &groupRes)

		adminRes.Group = groupRes
		res = append(res, adminRes)
	}
	return res, nil
}

func (u *Admin) UpdateAdmin(key string, value interface{}, data map[string]interface{}) (AdminResource, error) {
	for _, v := range []string{"id", "created_at"} { // These fields cannot be updated
		delete(data, v)
	}
	userAdapter := u.dbPort.NewUserAdapter()

	// Confirm admin exists
	_, err := userAdapter.Get(key, value)
	if err != nil {
		return AdminResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("admin not found"),
		}
	}

	// Update admin
	admin, err := userAdapter.Update(key, value, data)
	if err != nil {
		return AdminResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	err,
		}
	}
	var adminRes AdminResource
	adminJson, _ := json.Marshal(admin)
	json.Unmarshal(adminJson, &adminRes)
	return adminRes, nil
}

func (u *Admin) DeleteAdmin(id int) error {
	// First check if user exists
	userAdapter := u.dbPort.NewUserAdapter()
	_, err := userAdapter.Get("id", id)
	if err != nil {
		return &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("admin not found"),
		}
	}

	// Delete all references to the user
	// 1. "user_permissions" table
	userPermissionAdapter := u.dbPort.NewUserPermissionsAdapter()
	userPermissionAdapter.Delete("user_id", id)
	
	// 2. "user_groups" table
	userGroupAdapter := u.dbPort.NewUserGroupAdapter()
	userGroupAdapter.Delete("user_id", id)

	// 3. "user_accounts" table
	userAccountAdapter := u.dbPort.NewAccountAdapter()
	userAccountAdapter.Delete("user_id", id)

	// 4. "user_transactions" table
	userTransactionAdapter := u.dbPort.NewTransactionAdapter()
	userTransactionAdapter.Delete("user_id", id)

	// Delete the user
	userAdapter.Delete("id", id)
	return nil
}