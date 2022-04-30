package core

import "strconv"

const BankNo = 1220000000

func (core *CoreApplication) MakeAccountNumber(id int) string {
	accountNo := strconv.Itoa(BankNo+id) 
	return accountNo
}

func (core *CoreApplication) AccountNumberIsValid(accountNo string) bool {
	accountNoInt, _ := strconv.Atoi(accountNo)
	return accountNoInt > BankNo && accountNoInt <= BankNo+9999999
}

func (core *CoreApplication) GetIdFromAccountNumber(accountNo string) int {
	accountNoInt, _ := strconv.Atoi(accountNo)
	if accountNoInt > BankNo && accountNoInt <= BankNo+9999999 {
		return accountNoInt - BankNo
	}
	return 0
}