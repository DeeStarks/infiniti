package core_tests

import (
	"testing"

	"github.com/deestarks/infiniti/application/core"
)

func TestMakeAccountNumber(t *testing.T) {
	cases := []struct {
		id 			int
		expected	string
	}{
		{12, "1220000012"},
		{20, "1220000020"},
		{312, "1220000312"},
		{420, "1220000420"},
		{500, "1220000500"},
	}

	coreApp := core.NewCoreApplication()
	for _, c := range cases {
		accountNo := coreApp.MakeAccountNumber(c.id)
		if accountNo != c.expected {
			t.Errorf("Expected %s but got %s", c.expected, accountNo)
		}
	}
}

func TestAccountNumberIsValid(t *testing.T) {
	// Expected between 1220000000 and 1229999999
	cases := []struct {
		accountNo	string
		expected	bool
	}{
		{"1220000012", true},
		{"1220000120", true},
		{"1220000200", true},
		{"1220000300", true},
		{"1220000400", true},
		{"1220000500", true},
		{"1234567890", false},
		{"1229999999", true},
		{"1220000000", false},
		{"1123123312", false},
		{"123222323", false},
	}

	coreApp := core.NewCoreApplication()
	for _, c := range cases {
		valid := coreApp.AccountNumberIsValid(c.accountNo)
		if valid != c.expected {
			t.Errorf("Expected %t but got %t", c.expected, valid)
		}
	}
}

func TestGetIdFromAccountNumber(t *testing.T) {
	cases := []struct {
		accountNo	string
		expected	int
	}{
		{"1220000012", 12},
		{"1220000120", 120},
		{"1220000200", 200},
		{"1220000300", 300},
		{"1220000400", 400},
		{"1220000500", 500},
		{"1234567890", 0},
		{"1229999999", 9999999},
		{"1220000000", 0},
		{"1123123312", 0},
		{"123222323", 0},
	}

	coreApp := core.NewCoreApplication()
	for _, c := range cases {
		id := coreApp.GetIdFromAccountNumber(c.accountNo)
		if id != c.expected {
			t.Errorf("Expected %d but got %d", c.expected, id)
		}
	}
}