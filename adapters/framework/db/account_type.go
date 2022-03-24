package db

import (
	"fmt"

	"github.com/deestarks/infiniti/utils"
)

type (
	AccountTypeModel struct {
		Id		int		`json:"id"`
		Name	string	`json:"name"`
	}
	
	AccountTypeAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewAccountTypeAdapter() *AccountTypeAdapter {
	return &AccountTypeAdapter{
		adapter: adpt,
		tableName: "account_types",
	}
}

func (mAdapt *AccountTypeAdapter) Create(data map[string]interface{}) (*AccountTypeModel, error) {
	var accountType AccountTypeModel

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

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&accountType.Id, &accountType.Name)
	if err != nil {
		return nil, err
	}
	return &accountType, nil
}

func (mAdapt *AccountTypeAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*AccountTypeModel, error) {
	var (
		accountType 	AccountTypeModel
		valArr		[]interface{}
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

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&accountType.Id, &accountType.Name)
	if err != nil {
		return nil, err
	}
	return &accountType, nil
}

func (mAdapt *AccountTypeAdapter) Delete(colName string, value interface{}) (*AccountTypeModel, error) {
	var (
		accountType	AccountTypeModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, name
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&accountType.Id, &accountType.Name)
	if err != nil {
		return nil, err
	}
	return &accountType, nil
}

func (mAdapt *AccountTypeAdapter) Get(colName string, value interface{}) (*AccountTypeModel, error) {
	var (
		accountType	AccountTypeModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT id, name
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&accountType.Id, &accountType.Name)
	if err != nil {
		return nil, err
	}
	return &accountType, nil
}

func (mAdapt *AccountTypeAdapter) Filter(colName string, value interface{}) (*[]AccountTypeModel, error) {
	var (
		accountTypes	[]AccountTypeModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE %s = $1", mAdapt.tableName, colName)
	rows, err := mAdapt.adapter.db.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var accountType AccountTypeModel
		err := rows.Scan(&accountType.Id, &accountType.Name)
		if err != nil {
			return nil, err
		}
		accountTypes = append(accountTypes, accountType)
	}
	return &accountTypes, nil
}

func (mAdapt *AccountTypeAdapter) List() (*[]AccountTypeModel, error) {
	var (
		accountTypes	[]AccountTypeModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, name FROM %s", mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var accountType AccountTypeModel
		err := rows.Scan(&accountType.Id, &accountType.Name)
		if err != nil {
			return nil, err
		}
		accountTypes = append(accountTypes, accountType)
	}
	return &accountTypes, nil
}