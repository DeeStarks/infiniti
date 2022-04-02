package db

import (
	"fmt"

	"github.com/deestarks/infiniti/utils"
	"github.com/lib/pq"
)

type (
	GroupModel struct {
		Id 		int 	`json:"id"`
		Name 	string 	`json:"name"`
	}
	
	GroupAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewGroupAdapter() *GroupAdapter {
	return &GroupAdapter{
		adapter: adpt,
		tableName: "groups",
	}
}

func (mAdapt *GroupAdapter) Create(data map[string]interface{}) (*GroupModel, error) {
	var group GroupModel

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
		RETURNING id, name
	`, mAdapt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&group.Id, &group.Name)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &group, nil
}

func (mAdapt *GroupAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*GroupModel, error) {
	var (
		group 	GroupModel
		valArr	[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, name
	`, mAdapt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1)

	valArr = append(valArr, colValue)

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&group.Id, &group.Name)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &group, nil
}

func (mAdapt *GroupAdapter) Delete(colName string, value interface{}) (*GroupModel, error) {
	var (
		group	GroupModel
		err		error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, name
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&group.Id, &group.Name)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &group, nil
}

// "Get" returns a single group
func (mAdapt *GroupAdapter) Get(colName string, value interface{}) (*GroupModel, error) {
	var (
		group	GroupModel
		err		error
	)

	query := fmt.Sprintf(`
		SELECT id, name
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, colName)
	// `, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&group.Id, &group.Name)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &group, nil
}

func (mAdapt *GroupAdapter) Filter(colName string, value interface{}) (*[]GroupModel, error) {
	var (
		groups	[]GroupModel
		err		error
	)
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE %s = $1", mAdapt.tableName, colName)
	rows, err := mAdapt.adapter.db.Query(query, value)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var group GroupModel
		err := rows.Scan(&group.Id, &group.Name)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
		groups = append(groups, group)
	}
	return &groups, nil
}

// "List" returns all groups
func (mAdapt *GroupAdapter) List() (*[]GroupModel, error) {
	var (
		groups	[]GroupModel
		err		error
	)
	query := fmt.Sprintf("SELECT id, name FROM %s", mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var group GroupModel
		err := rows.Scan(&group.Id, &group.Name)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
		groups = append(groups, group)
	}
	return &groups, nil
}