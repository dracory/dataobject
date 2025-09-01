package dataobject

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/dracory/uid"
)

// generateID generates a new ID for a data object
func generateID() string {
	return uid.HumanUid()
}

// isDataObjectJSON checks if the passed JSON string is a valid DataObject JSON object
//
// Business logic:
// - checks if the JSON string is empty or "{}", or "null", returns false
// - checks if the JSON string is invalid JSON string (does not start with '{' or end with '}', returns false
// - checks if the JSON string contains the required id property, returns false
// - returns true
//
// Returns:
// - true if the JSON string is a valid DataObject JSON string
// - false otherwise
func isValidDataObjectJSON(jsonString string) bool {
	if jsonString == "" || jsonString == "{}" || jsonString == "null" {
		return false
	}

	if !strings.HasPrefix(jsonString, "{") || !strings.HasSuffix(jsonString, "}") {
		return false
	}

	if strings.Contains(jsonString, `"`+propertyId+`"`) {
		return true
	}

	return false
}

// mapStringAnyToMapStringString converts a map[string]any to map[string]string
func mapStringAnyToMapStringString(data map[string]any) map[string]string {
	result := map[string]string{}
	for k, v := range data {
		result[k] = toString(v)
	}
	return result
}

// toString converts an interface to string
func toString(v any) string {
	switch v := v.(type) {
	case string:
		return v

	case nil:
		return ""

	case []byte:
		return btos(v)

	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)

	case float64:
		return strconv.FormatFloat(v, 'f', 4, 64)

	default:
		return fmt.Sprint(v)
	}
}

// btox converts bytes to string
func btos(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
