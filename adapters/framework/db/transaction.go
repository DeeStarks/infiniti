package db

import (
	"fmt"
	"time"

	"github.com/deestarks/infiniti/utils"
	"github.com/lib/pq"
)

type (
	TransactionModel struct {
		Id					int			`json:"id"`
		UserId				int			`json:"user_id"`
		TransactionTypeId	int			`json:"transaction_type_id"`
		Amount				float64		`json:"amount"`
		SenderId			int			`json:"sender_id"`
		ReceiverId			int			`json:"receiver_id"`
		Remark				string		`json:"remark"`
		CreatedAt			time.Time	`json:"created_at"`
	}
	
	TransactionAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewTransactionAdapter() *TransactionAdapter {
	return &TransactionAdapter{
		adapter: adpt,
		tableName: "user_transactions",
	}
}

func (mAdapt *TransactionAdapter) Create(data map[string]interface{}) (*TransactionModel, error) {
	var transaction TransactionModel

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
		RETURNING id, user_id, transaction_type_id, amount, sender_id, reciever_id, remark, created_at
	`, mAdapt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(
		&transaction.Id, &transaction.UserId, &transaction.TransactionTypeId, &transaction.Amount,
		&transaction.SenderId, &transaction.ReceiverId, &transaction.Remark, &transaction.CreatedAt,
	)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &transaction, nil
}

func (mAdapt *TransactionAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*TransactionModel, error) {
	var (
		transaction 	TransactionModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, user_id, transaction_type_id, amount, sender_id, reciever_id, remark, created_at
	`, mAdapt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1)

	valArr = append(valArr, colValue)

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(
		&transaction.Id, &transaction.UserId, &transaction.TransactionTypeId, &transaction.Amount,
		&transaction.SenderId, &transaction.ReceiverId, &transaction.Remark, &transaction.CreatedAt,
	)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &transaction, nil
}

func (mAdapt *TransactionAdapter) Delete(colName string, value interface{}) (*TransactionModel, error) {
	var (
		transaction	TransactionModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, user_id, transaction_type_id, amount, sender_id, reciever_id, remark, created_at
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(
		&transaction.Id, &transaction.UserId, &transaction.TransactionTypeId, &transaction.Amount,
		&transaction.SenderId, &transaction.ReceiverId, &transaction.Remark, &transaction.CreatedAt,
	)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &transaction, nil
}

func (mAdapt *TransactionAdapter) Get(colName string, value interface{}) (*TransactionModel, error) {
	var (
		transaction	TransactionModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT id, user_id, transaction_type_id, amount, sender_id, reciever_id, remark, created_at
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(
		&transaction.Id, &transaction.UserId, &transaction.TransactionTypeId, &transaction.Amount,
		&transaction.SenderId, &transaction.ReceiverId, &transaction.Remark, &transaction.CreatedAt,
	)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	return &transaction, nil
}

func (mAdapt *TransactionAdapter) Filter(colName string, value interface{}) (*[]TransactionModel, error) {
	var (
		transactions	[]TransactionModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, user_id, transaction_type_id, amount, sender_id, reciever_id, remark, created_at FROM %s WHERE %s = $1", mAdapt.tableName, colName)
	rows, err := mAdapt.adapter.db.Query(query, value)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var transaction TransactionModel
		err := rows.Scan(
			&transaction.Id, &transaction.UserId, &transaction.TransactionTypeId, &transaction.Amount,
			&transaction.SenderId, &transaction.ReceiverId, &transaction.Remark, &transaction.CreatedAt,
		)
		if err, ok := err.(*pq.Error); ok {
			return nil, fmt.Errorf("%s", err.Detail)
		}
		transactions = append(transactions, transaction)
	}
	return &transactions, nil
}

func (mAdapt *TransactionAdapter) List() (*[]TransactionModel, error) {
	var (
		transactions	[]TransactionModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, user_id, transaction_type_id, amount, sender_id, reciever_id, remark, created_at FROM %s", mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var transaction TransactionModel
		err := rows.Scan(
			&transaction.Id, &transaction.UserId, &transaction.TransactionTypeId, &transaction.Amount,
			&transaction.SenderId, &transaction.ReceiverId, &transaction.Remark, &transaction.CreatedAt,
		)
		if err, ok := err.(*pq.Error); ok {
			return nil, fmt.Errorf("%s", err.Detail)
		}
		transactions = append(transactions, transaction)
	}
	return &transactions, nil
}