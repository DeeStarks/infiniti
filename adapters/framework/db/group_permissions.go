package db

import (
	"fmt"

	"github.com/deestarks/infiniti/utils"
)

type (
	GroupPermissionsModel struct {
		Id 				int 	`json:"id"`
		GroupId 		int 	`json:"group_id"`
		PermissionId 	int 	`json:"permission_id"`
	}
	
	GroupPermissionsAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewGroupPermissionsAdapter() *GroupPermissionsAdapter {
	return &GroupPermissionsAdapter{
		adapter: adpt,
		tableName: "group_permissions",
	}
}

func (mAdapt *GroupPermissionsAdapter) Create(data map[string]interface{}) (*GroupPermissionsModel, error) {
	var permission GroupPermissionsModel

	mToS := utils.MapToStructSlice(data)
	var (
		colStr	string
		valArr	[]interface{}
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
		RETURNING id, group_id, permission_id
	`, mAdapt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&permission.Id, &permission.GroupId, &permission.PermissionId)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (mAdapt *GroupPermissionsAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*GroupPermissionsModel, error) {
	var (
		permission 	GroupPermissionsModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, group_id, permission_id
	`, mAdapt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1)

	valArr = append(valArr, colValue)

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&permission.Id, &permission.GroupId, &permission.PermissionId)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (mAdapt *GroupPermissionsAdapter) Delete(colName string, value interface{}) (*GroupPermissionsModel, error) {
	var (
		permission	GroupPermissionsModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, group_id, permission_id
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&permission.Id, &permission.GroupId, &permission.PermissionId)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (mAdapt *GroupPermissionsAdapter) Get(colName string, value interface{}) (*GroupPermissionsModel, error) {
	var (
		permission	GroupPermissionsModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT id, group_id, permission_id
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&permission.Id, &permission.GroupId, &permission.PermissionId)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (mAdapt *GroupPermissionsAdapter) Filter(colName string, value interface{}) (*[]GroupPermissionsModel, error) {
	var (
		permissions	[]GroupPermissionsModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, group_id, permission_id FROM %s WHERE %s = $1", mAdapt.tableName, colName)
	rows, err := mAdapt.adapter.db.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission GroupPermissionsModel
		err := rows.Scan(&permission.Id, &permission.GroupId, &permission.PermissionId)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return &permissions, nil
}

func (mAdapt *GroupPermissionsAdapter) List() (*[]GroupPermissionsModel, error) {
	var (
		permissions	[]GroupPermissionsModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, group_id, permission_id FROM %s", mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission GroupPermissionsModel
		err := rows.Scan(&permission.Id, &permission.GroupId, &permission.PermissionId)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return &permissions, nil
}