package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deestarks/infiniti/adapters/framework/db"
	"github.com/deestarks/infiniti/application/core"
	"github.com/deestarks/infiniti/utils"
)


type Transactions struct {
	dbPort 		db.DBPort
	corePort 	core.CoreAppPort
}

type TransactionsResource struct {
	Id					int							`json:"id"`
	UserId				int							`json:"user_id"`
	TransactionType		TransactionTypesResource	`json:"transaction_type_id"`
	Amount				float64						`json:"amount"`
	SenderId			int							`json:"sender_id"`
	ReceiverId			int							`json:"receiver_id"`
	Remark				string						`json:"remark"`
	CreatedAt			time.Time					`json:"created_at"`
}

type TransactionsTransferResource struct {
	Id					int			`json:"id"`
	UserId				int			`json:"user_id"`
	Amount				float64		`json:"amount"`
	SenderId			int			`json:"sender_id"`
	ReceiverId			int			`json:"receiver_id"`
	Remark				string		`json:"remark"`
	CreatedAt			time.Time	`json:"created_at"`
}

type TransactionsWithdrawalResource struct {
	Id					int			`json:"id"`
	UserId				int			`json:"user_id"`
	Amount				float64		`json:"amount"`
	CreatedAt			time.Time	`json:"created_at"`
}

type TransactionsDepositResource struct {
	Id					int			`json:"id"`
	UserId				int			`json:"user_id"`
	Amount				float64		`json:"amount"`
	CreatedAt			time.Time	`json:"created_at"`
}

func (service *Service) NewTransactionsService() *Transactions {
	return &Transactions{
		dbPort: 	service.dbPort,
		corePort:	service.corePort,
	}
}

func (trans *Transactions) ListTransactions(userId int) ([]TransactionsResource, error) {
	adapter := trans.dbPort.NewTransactionAdapter()
	conditions := map[string]interface{}{
		"user_transactions.user_id": userId,
	}
	selector := adapter.NewTransactionCustomSelector(conditions, "user_transactions.id", true).
		Join("transaction_types", "id", "user_transactions", "transaction_type_id", []string{"id", "name"})
	transactions, err := selector.Query()
	if err != nil {
		return nil, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err:	err,
		}
	}

	var res []TransactionsResource
	for _, transaction := range transactions {
		transactionData := map[string]interface{} {
			"id": transaction["user_transactions__id"],
			"user_id": transaction["user_transactions__user_id"],
			"amount": transaction["user_transactions__amount"],
			"sender_id": transaction["user_transactions__sender_id"],
			"receiver_id": transaction["user_transactions__receiver_id"],
			"remark": transaction["user_transactions__remark"],
			"created_at": transaction["user_transactions__created_at"],
		}

		transactionType := map[string]interface{} {
			"id": transaction["transaction_types__id"],
			"name": transaction["transaction_types__name"],
		}

		var trasRes TransactionsResource
		tJson, _ := json.Marshal(transactionData)
		json.Unmarshal(tJson, &trasRes)

		var ttRes TransactionTypesResource
		ttJson, _ := json.Marshal(transactionType)
		json.Unmarshal(ttJson, &ttRes)

		trasRes.TransactionType = ttRes
		res = append(res, trasRes)

	}
	return res, nil
}

