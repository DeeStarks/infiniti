package core_tests

import (
	"testing"
	"github.com/deestarks/infiniti/application/core"
)

func TestHashComparePassword(t *testing.T) {
	cases := []struct {
		password string
	}{
		{"password@*#"},
		{"password1^&@"},
		{"password2(&$"},
		{"password3@)*"},
		{"password4!)@*"},
	}

	coreApp := core.NewCoreApplication()

	for _, c := range cases {
		hash, err := coreApp.HashPassword(c.password)
		if err != nil {
			t.Error(err)
		}
		if hash == c.password {
			t.Error("Password and hash are the same")
		}

		err = coreApp.ComparePassword(hash, c.password)
		if err != nil {
			t.Error(err)
		}
	}
}