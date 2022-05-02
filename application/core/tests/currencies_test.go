package core_tests

import (
	"testing"

	"github.com/deestarks/infiniti/application/core"
)

func TestConvertCurrency(t *testing.T) {
	cases := []struct {
		amount, fromConversionRate, toConversionRate, expected float64
	}{
		{amount: 1000.20, fromConversionRate: 1.0, toConversionRate: 1.0, expected: 1000.20},
		{amount: 1000.20, fromConversionRate: 1.0, toConversionRate: 0.5, expected: 500.10},
		{amount: 1000.20, fromConversionRate: 0.5, toConversionRate: 1.0, expected: 2000.40},
		{amount: 1000.20, fromConversionRate: 414.0, toConversionRate: 1.0, expected: 2.42},
		{amount: 1000.20, fromConversionRate: 1.0, toConversionRate: 414.0, expected: 414082.80},
	}

	corePort := core.NewCoreApplication()
	for _, c := range cases {
		if actual := corePort.ConvertCurrency(c.amount, c.fromConversionRate, c.toConversionRate); actual - c.expected > 0.1 { // Allow for rounding errors
			t.Errorf("ConvertCurrency(%f, %f, %f) == %f, expected %f", c.amount, c.fromConversionRate, c.toConversionRate, actual, c.expected)
		}
	}
}