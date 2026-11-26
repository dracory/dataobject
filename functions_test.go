package dataobject

import (
	"testing"

	"uuid"
)

func TestToString(t *testing.T) {
	inputs := [][]any{
		{1, "1"},
		{"2", "2"},
		{true, "true"},
		{false, "false"},
		{0.123, "0.1230"}, // precission 4
		{nil, ""},
	}

	for _, v := range inputs {
		result := toString(v[0])
		expected := v[1]

		if result != expected {
			t.Error("Result", result, "DOES NOT match expected:", expected)
		}
	}
}

func TestMapStringAnyToMapStringString(t *testing.T) {
	// Test with various data types
	input := map[string]any{
		"string":  "value",
		"int":     42,
		"float":   3.14159,
		"bool":    true,
		"nil":     nil,
		"bytes":   []byte("byte array"),
		"int8":    int8(8),
		"int16":   int16(16),
		"int32":   int32(32),
		"int64":   int64(64),
		"uint":    uint(1),
		"uint8":   uint8(8),
		"uint16":  uint16(16),
		"uint32":  uint32(32),
		"uint64":  uint64(64),
		"float64": float64(64.64),
	}

	result := mapStringAnyToMapStringString(input)

	// Verify all values are converted to strings correctly
	expectedValues := map[string]string{
		"string":  "value",
		"int":     "42",
		"float":   "3.1416", // precision 4
		"bool":    "true",
		"nil":     "",
		"bytes":   "byte array",
		"int8":    "8",
		"int16":   "16",
		"int32":   "32",
		"int64":   "64",
		"uint":    "1",
		"uint8":   "8",
		"uint16":  "16",
		"uint32":  "32",
		"uint64":  "64",
		"float64": "64.6400", // precision 4
	}

	// Check that all keys exist and have the expected values
	for key, expectedValue := range expectedValues {
		if result[key] != expectedValue {
			t.Errorf("For key %s, expected %s but got %s", key, expectedValue, result[key])
		}
	}

	// Check that the result has the same number of entries as the input
	if len(result) != len(input) {
		t.Errorf("Expected result to have %d entries, but got %d", len(input), len(result))
	}
}

func TestMapStringAnyToMapStringStringWithNestedStructures(t *testing.T) {
	// Test with nested maps and slices
	input := map[string]any{
		"nested_map": map[string]any{
			"key1": "value1",
			"key2": 42,
		},
		"array": []any{1, "two", 3.0},
	}

	result := mapStringAnyToMapStringString(input)

	// Verify nested structures are converted to string representation
	if result["nested_map"] == "" {
		t.Error("Expected nested_map to be converted to a non-empty string, but got empty string")
	}

	if result["array"] == "" {
		t.Error("Expected array to be converted to a non-empty string, but got empty string")
	}
}

func TestBtos(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "hello"},
		{[]byte(""), ""},
		{[]byte("special chars: !@#$%^&*()"), "special chars: !@#$%^&*()"},
		{[]byte{72, 101, 108, 108, 111}, "Hello"}, // ASCII values for "Hello"
	}

	for _, tc := range testCases {
		result := btos(tc.input)
		if result != tc.expected {
			t.Errorf("Expected btos(%v) to be %s, but got %s", tc.input, tc.expected, result)
		}
	}
}

func TestToStringWithMoreTypes(t *testing.T) {
	// Test additional types not covered in the original test
	testCases := []struct {
		input    interface{}
		expected string
	}{
		// Test all numeric types
		{int8(127), "127"},
		{int16(32767), "32767"},
		{int32(2147483647), "2147483647"},
		{int64(9223372036854775807), "9223372036854775807"},
		{uint(42), "42"},
		{uint8(255), "255"},
		{uint16(65535), "65535"},
		{uint32(4294967295), "4294967295"},
		{uint64(18446744073709551615), "18446744073709551615"},
		{float64(3.14159265359), "3.1416"}, // precision 4

		// Test byte array
		{[]byte("byte array"), "byte array"},

		// Test complex types that use fmt.Sprint
		{struct{ Name string }{"test"}, "{test}"},
		{[]int{1, 2, 3}, "[1 2 3]"},
		{map[string]int{"one": 1}, "map[one:1]"},
	}

	for _, tc := range testCases {
		result := toString(tc.input)
		if result != tc.expected {
			t.Errorf("Expected toString(%v) to be %s, but got %s", tc.input, tc.expected, result)
		}
	}
}

func TestMapStringAnyToMapStringStringEmptyMap(t *testing.T) {
	// Test with empty map
	input := map[string]any{}

	result := mapStringAnyToMapStringString(input)

	// Verify result is an empty map, not nil
	if result == nil {
		t.Error("Expected result to be an empty map, but got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected result to be empty, but got %d entries", len(result))
	}
}

func TestIsDataObjectJSON(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Empty JSON",
			input:    "",
			expected: false,
		},
		{
			name:     "Empty JSON object",
			input:    "{}",
			expected: false,
		},
		{
			name:     "Null JSON",
			input:    "null",
			expected: false,
		},
		{
			name:     "Invalid JSON (not starting with {)",
			input:    "[1,2,3]",
			expected: false,
		},
		{
			name:     "Invalid JSON (not ending with })",
			input:    "{\"id\":\"123\"",
			expected: false,
		},
		{
			name:     "Valid JSON without id property",
			input:    "{\"name\":\"test\"}",
			expected: false,
		},
		{
			name:     "Valid JSON with id property",
			input:    "{\"id\":\"123\"}",
			expected: true,
		},
		{
			name:     "Valid JSON with id property and other properties",
			input:    "{\"id\":\"123\",\"name\":\"test\"}",
			expected: true,
		},
		{
			name:     "Valid JSON with id property in different position",
			input:    "{\"name\":\"test\",\"id\":\"123\"}",
			expected: true,
		},
		{
			name:     "Valid JSON with id property but empty value",
			input:    "{\"id\":\"\"}",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidDataObjectJSON(tc.input)
			if result != tc.expected {
				t.Errorf("Expected isDataObjectJSON(%s) to be %v, but got %v", tc.input, tc.expected, result)
			}
		})
	}
}

func Test_generateID(t *testing.T) {
	for range 10000 {
		uid1 := generateID()
		uid2 := generateID()

		if uid1 == uid2 {
			t.Error("generateID() generated the same ID twice")
		}
	}
}

func TestGenerateIDUsesUUIDv7(t *testing.T) {
	id := generateID()

	parsed, err := uuid.Parse(id)
	if err != nil {
		t.Fatalf("generateID() generated an invalid UUID: %v", err)
	}

	if parsed[6]>>4 != 7 {
		t.Errorf("generateID() generated UUID version %d, expected version 7", parsed[6]>>4)
	}
}
