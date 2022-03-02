package db

import (
	"fmt"
	"log"

	"github.com/deestarks/infiniti/utils"
)

type (
	AccountModel struct {
		Id				int		`json:"id"`
		UserId 			int		`json:"user_id"`
		AccountTypeId 	int		`json:"account_type_id"`
		AccountNumber 	uint64	`json:"account_number"`
		Balance 		float64	`json:"balance"`
		CurrencyId 		int		`json:"currency_id"`
	}
	
	AccountAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewAccountAdapter() *AccountAdapter {
	return &AccountAdapter{
		adapter: adpt,
		tableName: "user_accounts",
	}
}

func (mAdapt *AccountAdapter) Create(data map[string]interface{}) (*AccountModel, error) {
	var account AccountModel

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
		RETURNING id, user_id, account_type_id, account_number, balance, currency_id
	`, mAdapt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(
		&account.Id, &account.UserId, &account.AccountTypeId,
		&account.AccountNumber, &account.Balance, &account.CurrencyId,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &account, nil
}

func (mAdapt *AccountAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (*AccountModel, error) {
	var (
		account 	AccountModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, user_id, account_type_id, account_number, balance, currency_id
	`, mAdapt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1)

	valArr = append(valArr, colValue)

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(
		&account.Id, &account.UserId, &account.AccountTypeId,
		&account.AccountNumber, &account.Balance, &account.CurrencyId,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &account, nil
}

func (mAdapt *AccountAdapter) Delete(colName string, value interface{}) (*AccountModel, error) {
	var (
		account	AccountModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, user_id, account_type_id, account_number, balance, currency_id
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(
		&account.Id, &account.UserId, &account.AccountTypeId,
		&account.AccountNumber, &account.Balance, &account.CurrencyId,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &account, nil
}

func (mAdapt *AccountAdapter) Get(colName string, value interface{}) (*AccountModel, error) {
	var (
		account	AccountModel
		err			error
	)

	query := fmt.Sprintf(`
		SELECT id, user_id, account_type_id, account_number, balance, currency_id
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(
		&account.Id, &account.UserId, &account.AccountTypeId,
		&account.AccountNumber, &account.Balance, &account.CurrencyId,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &account, nil
}

func (mAdapt *AccountAdapter) Filter(colName string, value interface{}) (*[]AccountModel, error) {
	var (
		accounts	[]AccountModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, user_id, account_type_id, account_number, balance, currency_id FROM %s WHERE %s = $1", mAdapt.tableName, colName)
	rows, err := mAdapt.adapter.db.Query(query, value)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account AccountModel
		err := rows.Scan(
			&account.Id, &account.UserId, &account.AccountTypeId,
			&account.AccountNumber, &account.Balance, &account.CurrencyId,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return &accounts, nil
}

func (mAdapt *AccountAdapter) List() (*[]AccountModel, error) {
	var (
		accounts	[]AccountModel
		err			error
	)
	query := fmt.Sprintf("SELECT id, user_id, account_type_id, account_number, balance, currency_id FROM %s", mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account AccountModel
		err := rows.Scan(
			&account.Id, &account.UserId, &account.AccountTypeId,
			&account.AccountNumber, &account.Balance, &account.CurrencyId,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return &accounts, nil
}