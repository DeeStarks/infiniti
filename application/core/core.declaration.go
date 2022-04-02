package core

type CoreApplication struct {}

func NewCoreApplication() CoreAppPort {
	return &CoreApplication{}
}