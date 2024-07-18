package env

import (
	"reflect"
	"testing"
)

func assertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v %v", err, msgAndArgs)
	}
}

func assertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error, got nil %v", msgAndArgs...)
	}
}

func assertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v %v", expected, actual, msgAndArgs)
	}
}

func setEnvForTest(t *testing.T, name string, value string) {
	t.Helper()

	err := Set(name, value)
	assertNoError(t, err, "Set")

	t.Cleanup(func() {
		err = Unset(name) // clean after the test
		assertNoError(t, err, "Unset")
	})
}

func TestSetUnset(t *testing.T) {
	key, value := "TEST_KEY", "TEST_VALUE"
	err := Set(key, value)
	assertNoError(t, err, "Set")

	actual, err := Get(key)
	assertNoError(t, err, "Get")
	assertEqual(t, value, actual, "Get")

	err = Unset(key)
	assertNoError(t, err, "Unset")

	if _, ok := Lookup(key); ok {
		t.Errorf("Lookup: expected %s to be unset", key)
	}
}

func TestGetUnset(t *testing.T) {
	key := "TEST_UNSET"
	_, err := Get(key)
	assertError(t, err, "Get unset variable")
}

func TestRequire(t *testing.T) {
	key := "TEST_REQUIRED"

	err := Require(key)
	assertError(t, err, "Require")

	err = Set(key, "value")
	assertNoError(t, err, "Set")

	err = Require(key)
	assertNoError(t, err, "Require")

	err = Unset(key)
	assertNoError(t, err, "Unset")
}
