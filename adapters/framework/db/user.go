package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/deestarks/infiniti/utils"
	global "github.com/deestarks/infiniti/adapters/framework/db/global"
)

type (
	UserModel struct {
		Id 					int 					`json:"id"`
		FirstName 			string 					`json:"first_name"`
		LastName 			string 					`json:"last_name"`
		Email 				string 					`json:"email"`
		Password 			string 					`json:"password"`
		CreatedAt 			time.Time 				`json:"created_at"`
		Account 			AccountModel			`json:"account"` // One to One relationship
		AccountType 		AccountTypeModel		`json:"account_type"` // One to One relationship
		Currency			CurrencyModel			`json:"currency"` // One to One relationship
		UserPermissions 	[]UserPermissionsModel	`json:"user_permissions"` // One to Many relationship
		UserGroups 			[]UserGroupModel		`json:"user_groups"` // One to Many relationship - Many to Many relationship
	}
	
	UserAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

var (
	selection = `
		users.id, users.first_name, users.last_name, users.email, users.password, users.created_at,
		COALESCE(account.id, 0), COALESCE(account.user_id, 0), COALESCE(account.account_type_id, 0), 
		COALESCE(account.account_number, 0), COALESCE(account.balance, 0.0), COALESCE(account.currency_id, 0),
		COALESCE(account_type.id, 0), COALESCE(account_type.name, ''),
		COALESCE(currency.id, 0), COALESCE(currency.name, ''), COALESCE(currency.symbol, ''),
		COALESCE(permission.id, 0), COALESCE(permission.user_id, 0), COALESCE(permission.permission_id, 0),
		COALESCE(user_group.user_id, 0), COALESCE(user_group.group_id, 0)
	`

	selectionJoin = `
		LEFT OUTER JOIN user_accounts AS account ON account.user_id = users.id
		LEFT OUTER JOIN account_types AS account_type ON account_type.id = account.account_type_id
		LEFT OUTER JOIN currencies AS currency ON currency.id = account.currency_id
		LEFT OUTER JOIN user_permissions AS permission ON permission.user_id = users.id
		LEFT OUTER JOIN user_groups AS user_group ON user_group.user_id = users.id
	`
)

func populateUserModel(row *sql.Rows, many bool) (interface{}, error) {
	var (
		// notInitial 		bool // Used to determine if the first row is being populated
		permission 		UserPermissionsModel
		userGroup 		UserGroupModel
		SingleModel		UserModel
		ManyModels		[]UserModel
	)

	// The following is to avoid duplicates
	userHits 		:= make(global.PreviousIdHits)
	permissionHits 	:= make(global.PreviousIdHits)
	userGroupHits 	:= make(global.PreviousIdHits)

	for row.Next() {
		if err := row.Scan(
			&SingleModel.Id, &SingleModel.FirstName, &SingleModel.LastName,
			&SingleModel.Email, &SingleModel.Password, &SingleModel.CreatedAt,
			&SingleModel.Account.Id, &SingleModel.Account.UserId, &SingleModel.Account.AccountTypeId, 
			&SingleModel.Account.AccountNumber, &SingleModel.Account.Balance, &SingleModel.Account.CurrencyId,
			&SingleModel.AccountType.Id, &SingleModel.AccountType.Name,
			&SingleModel.Currency.Id, &SingleModel.Currency.Name, &SingleModel.Currency.Symbol,
			&permission.Id, &permission.UserId, &permission.PermissionId,
			&userGroup.UserId, &userGroup.GroupId,
		); err != nil {
			log.Fatal(err)
			return nil, err
		}

		if !many {
			if !permissionHits[permission.Id] { // If the permission has not been hit, add it to the user's permissions
				SingleModel.UserPermissions = append(SingleModel.UserPermissions, permission)
				permissionHits[permission.Id] = true
			}
			if !userGroupHits[userGroup.GroupId] { // If the user group has not been hit, add it to the user's user groups
				SingleModel.UserGroups = append(SingleModel.UserGroups, userGroup)
				userGroupHits[userGroup.GroupId] = true
			}
		} else {
			// <=========== TODO: Complete this section ===========>

			if !userHits[SingleModel.Id] { // If the user has not been hit
				if SingleModel.Id != 0 { // If the user model has been populated
					fmt.Println(SingleModel.Id)
					ManyModels 		= append(ManyModels, SingleModel) // Add the user model to the list of user models
					SingleModel 	= UserModel{} // Reset the user model
					permissionHits 	= make(global.PreviousIdHits) // Reset the permission hits
					userGroupHits 	= make(global.PreviousIdHits) // Reset the user group hits
				}

				SingleModel.UserPermissions = append(SingleModel.UserPermissions, permission)
				SingleModel.UserGroups = append(SingleModel.UserGroups, userGroup)
				permissionHits[permission.Id] = true // Add the permission to the permissionHits map to prevent duplicates
				userGroupHits[userGroup.GroupId] = true // Add the user group to the userGroupHits map to prevent duplicates
				userHits[SingleModel.Id] = true // Add the user to the userHits map to prevent duplicates
				fmt.Println(userHits[SingleModel.Id])
			} else { // If this is not the first row, populate with only the foreign keys
				if !permissionHits[permission.Id] { // If the permission has not been hit, add it to the user's permissions
					SingleModel.UserPermissions = append(SingleModel.UserPermissions, permission)
					permissionHits[permission.Id] = true
				}
				if !userGroupHits[userGroup.GroupId] { // If the user group has not been hit, add it to the user's user groups
					SingleModel.UserGroups = append(SingleModel.UserGroups, userGroup)
					userGroupHits[userGroup.GroupId] = true
				}
			}

			// <========================================>
		}
	}
	if many {
		if SingleModel.Id != 0 { // If the user model has been populated
			ManyModels = append(ManyModels, SingleModel) // Add the last user model to the list of user models
		}
		return &ManyModels, nil
	}
	return &SingleModel, nil
}

func (adpt *DBAdapter) NewUserAdapter() *UserAdapter {
	return &UserAdapter{
		adapter: adpt,
		tableName: "users",
	}
}

// Define the methods of the UserAdapter
func (mAdapt *UserAdapter) Get(col string, value interface{}) (*UserModel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM %s
		%s
		WHERE users.%s = $1
	`, selection, mAdapt.tableName, selectionJoin, col)

	rows, err := mAdapt.adapter.db.Query(query, value)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	arrangedUser, err := populateUserModel(rows, false)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	user := arrangedUser.(*UserModel)
	return user, nil
}

// TODO: Populate the user model with the user's permissions and user groups
func (mAdapt *UserAdapter) Filter(col string, value interface{}) (*[]UserModel, error) {
	query := fmt.Sprintf(`
		SELECT id, first_name, last_name, email, password, created_at FROM %s
		WHERE %s = $1 ORDER BY id ASC
	`, mAdapt.tableName, col)
	rows, err := mAdapt.adapter.db.Query(query, value)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	arrangedUsers, err := populateUserModel(rows, true)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	users := arrangedUsers.(*[]UserModel)
	return users, nil
}

// TODO: Populate the user model with the user's permissions and user groups
func (mAdapt *UserAdapter) List() (*[]UserModel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM %s
		%s
		ORDER BY id ASC
	`, selection, mAdapt.tableName, selectionJoin)
	rows, err := mAdapt.adapter.db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	
	arrangedUsers, err := populateUserModel(rows, true)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	users := arrangedUsers.(*[]UserModel)
	return users, nil
}

func (mAdapt *UserAdapter) Create(data map[string]interface{}) (*UserModel, error) {
	var user UserModel

	mToS := utils.MapToStructSlice(data)
	var (
		colStr		string
		valArr		[]interface{}
	)
	for i, s := range mToS {
		colStr += s.Key + ", "
		valArr = append(valArr, s.Value)
		if i == len(mToS)-1 {
			colStr = colStr[:len(colStr)-2] // remove the last ", "
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %s ( %s ) VALUES ( %s )
		RETURNING id, first_name, last_name, email, password, created_at
	`, mAdapt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(
		&user.Id, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.CreatedAt,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &user, nil
}

func (mAdapt *UserAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*UserModel, error) {
	var (
		user 	UserModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, first_name, last_name, email, password, created_at
	`, mAdapt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1)

	valArr = append(valArr, colValue)

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(
		&user.Id, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.CreatedAt,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &user, nil
}

func (mAdapt *UserAdapter) Delete(colName string, value interface{}) (*UserModel, error) {
	var (
		user	UserModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, first_name, last_name, email, password, created_at
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(
		&user.Id, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.CreatedAt,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &user, nil
}