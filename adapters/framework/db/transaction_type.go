package db

import (
	"fmt"
	"log"

	"github.com/deestarks/infiniti/utils"
)

type (
	TransactionTypeModel struct {
		Id		int		`json:"id"`
		Name	string	`json:"name"`
	}
	
	TransactionTypeAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewTransactionTypeAdapter() *TransactionTypeAdapter {
	return &TransactionTypeAdapter{
		adapter: adpt,
		tableName: "transaction_types",
	}
}

func (mAdapt *TransactionTypeAdapter) Create(data map[string]interface{}) (*TransactionTypeModel, error) {
	var transactionType TransactionTypeModel

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

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&transactionType.Id, &transactionType.Name)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &transactionType, nil
}

func (mAdapt *TransactionTypeAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*TransactionTypeModel, error) {
	var (
		transactionType 	TransactionTypeModel
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

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(&transactionType.Id, &transactionType.Name)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &transactionType, nil
}

func (mAdapt *TransactionTypeAdapter) Delete(colName string, value interface{}) (*TransactionTypeModel, error) {
	var (
		transactionType	TransactionTypeModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, name
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&transactionType.Id, &transactionType.Name)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &transactionType, nil
}

func (mAdapt *TransactionTypeAdapter) Get(colName string, value interface{}) (*TransactionTypeModel, error) {
	var (
		transactionType	TransactionTypeModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT id, name
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(&transactionType.Id, &transactionType.Name)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &transactionType, nil
}

func (mAdapt *TransactionTypeAdapter) Filter(colName string, value interface{}) (*[]TransactionTypeModel, error) {
	var (
		transactionTypes	[]TransactionTypeModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE %s = $1", mAdapt.tableName, colName)
	rows, err := mAdapt.adapter.db.Query(query, value)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transactionType TransactionTypeModel
		err := rows.Scan(&transactionType.Id, &transactionType.Name)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		transactionTypes = append(transactionTypes, transactionType)
	}
	return &transactionTypes, nil
}

func (mAdapt *TransactionTypeAdapter) List() (*[]TransactionTypeModel, error) {
	var (
		transactionTypes	[]TransactionTypeModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, name FROM %s", mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transactionType TransactionTypeModel
		err := rows.Scan(&transactionType.Id, &transactionType.Name)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		transactionTypes = append(transactionTypes, transactionType)
	}
	return &transactionTypes, nil
}