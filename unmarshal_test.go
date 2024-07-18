package env

import (
	"reflect"
	"testing"
)

type Config struct {
	StringField string   `env:"STRING_FIELD,default=default_value"`
	IntField    int      `env:"INT_FIELD,default=123"`
	BoolField   bool     `env:"BOOL_FIELD,default=true"`
	FloatField  float64  `env:"FLOAT_FIELD,default=1.23"`
	SliceField  []string `env:"SLICE_FIELD,default=item1"`
}

type NestedConfig struct {
	NestedField string `env:"NESTED_FIELD,default=nested"`
}

type ParentConfig struct {
	Config
	Nested NestedConfig `env:"NESTED"`
}

func TestUnmarshal(t *testing.T) {
	setEnvForTest(t, "STRING_FIELD", "string_value")
	setEnvForTest(t, "INT_FIELD", "456")
	setEnvForTest(t, "BOOL_FIELD", "false")
	setEnvForTest(t, "FLOAT_FIELD", "4.56")
	setEnvForTest(t, "SLICE_FIELD", "item1,item2")

	var cfg Config
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal")

	expected := Config{
		StringField: "string_value",
		IntField:    456,
		BoolField:   false,
		FloatField:  4.56,
		SliceField:  []string{"item1", "item2"},
	}

	assertEqual(t, expected, cfg, "Unmarshal")
}

func TestUnmarshalNested(t *testing.T) {
	setEnvForTest(t, "STRING_FIELD", "string_value")
	setEnvForTest(t, "INT_FIELD", "456")
	setEnvForTest(t, "BOOL_FIELD", "false")
	setEnvForTest(t, "FLOAT_FIELD", "4.56")
	setEnvForTest(t, "SLICE_FIELD", "item1,item2")
	setEnvForTest(t, "NESTED_NESTED_FIELD", "nested_value")

	var cfg ParentConfig
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal")

	expected := ParentConfig{
		Config: Config{
			StringField: "string_value",
			IntField:    456,
			BoolField:   false,
			FloatField:  4.56,
			SliceField:  []string{"item1", "item2"},
		},
		Nested: NestedConfig{
			NestedField: "nested_value",
		},
	}

	assertEqual(t, expected, cfg, "UnmarshalNested")
}

func TestUnmarshalUntaggedFields(t *testing.T) {
	type UntaggedFieldsConfig struct {
		UntaggedField1 string
		UntaggedField2 int
		TaggedField    string `env:"TAGGED_FIELD,default=default_value"`
	}

	setEnvForTest(t, "TAGGED_FIELD", "value")

	var cfg UntaggedFieldsConfig
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal")

	expected := UntaggedFieldsConfig{
		UntaggedField1: "",
		UntaggedField2: 0,
		TaggedField:    "value",
	}

	assertEqual(t, expected, cfg, "UnmarshalUntaggedFields")
}

func TestUnmarshalFloat(t *testing.T) {
	setEnvForTest(t, "FLOAT32_VALUE", "3.14")
	setEnvForTest(t, "FLOAT64_VALUE", "3.14159")

	var cfg struct {
		Float32Value float32 `env:"FLOAT32_VALUE"`
		Float64Value float64 `env:"FLOAT64_VALUE"`
	}
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal FloatConfig")

	expected := struct {
		Float32Value float32
		Float64Value float64
	}{
		Float32Value: 3.14,
		Float64Value: 3.14159,
	}

	assertEqual(t, expected.Float32Value, cfg.Float32Value, "Float32Value")
	assertEqual(t, expected.Float64Value, cfg.Float64Value, "Float64Value")
}

