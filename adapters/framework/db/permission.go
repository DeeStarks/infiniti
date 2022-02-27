package db

import (
	"fmt"
	"log"
	"github.com/deestarks/infiniti/utils"
)

type (
	PermissionsModel struct {
		Id					int 					`json:"id"`
		TableId				int 					`json:"table_id"`
		Method				string 					`json:"method"`
		// UserPermissions   	UserPermissionsModel	`json:"user_permissions"` // One to many relationship
		// GroupPermissions  	GroupPermissionsModel	`json:"group_permissions"` // One to many relationship
	}

	PermissionsAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}

	_permissionsEnum struct {
		Get, Create, Update, Delete 	string
	}
)

func PermissionsEnum() _permissionsEnum {
	return _permissionsEnum{
		Get: "get",
		Create: "post",
		Update: "put",
		Delete: "delete",
	}
}

func (adpt *DBAdapter) NewPermissionsAdapter() *PermissionsAdapter {
	return &PermissionsAdapter{
		adapter: adpt,
		tableName: "permissions",
	}
}

func (pAdpt *PermissionsAdapter) Create(data map[string]interface{}) (*PermissionsModel, error) {
	var permission PermissionsModel

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
		RETURNING id, table_id, method
	`, pAdpt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := pAdpt.adapter.db.QueryRow(query, valArr...).Scan(&permission.Id, &permission.TableId, &permission.Method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &permission, nil
}

func (pAdpt *PermissionsAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*PermissionsModel, error) {
	var (
		permission 	PermissionsModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, table_id, method
	`, pAdpt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1) // CreateSetConditions will create
	// conditions and create placeholders for the values (e.g. id = $1, name = $2...$n)
	// The last "len(data)+1" is for the value of the column to be updated.

	// Add the value of the column to be updated to the end of the array of values.
	valArr = append(valArr, colValue)

	err := pAdpt.adapter.db.QueryRow(query, valArr...).Scan(&permission.Id, &permission.TableId, &permission.Method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &permission, nil
}

func (pAdpt *PermissionsAdapter) Delete(colName string, value interface{}) (*PermissionsModel, error) {
	var (
		permission	PermissionsModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, table_id, method
	`, pAdpt.tableName, colName)
	err = pAdpt.adapter.db.QueryRow(query, value).Scan(&permission.Id, &permission.TableId, &permission.Method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &permission, nil
}

// "Get" returns a single permission
func (pAdpt *PermissionsAdapter) Get(colName string, value interface{}) (*PermissionsModel, error) {
	var (
		permission	PermissionsModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT id, table_id, method
		FROM %s
		WHERE %s = $1
	`, pAdpt.tableName, colName)

	// This will replace the above query when the UserPermissionsModel and GroupPermissionsModel are implemented
	// query := fmt.Sprintf(`
	// 	SELECT %[1]s.id, %[1]s.table_id, %[1]s.method, user_permissions.id, user_permissions.user_id, user_permissions.permission_id,
	// 		group_permissions.id, group_permissions.group_id, group_permissions.permission_id
	// 	FROM %[1]s	
	// 	WHERE %[2]s = $1
	// 	OUTER JOIN user_permissions ON user_permissions.permission_id = permissions.id
	// 	OUTER JOIN group_permissions ON group_permissions.permission_id = permissions.id
	// `, pAdpt.tableName, colName)
	err = pAdpt.adapter.db.QueryRow(query, value).Scan(
		&permission.Id, &permission.TableId, &permission.Method,
		// &permission.UserPermissions.Id, &permission.UserPermissions.UserId, &permission.UserPermissions.PermissionId,
		// &permission.GroupPermissions.Id, &permission.GroupPermissions.GroupId, &permission.GroupPermissions.PermissionId,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &permission, nil
}

func (pAdpt *PermissionsAdapter) Filter(colName string, value interface{}) (*[]PermissionsModel, error) {
	var (
		permissions	[]PermissionsModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, table_id, method FROM %s WHERE %s = $1", pAdpt.tableName, colName)

	// This will replace the above query when the UserPermissionsModel and GroupPermissionsModel are implemented
	// query := fmt.Sprintf(`
	// 	SELECT %[1]s.id, %[1]s.table_id, %[1]s.method, user_permissions.id, user_permissions.user_id, user_permissions.permission_id,
	// 		group_permissions.id, group_permissions.group_id, group_permissions.permission_id
	// 	FROM %[1]s 
	// 	WHERE %s = $1
	// 	OUTER JOIN user_permissions ON user_permissions.permission_id = permissions.id
	// 	OUTER JOIN group_permissions ON group_permissions.permission_id = permissions.id
	// `, pAdpt.tableName, colName)
	rows, err := pAdpt.adapter.db.Query(query, value)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission PermissionsModel
		err := rows.Scan(
			&permission.Id, &permission.TableId, &permission.Method,
			// &permission.UserPermissions.Id, &permission.UserPermissions.UserId, &permission.UserPermissions.PermissionId,
			// &permission.GroupPermissions.Id, &permission.GroupPermissions.GroupId, &permission.GroupPermissions.PermissionId,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return &permissions, nil
}

// "List" returns all permissions
func (pAdpt *PermissionsAdapter) List() (*[]PermissionsModel, error) {
	var (
		permissions	[]PermissionsModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, table_id, method FROM %s", pAdpt.tableName)
	// This will replace the above query when the UserPermissionsModel and GroupPermissionsModel are implemented
	// query := fmt.Sprintf(`
	// 	SELECT %[1]s.id, %[1]s.table_id, %[1]s.method, user_permissions.id, user_permissions.user_id, user_permissions.permission_id,
	// 		group_permissions.id, group_permissions.group_id, group_permissions.permission_id
	//	FROM %s
	// 	OUTER JOIN user_permissions ON user_permissions.permission_id = permissions.id
	// 	OUTER JOIN group_permissions ON group_permissions.permission_id = permissions.id
	// `, pAdpt.tableName)
	rows, err := pAdpt.adapter.db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission PermissionsModel
		err := rows.Scan(
			&permission.Id, &permission.TableId, &permission.Method,
			// &permission.UserPermissions.Id, &permission.UserPermissions.UserId, &permission.UserPermissions.PermissionId,
			// &permission.GroupPermissions.Id, &permission.GroupPermissions.GroupId, &permission.GroupPermissions.PermissionId,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return &permissions, nil
}