func (trans *Transactions) GetTransaction(userId int, colName string, colValue interface{}) (TransactionsResource, error) {
	adapter := trans.dbPort.NewTransactionAdapter()
	conditions := map[string]interface{}{
		"user_transactions.user_id": userId,
		"user_transactions."+colName: colValue,
	}
	selector := adapter.NewTransactionCustomSelector(conditions, "user_transactions.id", true).
		Join("transaction_types", "id", "user_transactions", "transaction_type_id", []string{"id", "name"})
	data, err := selector.Query()
	if err != nil {
		return TransactionsResource{}, &utils.RequestError{
			Code:	http.StatusInternalServerError,
			Err:	err,
		}
	}

	if len(data) < 1 {
		return TransactionsResource{}, &utils.RequestError{
			Code:	http.StatusBadRequest,
			Err: 	fmt.Errorf("user transaction not found"),
		}
	}
	transaction := data[0] // There should be one if valid values are passed to "colName" and "colValue"

	transactionData := map[string]interface{} {
		"id": transaction["user_transactions__id"],
		"user_id": transaction["user_transactions__user_id"],
		"amount": transaction["user_transactions__amount"],
		"sender_id": transaction["user_transactions__sender_id"],
		"receiver_id": transaction["user_transactions__receiver_id"],
		"remark": transaction["user_transactions__remark"],
		"created_at": transaction["user_transactions__created_at"],
	}

	transactionType := map[string]interface{} {
		"id": transaction["transaction_types__id"],
		"name": transaction["transaction_types__name"],
	}

	var trasRes TransactionsResource
	tJson, _ := json.Marshal(transactionData)
	json.Unmarshal(tJson, &trasRes)

	var ttRes TransactionTypesResource
	ttJson, _ := json.Marshal(transactionType)
	json.Unmarshal(ttJson, &ttRes)

	trasRes.TransactionType = ttRes
	return trasRes, nil
}

func (trans *Transactions) Deposit(userId int, data map[string]interface{}) (TransactionsDepositResource, error) {
	rsrts := []string{"id", "user_id", "sender_id", "receiver_id", "created_at", "transaction_type_id"} // Restricted columns to create
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range rsrts {
		delete(data, field)
	}

	// Add the userId
	data["user_id"] = userId

	requires := []string{"amount"} // Required fields
	var notFound string
	for _, field := range requires {
		if _, ok := data[field]; !ok {
			notFound += " ,"+field
		}
	}
	if notFound != "" {
		return TransactionsDepositResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("missing fields: %s", notFound[2:]),
		}
	}

	// Retrieve transfer type and add to data
	transTypeAdapter := trans.dbPort.NewTransactionTypeAdapter()
	tt, err := transTypeAdapter.Get("name", "deposit")
	if err != nil {
		return TransactionsDepositResource{}, &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: fmt.Errorf("\"deposit\" transaction type not found"),
		}
	}
	data["transaction_type_id"] = tt.Id

	// Validate amount
	amount, ok := data["amount"].(float64)
	if !ok {
		return TransactionsDepositResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("\"amount\" must be a floating point number"),
		}
	}

	// Confirm user's account existence
	acctAdapter := trans.dbPort.NewAccountAdapter()
	userAcct, err := acctAdapter.Get("user_id", userId)
	if err != nil {
		return TransactionsDepositResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("cannot deposit - user account not found"),
		}
	}

	// Add to balance
	acctAdapter.Update("user_id", userId, map[string]interface{}{
		"balance": userAcct.Balance+amount,
	})

	// Create transaction data
	transactionAdapter := trans.dbPort.NewTransactionAdapter()
	transRes, err := transactionAdapter.Create(data)
	if err != nil {
		return TransactionsDepositResource{}, &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: err,
		}
	}

	var res TransactionsDepositResource
	jsonRes, _ := json.Marshal(transRes)
	json.Unmarshal(jsonRes, &res)
	return res, nil
}

