package db

import (
	"fmt"

	"github.com/deestarks/infiniti/utils"
)

type (
	UserPermissionsModel struct {
		Id 				int 	`json:"id"`
		UserId 			int 	`json:"user_id"`
		PermissionId 	int 	`json:"permission_id"`
	}
	
	UserPermissionsAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewUserPermissionsAdapter() *UserPermissionsAdapter {
	return &UserPermissionsAdapter{
		adapter: adpt,
		tableName: "user_permissions",
	}
}

func (mAdapt *UserPermissionsAdapter) Create(data map[string]interface{}) (*UserPermissionsModel, error) {
	var permission UserPermissionsModel

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
		RETURNING id, user_id, permission_id
	`, mAdapt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&permission.Id, &permission.UserId, &permission.PermissionId)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (mAdapt *UserPermissionsAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*UserPermissionsModel, error) {
	var (
		permission 	UserPermissionsModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, user_id, permission_id
	`, mAdapt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1)

	valArr = append(valArr, colValue)

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&permission.Id, &permission.UserId, &permission.PermissionId)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (mAdapt *UserPermissionsAdapter) Delete(colName string, value interface{}) (*UserPermissionsModel, error) {
	var (
		permission	UserPermissionsModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, user_id, permission_id
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&permission.Id, &permission.UserId, &permission.PermissionId)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// "Get" returns a single permission
func (mAdapt *UserPermissionsAdapter) Get(colName string, value interface{}) (*UserPermissionsModel, error) {
	var (
		permission	UserPermissionsModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT id, user_id, permission_id
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, colName)
	// `, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(
		&permission.Id, &permission.UserId, &permission.PermissionId,
	)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (mAdapt *UserPermissionsAdapter) Filter(colName string, value interface{}) (*[]UserPermissionsModel, error) {
	var (
		permissions	[]UserPermissionsModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, user_id, permission_id FROM %s WHERE %s = $1", mAdapt.tableName, colName)
	rows, err := mAdapt.adapter.db.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission UserPermissionsModel
		err := rows.Scan(
			&permission.Id, &permission.UserId, &permission.PermissionId,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return &permissions, nil
}

// "List" returns all permissions
func (mAdapt *UserPermissionsAdapter) List() (*[]UserPermissionsModel, error) {
	var (
		permissions	[]UserPermissionsModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, user_id, permission_id FROM %s", mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission UserPermissionsModel
		err := rows.Scan(
			&permission.Id, &permission.UserId, &permission.PermissionId,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return &permissions, nil
}