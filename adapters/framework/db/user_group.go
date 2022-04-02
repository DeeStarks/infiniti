package db

import (
	"fmt"

	"github.com/deestarks/infiniti/utils"
	"github.com/lib/pq"
)

type (
	UserGroupModel struct {
		GroupId 	int 	`json:"group_id"`
		UserId 		int 	`json:"user_id"`
	}
	
	UserGroupAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewUserGroupAdapter() *UserGroupAdapter {
	return &UserGroupAdapter{
		adapter: adpt,
		tableName: "user_groups",
	}
}

func (mAdapt *UserGroupAdapter) Create(data map[string]interface{}) (*UserGroupModel, error) {
	var userGroup UserGroupModel

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
		RETURNING group_id, user_id
	`, mAdapt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&userGroup.GroupId, &userGroup.UserId)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &userGroup, nil
}

func (mAdapt *UserGroupAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*UserGroupModel, error) {
	var (
		userGroup 	UserGroupModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING group_id, user_id
	`, mAdapt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1)

	valArr = append(valArr, colValue)

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&userGroup.GroupId, &userGroup.UserId)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &userGroup, nil
}

func (mAdapt *UserGroupAdapter) Delete(colName string, value interface{}) (*UserGroupModel, error) {
	var (
		userGroup	UserGroupModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING group_id, user_id
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&userGroup.GroupId, &userGroup.UserId)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &userGroup, nil
}

func (mAdapt *UserGroupAdapter) Get(colName string, value interface{}) (*UserGroupModel, error) {
	var (
		userGroup	UserGroupModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT group_id, user_id
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&userGroup.GroupId, &userGroup.UserId)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &userGroup, nil
}

func (mAdapt *UserGroupAdapter) Filter(colName string, value interface{}) (*[]UserGroupModel, error) {
	var (
		userGroups	[]UserGroupModel
		err			error
	)
	query := fmt.Sprintf("SELECT group_id, user_id FROM %s WHERE %s = $1", mAdapt.tableName, colName)
	rows, err := mAdapt.adapter.db.Query(query, value)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var userGroup UserGroupModel
		err := rows.Scan(&userGroup.GroupId, &userGroup.UserId)
		if err, ok := err.(*pq.Error); ok {
			return nil, fmt.Errorf("%s", err.Detail)
		}
		userGroups = append(userGroups, userGroup)
	}
	return &userGroups, nil
}

func (mAdapt *UserGroupAdapter) List() (*[]UserGroupModel, error) {
	var (
		userGroups	[]UserGroupModel
		err			error
	)
	query := fmt.Sprintf("SELECT group_id, user_id FROM %s", mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var userGroup UserGroupModel
		err := rows.Scan(&userGroup.GroupId, &userGroup.UserId)
		if err, ok := err.(*pq.Error); ok {
			return nil, fmt.Errorf("%s", err.Detail)
		}
		userGroups = append(userGroups, userGroup)
	}
	return &userGroups, nil
}