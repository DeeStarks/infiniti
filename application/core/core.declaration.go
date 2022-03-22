package core

import (
	"github.com/deestarks/infiniti/ports"
)

type CoreApplication struct {}

func NewCoreApplication() ports.CoreAppPort {
	return CoreApplication{}
}