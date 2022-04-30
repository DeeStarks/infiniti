package core

import "strconv"

func (core *CoreApplication) MakeAccountNumber(id int) string {
	bankNo := 1220000000
	accountNo := strconv.Itoa(bankNo+id) 
	return accountNo
}

func (core *CoreApplication) AccountNumberIsValid(accountNo string) bool {
	bankNo := 1220000000
	accountNoInt, _ := strconv.Atoi(accountNo)
	return accountNoInt > bankNo
}

func (core *CoreApplication) GetIdFromAccountNumber(accountNo string) int {
	bankNo := 1220000000
	accountNoInt, _ := strconv.Atoi(accountNo)
	return accountNoInt - bankNo
}