func TestUnmarshalUnsupportedKind(t *testing.T) {
	setEnvForTest(t, "UNSUPPORTED", "invalid")

	var cfg struct {
		Unsupported complex64 `env:"UNSUPPORTED"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal Unsupported kind")
}

func TestUnmarshalSetFieldIntError(t *testing.T) {
	setEnvForTest(t, "INVALID_INT", "invalid")

	var cfg struct {
		InvalidInt int `env:"INVALID_INT"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal InvalidInt")
}

func TestUnmarshalSetFieldFloatError(t *testing.T) {
	setEnvForTest(t, "INVALID_FLOAT", "invalid")

	var cfg struct {
		InvalidFloat float64 `env:"INVALID_FLOAT"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal InvalidFloat")
}

func TestUnmarshalSetFieldNestedError(t *testing.T) {
	type NestedConfig struct {
		NestedField string `env:"NESTED_FIELD,required"`
	}

	setEnvForTest(t, "NESTED_FIELD", "") // Setting an empty value to trigger required error

	var cfg struct {
		Nested NestedConfig
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal NestedConfig")
}

func TestUnmarshalRequired(t *testing.T) {
	type RequiredConfig struct {
		RequiredVar string `env:"REQUIRED_VAR,required"`
	}

	var cfg RequiredConfig
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal Required")

	expectedErr := "required environment variable REQUIRED_VAR is not set"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("expected error %s, got %s", expectedErr, err.Error())
	}

	setEnvForTest(t, "REQUIRED_VAR", "value")
	err = Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal Required")
	assertEqual(t, "value", cfg.RequiredVar, "Unmarshal Required")
}

func TestUnmarshalSetFieldErrors(t *testing.T) {
	tests := []struct {
		envKey    string
		envValue  string
		fieldType string
	}{
		{"INVALID_UINT", "invalid", "uint"},
		{"INVALID_FLOAT", "invalid", "float32"},
		{"UNSUPPORTED", "invalid", "complex64"},
	}

	for _, tt := range tests {
		t.Run(tt.fieldType, func(t *testing.T) {
			setEnvForTest(t, tt.envKey, tt.envValue)

			var cfg struct {
				InvalidUint  uint      `env:"INVALID_UINT"`
				InvalidFloat float32   `env:"INVALID_FLOAT"`
				Unsupported  complex64 `env:"UNSUPPORTED"`
			}

			err := Unmarshal(&cfg)
			assertError(t, err, "Unmarshal "+tt.fieldType)
		})
	}
}

func TestSetFieldUint(t *testing.T) {
	type Config struct {
		UIntField uint `env:"UINT_FIELD"`
	}

	var cfg Config

	err := Set("UINT_FIELD", "42")
	assertNoError(t, err, "Set UINT_FIELD")

	field := reflect.ValueOf(&cfg).Elem().FieldByName("UIntField")
	err = setField(field, "42")
	assertNoError(t, err, "setField Uint")

	assertEqual(t, uint(42), cfg.UIntField, "UintField value")
}

func TestSetFieldBoolError(t *testing.T) {
	type Config struct {
		BoolField bool `env:"BOOL_FIELD"`
	}

	var cfg Config

	err := Set("BOOL_FIELD", "invalid")
	assertNoError(t, err, "Set BOOL_FIELD")

	field := reflect.ValueOf(&cfg).Elem().FieldByName("BoolField")
	err = setField(field, "invalid")
	assertError(t, err, "setField Bool")
}

func TestUnsupportedSliceElementKind(t *testing.T) {
	type Config struct {
		Unsupported []complex64 `env:"UNSUPPORTED"`
	}

	var cfg Config

	err := Set("UNSUPPORTED", "invalid")
	assertNoError(t, err, "Set UNSUPPORTED")

	field := reflect.ValueOf(&cfg).Elem().FieldByName("Unsupported")
	err = setField(field, "invalid")
	assertError(t, err, "setField Unsupported")
}

func TestUnmarshalFallback(t *testing.T) {
	type FallbackConfig struct {
		FallbackField string `env:"FALLBACK_FIELD,fallback=fallback_value"`
	}

	var cfg FallbackConfig
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal Fallback")

	expected := FallbackConfig{
		FallbackField: "fallback_value",
	}

	assertEqual(t, expected, cfg, "UnmarshalFallback")
}

// New tests for additional slice types

func TestUnmarshalSliceBool(t *testing.T) {
	setEnvForTest(t, "SLICE_BOOL", "true,false,true")

	var cfg struct {
		SliceBool []bool `env:"SLICE_BOOL"`
	}
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal SliceBool")

	expected := struct {
		SliceBool []bool
	}{
		SliceBool: []bool{true, false, true},
	}

	assertEqual(t, expected.SliceBool, cfg.SliceBool, "SliceBool")
}

func TestUnmarshalSliceInt(t *testing.T) {
	setEnvForTest(t, "SLICE_INT", "1,2,3")

	var cfg struct {
		SliceInt []int `env:"SLICE_INT"`
	}
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal SliceInt")

	expected := struct {
		SliceInt []int
	}{
		SliceInt: []int{1, 2, 3},
	}

	assertEqual(t, expected.SliceInt, cfg.SliceInt, "SliceInt")
}

func TestUnmarshalSliceUint(t *testing.T) {
	setEnvForTest(t, "SLICE_UINT", "1,2,3")

	var cfg struct {
		SliceUint []uint `env:"SLICE_UINT"`
	}
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal SliceUint")

	expected := struct {
		SliceUint []uint
	}{
		SliceUint: []uint{1, 2, 3},
	}

	assertEqual(t, expected.SliceUint, cfg.SliceUint, "SliceUint")
}

func TestUnmarshalSliceFloat(t *testing.T) {
	setEnvForTest(t, "SLICE_FLOAT", "1.1,2.2,3.3")

	var cfg struct {
		SliceFloat []float64 `env:"SLICE_FLOAT"`
	}
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal SliceFloat")

	expected := struct {
		SliceFloat []float64
	}{
		SliceFloat: []float64{1.1, 2.2, 3.3},
	}

	assertEqual(t, expected.SliceFloat, cfg.SliceFloat, "SliceFloat")
}

func TestUnmarshalSliceBoolError(t *testing.T) {
	setEnvForTest(t, "SLICE_BOOL_ERROR", "true,notabool,true")

	var cfg struct {
		SliceBool []bool `env:"SLICE_BOOL_ERROR"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal SliceBoolError")
}

func TestUnmarshalSliceIntError(t *testing.T) {
	setEnvForTest(t, "SLICE_INT_ERROR", "1,two,3")

	var cfg struct {
		SliceInt []int `env:"SLICE_INT_ERROR"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal SliceIntError")
}

func TestUnmarshalSliceUintError(t *testing.T) {
	setEnvForTest(t, "SLICE_UINT_ERROR", "1,-2,3")

	var cfg struct {
		SliceUint []uint `env:"SLICE_UINT_ERROR"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal SliceUintError")
}

func TestUnmarshalSliceFloatError(t *testing.T) {
	setEnvForTest(t, "SLICE_FLOAT_ERROR", "1.1,two.two,3.3")

	var cfg struct {
		SliceFloat []float64 `env:"SLICE_FLOAT_ERROR"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal SliceFloatError")
}
