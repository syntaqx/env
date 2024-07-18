package env

import (
	"testing"
)

func TestGet(t *testing.T) {
	setEnvForTest(t, "TEST_GET", "value")

	value, err := Get("TEST_GET")
	assertNoError(t, err, "Get TEST_GET")
	assertEqual(t, "value", value, "Get TEST_GET")
}

func TestGetWithFallback(t *testing.T) {
	setEnvForTest(t, "TEST_GET_WITH_FALLBACK", "value")

	value := GetWithFallback("TEST_GET_WITH_FALLBACK", "fallback")
	assertEqual(t, "value", value, "GetWithFallback TEST_GET_WITH_FALLBACK")

	Unset("TEST_GET_WITH_FALLBACK")

	value = GetWithFallback("TEST_GET_WITH_FALLBACK", "fallback")
	assertEqual(t, "fallback", value, "GetWithFallback TEST_GET_WITH_FALLBACK")
}

func TestGetBool(t *testing.T) {
	setEnvForTest(t, "TEST_BOOL", "true")

	value, err := GetBool("TEST_BOOL")
	assertNoError(t, err, "GetBool TEST_BOOL")
	assertEqual(t, true, value, "GetBool TEST_BOOL")

	setEnvForTest(t, "TEST_BOOL", "invalid")

	_, err = GetBool("TEST_BOOL")
	assertError(t, err, "GetBool TEST_BOOL invalid")
}

func TestGetBoolWithFallback(t *testing.T) {
	setEnvForTest(t, "TEST_BOOL_WITH_FALLBACK", "true")

	value, err := GetBoolWithFallback("TEST_BOOL_WITH_FALLBACK", false)
	assertNoError(t, err, "GetBoolWithFallback TEST_BOOL_WITH_FALLBACK")
	assertEqual(t, true, value, "GetBoolWithFallback TEST_BOOL_WITH_FALLBACK")

	Unset("TEST_BOOL_WITH_FALLBACK")

	value, err = GetBoolWithFallback("TEST_BOOL_WITH_FALLBACK", false)
	assertNoError(t, err, "GetBoolWithFallback TEST_BOOL_WITH_FALLBACK")
	assertEqual(t, false, value, "GetBoolWithFallback TEST_BOOL_WITH_FALLBACK")
}

func TestGetInt(t *testing.T) {
	setEnvForTest(t, "TEST_INT", "42")

	value, err := GetInt("TEST_INT")
	assertNoError(t, err, "GetInt TEST_INT")
	assertEqual(t, 42, value, "GetInt TEST_INT")

	setEnvForTest(t, "TEST_INT", "invalid")

	_, err = GetInt("TEST_INT")
	assertError(t, err, "GetInt TEST_INT invalid")
}

func TestGetIntWithFallback(t *testing.T) {
	setEnvForTest(t, "TEST_INT_WITH_FALLBACK", "42")

	value, err := GetIntWithFallback("TEST_INT_WITH_FALLBACK", 10)
	assertNoError(t, err, "GetIntWithFallback TEST_INT_WITH_FALLBACK")
	assertEqual(t, 42, value, "GetIntWithFallback TEST_INT_WITH_FALLBACK")

	Unset("TEST_INT_WITH_FALLBACK")

	value, err = GetIntWithFallback("TEST_INT_WITH_FALLBACK", 10)
	assertNoError(t, err, "GetIntWithFallback TEST_INT_WITH_FALLBACK")
	assertEqual(t, 10, value, "GetIntWithFallback TEST_INT_WITH_FALLBACK")
}

func TestGetFloat(t *testing.T) {
	setEnvForTest(t, "TEST_FLOAT", "42.42")

	value, err := GetFloat("TEST_FLOAT")
	assertNoError(t, err, "GetFloat TEST_FLOAT")
	assertEqual(t, 42.42, value, "GetFloat TEST_FLOAT")

	setEnvForTest(t, "TEST_FLOAT", "invalid")

	_, err = GetFloat("TEST_FLOAT")
	assertError(t, err, "GetFloat TEST_FLOAT invalid")
}

func TestGetFloatWithFallback(t *testing.T) {
	setEnvForTest(t, "TEST_FLOAT_WITH_FALLBACK", "42.42")

	value, err := GetFloatWithFallback("TEST_FLOAT_WITH_FALLBACK", 10.1)
	assertNoError(t, err, "GetFloatWithFallback TEST_FLOAT_WITH_FALLBACK")
	assertEqual(t, 42.42, value, "GetFloatWithFallback TEST_FLOAT_WITH_FALLBACK")

	Unset("TEST_FLOAT_WITH_FALLBACK")

	value, err = GetFloatWithFallback("TEST_FLOAT_WITH_FALLBACK", 10.1)
	assertNoError(t, err, "GetFloatWithFallback TEST_FLOAT_WITH_FALLBACK")
	assertEqual(t, 10.1, value, "GetFloatWithFallback TEST_FLOAT_WITH_FALLBACK")
}

func TestGetStringSlice(t *testing.T) {
	setEnvForTest(t, "TEST_STRING_SLICE", "value1,value2")

	value, err := GetStringSlice("TEST_STRING_SLICE")
	assertNoError(t, err, "GetStringSlice TEST_STRING_SLICE")
	assertEqual(t, []string{"value1", "value2"}, value, "GetStringSlice TEST_STRING_SLICE")
}

func TestGetStringSliceWithFallback(t *testing.T) {
	setEnvForTest(t, "TEST_STRING_SLICE_WITH_FALLBACK", "value1,value2")

	value, err := GetStringSliceWithFallback("TEST_STRING_SLICE_WITH_FALLBACK", []string{"fallback1", "fallback2"})
	assertNoError(t, err, "GetStringSliceWithFallback TEST_STRING_SLICE_WITH_FALLBACK")
	assertEqual(t, []string{"value1", "value2"}, value, "GetStringSliceWithFallback TEST_STRING_SLICE_WITH_FALLBACK")

	Unset("TEST_STRING_SLICE_WITH_FALLBACK")

	value, err = GetStringSliceWithFallback("TEST_STRING_SLICE_WITH_FALLBACK", []string{"fallback1", "fallback2"})
	assertNoError(t, err, "GetStringSliceWithFallback TEST_STRING_SLICE_WITH_FALLBACK")
	assertEqual(t, []string{"fallback1", "fallback2"}, value, "GetStringSliceWithFallback TEST_STRING_SLICE_WITH_FALLBACK")
}
