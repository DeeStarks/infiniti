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
	ReceiverId			int							`json:"reciever_id"`
	Remark				string						`json:"remark"`
	CreatedAt			time.Time					`json:"created_at"`
}

type TransactionsTransferResource struct {
	Id					int			`json:"id"`
	UserId				int			`json:"user_id"`
	Amount				float64		`json:"amount"`
	SenderId			int			`json:"sender_id"`
	ReceiverId			int			`json:"reciever_id"`
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
			"reciever_id": transaction["user_transactions__reciever_id"],
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
		"reciever_id": transaction["user_transactions__reciever_id"],
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

func (trans *Transactions) Deposit(data map[string]interface{}) (TransactionsDepositResource, error) {
	restricted := []string{"id", "user_id", "sender_id", "reciever_id", "created_at", "transaction_type_id"} // Restricted columns to create
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range restricted {
		delete(data, field)
	}
	requires := []string{"amount", "account_number"} // Required fields
	var notFound string
	for _, field := range requires {
		if _, ok := data[field]; !ok {
			notFound += ", "+field
		}
	}
	if notFound != "" {
		return TransactionsDepositResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("missing fields: %s", notFound[2:]),
		}
	}

	acctNo, ok := data["account_number"].(string)
	if !ok {
		return TransactionsDepositResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("\"account_number\" must be a string"),
		}
	}
	if !trans.corePort.AccountNumberIsValid(acctNo) {
		return  TransactionsDepositResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("invalid account number"),
		}
	}
	data["user_id"] = trans.corePort.GetIdFromAccountNumber(acctNo)
	delete(data, "account_number") // Account number won't be used anymore

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

	userId := data["user_id"]
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
		fmt.Println(err.Error())
		return TransactionsDepositResource{}, &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: err,
		}
	}

	var res TransactionsDepositResource
	jsonRes, _ := json.Marshal(transRes)
	json.Unmarshal(jsonRes, &res)
	fmt.Println(res)
	return res, nil
}

func (trans *Transactions) Withdrawal(data map[string]interface{}) (TransactionsWithdrawalResource, error) {
	restricted := []string{"id", "sender_id", "reciever_id", "created_at", "transaction_type_id"} // Restricted columns to create
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range restricted {
		delete(data, field)
	}

	// Allow the use of both user_id and account_number
	// If user_id is passed, account_number is ignored
	if _, ok := data["user_id"]; ok {
		delete(data, "account_number")
	}

	if _, ok := data["account_number"]; ok {
		// Retrieve user_id from account_number
		acctNo, ok := data["account_number"].(string)
		if isValid := trans.corePort.AccountNumberIsValid(acctNo); !isValid || !ok {
			return TransactionsWithdrawalResource{}, &utils.RequestError{
				Code: http.StatusBadRequest,
				Err: fmt.Errorf("'account_number' not valid"),
			}
		}
		data["user_id"] = trans.corePort.GetIdFromAccountNumber(acctNo)
		delete(data, "account_number")
	}

	requires := []string{"amount", "user_id"} // Required fields
	var notFound string
	for _, field := range requires {
		if _, ok := data[field]; !ok {
			if field == "user_id" { 
				// Make sure either user_id or account_number is passed
				notFound += ", "+field+" or account_number"
			} else {
				notFound += ", "+field
			}
		}
	}
	if notFound != "" {
		return TransactionsWithdrawalResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("missing fields: %s", notFound[2:]),
		}
	}

	// Validate amount
	amount, ok := data["amount"].(float64)
	if !ok {
		return TransactionsWithdrawalResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("\"amount\" must be a floating point number"),
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

	// Confirm user's account existence
	userId := data["user_id"]

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
	_, err = acctAdapter.Update("user_id", userId, map[string]interface{}{
		"balance": userAcct.Balance-amount,
	})
	if err != nil {
		return TransactionsWithdrawalResource{}, &utils.RequestError{
			Code: http.StatusInternalServerError,
			Err: err,
		}
	}

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

func (trans *Transactions) Transfer(data map[string]interface{}) (TransactionsTransferResource, error) {
	restricted := []string{"id", "user_id", "created_at", "reciever_id", "transaction_type_id"} // Restricted columns to create
	// Check if the data contains restricted columns
	// then remove them from the data
	for _, field := range restricted {
		delete(data, field)
	}

	requires := []string{"amount", "sender_id", "recipient_account_number"} // Required fields
	var notFound string
	for _, field := range requires {
		if _, ok := data[field]; !ok {
			notFound += ", "+field
		}
	}
	if notFound != "" {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("missing fields: %s", notFound[2:]),
		}
	}

	// Validate amount
	amount, ok := data["amount"].(float64)
	if !ok {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("\"amount\" must be a floating point number"),
		}
	}

	// Confirm receiver's account number validity
	receiverAcctNo, ok := data["recipient_account_number"].(string)
	if !ok {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("\"recipient_account_number\" must be a string"),
		}
	}
	if !trans.corePort.AccountNumberIsValid(receiverAcctNo) {
		return  TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("recipient's account number is not valid"),
		}
	}
	data["reciever_id"] = trans.corePort.GetIdFromAccountNumber(receiverAcctNo)
	delete(data, "recipient_account_number") // Account number won't be used anymore

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
	userId := data["sender_id"]
	data["user_id"] = userId

	senderAcct, err := accountAdapter.Get("user_id", userId)
	if err != nil {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("cannot transfer - sender unknown"),
		}
	}
	// Compare balance
	if amount > senderAcct.Balance {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("insufficient balance"),
		}
	}

	// Make sure receiver exists
	reciever_id := data["reciever_id"]
	receiverAcct, err := accountAdapter.Get("user_id", reciever_id)
	if err != nil {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("cannot transfer - receiver unknown"),
		}
	}

	// Exchange currency
	currencyAdapter := trans.dbPort.NewCurrencyAdapter()
	// Sender's currency
	senderCurrency, err := currencyAdapter.Get("id", senderAcct.CurrencyId)
	if err != nil {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("sender's currency unknown"),
		}
	}
	// Receiver's currency
	receiverCurrency, err := currencyAdapter.Get("id", receiverAcct.CurrencyId)
	if err != nil {
		return TransactionsTransferResource{}, &utils.RequestError{
			Code: http.StatusBadRequest,
			Err: fmt.Errorf("receiver's currency unknown"),
		}
	}
	// Exchange amount
	exchangedAmount := trans.corePort.ConvertCurrency(amount, senderCurrency.ConversionRateToUSD, receiverCurrency.ConversionRateToUSD)

	// Make transfer
	// Deduct the amount from sender's account balance
	accountAdapter.Update("user_id", userId, map[string]interface{}{
		"balance": senderAcct.Balance-amount,
	})
	// Add the converted amount to receiver's account balance
	accountAdapter.Update("user_id", reciever_id, map[string]interface{}{
		"balance": receiverAcct.Balance+exchangedAmount,
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