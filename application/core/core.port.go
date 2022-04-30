package core

type CoreAppPort interface {
	HashPassword(string) 				(string, error)
	ComparePassword(string, string) 	error
	MakeAccountNumber(int) 				string
	AccountNumberIsValid(string) 		bool
	GetIdFromAccountNumber(string) 		int
}