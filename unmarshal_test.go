package env

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

type Config struct {
	StringField              string   `env:"STRING_FIELD,default=default_value"`
	IntField                 int      `env:"INT_FIELD,default=123"`
	BoolField                bool     `env:"BOOL_FIELD,default=true"`
	FloatField               float64  `env:"FLOAT_FIELD,default=1.23"`
	StringSliceField         []string `env:"SLICE_FIELD,default=item1"`
	StringSliceFieldMultiple []string `env:"SLICE_FIELD_MULTI,default=[item1,item2,item3]"`
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
	setEnvForTest(t, "SLICE_FIELD_MULTI", "item1,item2,item3,item4")

	var cfg Config
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal")

	expected := Config{
		StringField:              "string_value",
		IntField:                 456,
		BoolField:                false,
		FloatField:               4.56,
		StringSliceField:         []string{"item1", "item2"},
		StringSliceFieldMultiple: []string{"item1", "item2", "item3", "item4"},
	}

	assertEqual(t, expected, cfg, "Unmarshal")
}

func TestUnmarshalWithDefaults(t *testing.T) {
	var cfg Config
	err := Unmarshal(&cfg)
	assertNoError(t, err)
	assertEqual(t, "default_value", cfg.StringField)
	assertEqual(t, 123, cfg.IntField)
	assertEqual(t, true, cfg.BoolField)
	assertEqual(t, 1.23, cfg.FloatField)
	assertEqual(t, []string{"item1"}, cfg.StringSliceField)
	assertEqual(t, []string{"item1", "item2", "item3"}, cfg.StringSliceFieldMultiple)
}

func TestUnmarshalDefaultsFromCode(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,default=localhost"`
		Port int    `env:"PORT"`
	}

	setEnvForTest(t, "HOST", "envhost")
	setEnvForTest(t, "PORT", "8080")

	cfg := &Config{
		Host: "syntaqx.com",
		Port: 3306,
	}

	err := Unmarshal(cfg)
	assertNoError(t, err, "Unmarshal")

	expected := Config{
		Host: "envhost",
		Port: 8080,
	}

	assertEqual(t, expected, *cfg, "UnmarshalDefaultsFromCode")
}

func TestUnmarshalDefaultsFromTags(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,default=localhost"`
		Port int    `env:"PORT"`
	}

	cfg := &Config{}

	err := Unmarshal(cfg)
	assertNoError(t, err, "Unmarshal")

	expected := Config{
		Host: "localhost",
		Port: 0,
	}

	assertEqual(t, expected, *cfg, "UnmarshalDefaultsFromTags")
}

func TestUnmarshalDefaultsFromCodeAndTags(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,default=localhost"`
		Port int    `env:"PORT"`
	}

	cfg := &Config{
		Host: "syntaqx.com",
		Port: 3306,
	}

	err := Unmarshal(cfg)
	assertNoError(t, err, "Unmarshal")

	expected := Config{
		Host: "localhost",
		Port: 3306,
	}

	assertEqual(t, expected, *cfg, "UnmarshalDefaultsFromCodeAndTags")
}

