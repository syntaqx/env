//go:build windows
// +build windows

package env

import (
	"testing"
)

func TestUnsetError(t *testing.T) {
	err := Unset("")
	t.Fail()
	assertError(t, err, "Unset")
}