func (trans *Transactions) Withdrawal(userId int, data map[string]interface{}) (TransactionsWithdrawalResource, error) {
	rsrts := []string{"id", "user_id", "sender_id", "receiver_id", "created_at", "transaction_type_id"} // Restricted columns to create
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range rsrts {
		delete(data, field)
	}

	// Add the userId
	data["user_id"] = userId

	requires := []string{"amount"} // Required fields
	var notFound string
	for _, field := range requires {
		if _, ok := data[field]; !ok {
			notFound += " ,"+field
		}
	}
	if notFound != "" {
		return TransactionsWithdrawalResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("missing fields: %s", notFound[2:]),
		}
	}

	// Retrieve transfer type and add to data
	transTypeAdapter := trans.dbPort.NewTransactionTypeAdapter()
	tt, err := transTypeAdapter.Get("name", "withdrawal")
	if err != nil {
		return TransactionsWithdrawalResource{}, &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: fmt.Errorf("\"withdrawal\" transaction type not found"),
		}
	}
	data["transaction_type_id"] = tt.Id

	// Validate amount
	amount, ok := data["amount"].(float64)
	if !ok {
		return TransactionsWithdrawalResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("\"amount\" must be a floating point number"),
		}
	}

	// Confirm user's account existence
	acctAdapter := trans.dbPort.NewAccountAdapter()
	userAcct, err := acctAdapter.Get("user_id", userId)
	if err != nil {
		return TransactionsWithdrawalResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("cannot withdraw - user account not found"),
		}
	}

	// Confirm sufficiency
	if amount > userAcct.Balance {
		return TransactionsWithdrawalResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("insufficient balance"),
		}
	}

	// Withdraw from balance
	acctAdapter.Update("user_id", userId, map[string]interface{}{
		"balance": userAcct.Balance-amount,
	})

	// Create transaction data
	transactionAdapter := trans.dbPort.NewTransactionAdapter()
	transRes, err := transactionAdapter.Create(data)
	if err != nil {
		return TransactionsWithdrawalResource{}, &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: err,
		}
	}

	var res TransactionsWithdrawalResource
	jsonRes, _ := json.Marshal(transRes)
	json.Unmarshal(jsonRes, &res)
	return res, nil
}

func (trans *Transactions) Transfer(userId int, data map[string]interface{}) (TransactionsTransferResource, error) {
	rsrts := []string{"id", "user_id", "sender_id", "created_at", "transaction_type_id"} // Restricted columns to create
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range rsrts {
		delete(data, field)
	}

	// Add the userId
	data["user_id"] = userId
	data["sender_id"] = userId

	requires := []string{"amount", "receiver_id"} // Required fields
	var notFound string
	for _, field := range requires {
		if _, ok := data[field]; !ok {
			notFound += " ,"+field
		}
	}
	if notFound != "" {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("missing fields: %s", notFound[2:]),
		}
	}

	// Retrieve transfer type and add to data
	transTypeAdapter := trans.dbPort.NewTransactionTypeAdapter()
	tt, err := transTypeAdapter.Get("name", "transfer")
	if err != nil {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: fmt.Errorf("\"transfer\" transaction type not found"),
		}
	}
	data["transaction_type_id"] = tt.Id

	// Initialize the account database adapter
	accountAdapter := trans.dbPort.NewAccountAdapter() 
	// Get sender and check balance
	senderAcct, err := accountAdapter.Get("user_id", userId)
	if err != nil {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("cannot transfer - sender unknown"),
		}
	}
	// Compare balance
	amount, ok := data["amount"].(float64)
	if !ok {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("\"amount\" must be a floating point number"),
		}
	}
	if amount > senderAcct.Balance {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("insufficient balance"),
		}
	}

	// Make sure receiver exists
	receiver_id := data["receiver_id"]
	receiverAcct, err := accountAdapter.Get("user_id", receiver_id)
	if err != nil {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("cannot transfer - receiver unknown"),
		}
	}

	// Make transfer
	accountAdapter.Update("user_id", userId, map[string]interface{}{
		"balance": senderAcct.Balance-amount,
	})
	accountAdapter.Update("user_id", receiver_id, map[string]interface{}{
		"balance": receiverAcct.Balance+amount,
	})

	// Save to transactions
	transactionAdapter := trans.dbPort.NewTransactionAdapter()
	transRes, err := transactionAdapter.Create(data)
	if err != nil {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: err,
		}
	}

	var res TransactionsTransferResource
	jsonRes, _ := json.Marshal(transRes)
	json.Unmarshal(jsonRes, &res)
	return res, nil
}