func TestUnmarshalNested(t *testing.T) {
	setEnvForTest(t, "STRING_FIELD", "string_value")
	setEnvForTest(t, "INT_FIELD", "456")
	setEnvForTest(t, "BOOL_FIELD", "false")
	setEnvForTest(t, "FLOAT_FIELD", "4.56")
	setEnvForTest(t, "SLICE_FIELD", "item1,item2")
	setEnvForTest(t, "SLICE_FIELD_MULTI", "item1,item2,item3,item4")
	setEnvForTest(t, "NESTED_NESTED_FIELD", "nested_value")

	var cfg ParentConfig
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal")

	expected := ParentConfig{
		Config: Config{
			StringField:              "string_value",
			IntField:                 456,
			BoolField:                false,
			FloatField:               4.56,
			StringSliceField:         []string{"item1", "item2"},
			StringSliceFieldMultiple: []string{"item1", "item2", "item3", "item4"},
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

func TestSetFieldEmptyValue(t *testing.T) {
	type Config struct {
		StringField string `env:"STRING_FIELD"`
		IntField    int    `env:"INT_FIELD"`
		BoolField   bool   `env:"BOOL_FIELD"`
	}

	var cfg Config

	field := reflect.ValueOf(&cfg).Elem().FieldByName("StringField")
	err := setField(field, "")
	assertNoError(t, err, "setField empty string")
	assertEqual(t, "", cfg.StringField, "StringField")

	field = reflect.ValueOf(&cfg).Elem().FieldByName("IntField")
	err = setField(field, "")
	assertNoError(t, err, "setField empty int")
	assertEqual(t, 0, cfg.IntField, "IntField")

	field = reflect.ValueOf(&cfg).Elem().FieldByName("BoolField")
	err = setField(field, "")
	assertNoError(t, err, "setField empty bool")
	assertEqual(t, false, cfg.BoolField, "BoolField")
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

func TestParseTag(t *testing.T) {
	type TestCase struct {
		Tag          string
		ExpectedOpts tagOptions
	}

	testCases := []TestCase{
		{
			Tag: "NOT_REQUIRED,default=required",
			ExpectedOpts: tagOptions{
				keys:     []string{"NOT_REQUIRED"},
				fallback: "required",
				required: false,
			},
		},
		{
			Tag: "REQUIRED,required",
			ExpectedOpts: tagOptions{
				keys:     []string{"REQUIRED"},
				fallback: "",
				required: true,
			},
		},
		{
			Tag: "REQUIRED_WITH_DEFAULT,default=default,required",
			ExpectedOpts: tagOptions{
				keys:     []string{"REQUIRED_WITH_DEFAULT"},
				fallback: "default",
				required: true,
			},
		},
		{
			Tag: "SINGLE_KEY,required,default=default",
			ExpectedOpts: tagOptions{
				keys:     []string{"SINGLE_KEY"},
				fallback: "default",
				required: true,
			},
		},
		{
			Tag: "MULTI_KEY1|MULTI_KEY2|MULTI_KEY3,required,default=default",
			ExpectedOpts: tagOptions{
				keys:     []string{"MULTI_KEY1", "MULTI_KEY2", "MULTI_KEY3"},
				fallback: "default",
				required: true,
			},
		},
		{
			Tag: "SQUARE_BRACKETS,default=[item1,item2,item3]",
			ExpectedOpts: tagOptions{
				keys:     []string{"SQUARE_BRACKETS"},
				fallback: "item1,item2,item3",
				required: false,
			},
		},
		{
			Tag: "SQUARE_BRACKETS,default=[item1,item2,item3],required",
			ExpectedOpts: tagOptions{
				keys:     []string{"SQUARE_BRACKETS"},
				fallback: "item1,item2,item3",
				required: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Tag, func(t *testing.T) {
			opts := parseTag(tc.Tag)
			if !reflect.DeepEqual(opts, tc.ExpectedOpts) {
				t.Errorf("parseTag(%s) returned %+v, expected %+v", tc.Tag, opts, tc.ExpectedOpts)
			}
		})
	}
}

func TestIsZeroValue(t *testing.T) {
	tests := []struct {
		name  string
		field reflect.Value
		want  bool
	}{
		{"ZeroString", reflect.ValueOf(""), true},
		{"NonZeroString", reflect.ValueOf("non-zero"), false},
		{"ZeroInt", reflect.ValueOf(0), true},
		{"NonZeroInt", reflect.ValueOf(123), false},
		{"ZeroBool", reflect.ValueOf(false), true},
		{"NonZeroBool", reflect.ValueOf(true), false},
		{"ZeroFloat", reflect.ValueOf(0.0), true},
		{"NonZeroFloat", reflect.ValueOf(1.23), false},
		{"ZeroSlice", reflect.ValueOf([]string(nil)), true},
		{"NonZeroSlice", reflect.ValueOf([]string{"item1"}), false},
		{"InvalidValue", reflect.Value{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isZeroValue(tt.field); got != tt.want {
				t.Errorf("isZeroValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalFileOption(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,default=localhost"`
		Port int    `env:"PORT"`
		Key  string `env:"KEY,file"`
	}

	// Create a temporary file with the content
	fileContent := "file_content"
	tmpFile, err := os.CreateTemp("", "example")
	assertNoError(t, err, "CreateTemp")
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(fileContent)
	assertNoError(t, err, "WriteString")

	err = tmpFile.Close()
	assertNoError(t, err, "Close")

	setEnvForTest(t, "HOST", "envhost")
	setEnvForTest(t, "PORT", "8080")
	setEnvForTest(t, "KEY", tmpFile.Name())

	cfg := &Config{
		Host: "syntaqx.com",
		Port: 3306,
	}

	err = Unmarshal(cfg)
	assertNoError(t, err, "Unmarshal")

	expected := Config{
		Host: "envhost",
		Port: 8080,
		Key:  fileContent,
	}

	assertEqual(t, expected, *cfg, "UnmarshalFileOption")
}

func TestUnmarshalFieldFileError(t *testing.T) {
	type Config struct {
		Key string `env:"KEY,file"`
	}

	// Set an invalid file path
	setEnvForTest(t, "KEY", "/invalid/path/to/file")

	var cfg Config
	err := Unmarshal(&cfg)
	assertError(t, err, "UnmarshalFieldFileError")

	expectedErrPrefix := "open /invalid/path/to/file"
	if err != nil && !strings.HasPrefix(err.Error(), expectedErrPrefix) {
		t.Errorf("expected error to start with %s, got %s", expectedErrPrefix, err.Error())
	}
}

func TestReadFileContentError(t *testing.T) {
	_, err := readFileContent("/invalid/path/to/file")
	assertError(t, err, "readFileContentError")

	expectedErrPrefix := "open /invalid/path/to/file"
	if err != nil && !strings.HasPrefix(err.Error(), expectedErrPrefix) {
		t.Errorf("expected error to start with %s, got %s", expectedErrPrefix, err.Error())
	}
}

func TestUnmarshalExpand(t *testing.T) {
	setEnvForTest(t, "HOST", "localhost")
	setEnvForTest(t, "PORT", "8080")
	setEnvForTest(t, "BASE_URL", "http://${HOST}:${PORT}")

	type Config struct {
		BaseURL string `env:"BASE_URL,expand"`
	}

	var cfg Config
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal with expand")

	expected := Config{
		BaseURL: "http://localhost:8080",
	}

	assertEqual(t, expected, cfg, "UnmarshalExpand")
}

func TestUnmarshalExpandWithDefault(t *testing.T) {
	setEnvForTest(t, "HOST", "localhost")
	setEnvForTest(t, "PORT", "8080")

	type Config struct {
		BaseURL string `env:"BASE_URL,default=http://${HOST}:${PORT}/api,expand"`
	}

	var cfg Config
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal with expand and default")

	expected := Config{
		BaseURL: "http://localhost:8080/api",
	}

	assertEqual(t, expected, cfg, "UnmarshalExpandWithDefault")
}

func TestUnmarshalExpandWithMissingEnv(t *testing.T) {
	type Config struct {
		BaseURL string `env:"BASE_URL,default=http://${HOST}:${PORT}/api,expand"`
	}

	var cfg Config
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal with expand and missing env variables")

	expected := Config{
		BaseURL: "http://:/api", // Expanded with empty strings for missing HOST and PORT
	}

	assertEqual(t, expected, cfg, "UnmarshalExpandWithMissingEnv")
}

func TestGetDefaultFromStructWithFallback(t *testing.T) {
	type Config struct {
		Host    string `env:"HOST,default=localhost"`
		Port    string `env:"PORT,default=8080"`
		Address string `env:"ADDRESS,default=${HOST}:${PORT},expand"`
	}

	var cfg Config
	defaultHost := getDefaultFromStruct("HOST", &cfg)
	defaultPort := getDefaultFromStruct("PORT", &cfg)

	if defaultHost != "localhost" {
		t.Errorf("expected default host to be 'localhost', got '%s'", defaultHost)
	}

	if defaultPort != "8080" {
		t.Errorf("expected default port to be '8080', got '%s'", defaultPort)
	}
}

func TestGetDefaultFromStructWithNestedStruct(t *testing.T) {
	type NestedConfig struct {
		NestedField string `env:"NESTED_FIELD,default=nested_default"`
	}

	type Config struct {
		Host   string       `env:"HOST,default=localhost"`
		Nested NestedConfig `env:"NESTED"`
	}

	var cfg Config
	defaultNestedField := getDefaultFromStruct("NESTED_FIELD", &cfg)

	if defaultNestedField != "nested_default" {
		t.Errorf("expected default nested field to be 'nested_default', got '%s'", defaultNestedField)
	}
}

func TestExpandVariables(t *testing.T) {
	setEnvForTest(t, "HOST", "localhost")
	setEnvForTest(t, "PORT", "8080")

	type Config struct {
		BaseURL1 string `env:"BASE_URL1,default=http://${HOST}:${PORT}/api,expand"`
		BaseURL2 string `env:"BASE_URL2,default=http://$HOST:$PORT/api,expand"`
	}

	var cfg Config
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal with expand and default")

	expected := Config{
		BaseURL1: "http://localhost:8080/api",
		BaseURL2: "http://localhost:8080/api",
	}

	assertEqual(t, expected, cfg, "ExpandVariables")
}
