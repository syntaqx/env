//go:build !windows

package env

import (
	"testing"
)

func TestUnsetError(t *testing.T) {
	err := Unset("")
	assertNoError(t, err, "Unset")
}
