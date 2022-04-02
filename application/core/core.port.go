package core

type CoreAppPort interface {
	HashPassword(string) 				(string, error)
	ComparePassword(string, string) 	(error)
}