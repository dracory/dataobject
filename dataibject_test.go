package dataobject

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strings"
	"testing"
)

func Test_DataObject_SetAndGet(t *testing.T) {
	user := NewDataObject()
	user.Set("first_name", "Jon")
	user.Set("last_name", "Doe")

	if user.Get("first_name") != "Jon" {
		t.Error("Expected: Jon, but found:", user.Get("first_name"))
	}

	if user.Get("last_name") != "Doe" {
		t.Error("Expected: Doe, but found:", user.Get("last_name"))
	}
}

func Test_DataObject_ToJson(t *testing.T) {
	user := NewDataObject()
	user.Set("first_name", "Jon")
	user.Set("last_name", "Doe")

	json, err := user.ToJSON()

	if err != nil {
		t.Error("Error must be nil, but found:", err.Error())
	}

	if !strings.Contains(json, `"first_name":"Jon"`) {
		t.Error(`Expected to contain: "first_name":"Jon", but found:`, json)
	}

	if !strings.Contains(json, `"last_name":"Doe"`) {
		t.Error(`Expected to contain: "last_name":"Doe", but found:`, json)
	}
}

func Test_DataObject_ToJsonAndBack(t *testing.T) {
	// Create a data object with some test data
	original := NewDataObject()
	original.Set("first_name", "Alice")
	original.Set("last_name", "Johnson")
	original.Set("email", "alice@example.com")

	// Convert to JSON
	jsonStr, err := original.ToJSON()
	if err != nil {
		t.Error("Failed to convert to JSON:", err.Error())
		return
	}

	// Verify JSON contains all fields
	if !strings.Contains(jsonStr, `"first_name":"Alice"`) {
		t.Error(`Expected JSON to contain "first_name":"Alice", but found:`, jsonStr)
	}
	if !strings.Contains(jsonStr, `"last_name":"Johnson"`) {
		t.Error(`Expected JSON to contain "last_name":"Johnson", but found:`, jsonStr)
	}
	if !strings.Contains(jsonStr, `"email":"alice@example.com"`) {
		t.Error(`Expected JSON to contain "email":"alice@example.com", but found:`, jsonStr)
	}

	// Create a new data object from JSON
	var data map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		t.Error("Failed to unmarshal JSON:", err.Error())
		return
	}

	newObj := NewDataObject()
	newObj.Hydrate(data)

	// Verify all data is preserved
	if newObj.Get("first_name") != "Alice" {
		t.Error("Expected first_name to be Alice, but found:", newObj.Get("first_name"))
	}
	if newObj.Get("last_name") != "Johnson" {
		t.Error("Expected last_name to be Johnson, but found:", newObj.Get("last_name"))
	}
	if newObj.Get("email") != "alice@example.com" {
		t.Error("Expected email to be alice@example.com, but found:", newObj.Get("email"))
	}

	// Verify ID is preserved
	if original.Get("id") != newObj.Get("id") {
		t.Error("Expected IDs to match after JSON round-trip, but they differ")
	}

	newObj, err = NewFromJSON(jsonStr)

	// Verify no error is returned
	if err != nil {
		t.Error("Error must be nil, but found:", err.Error())
		return
	}

	// Verify all data is preserved
	if newObj.Get("first_name") != "Alice" {
		t.Error("Expected first_name to be Alice, but found:", newObj.Get("first_name"))
	}
	if newObj.Get("last_name") != "Johnson" {
		t.Error("Expected last_name to be Johnson, but found:", newObj.Get("last_name"))
	}
	if newObj.Get("email") != "alice@example.com" {
		t.Error("Expected email to be alice@example.com, but found:", newObj.Get("email"))
	}

	// Verify ID is preserved
	if original.Get("id") != newObj.Get("id") {
		t.Error("Expected IDs to match after JSON round-trip, but they differ")
	}
}

func Test_DataObject_ToGob(t *testing.T) {
	// Create a data object with some test data
	user := NewDataObject()
	user.Set("first_name", "Jon")
	user.Set("last_name", "Doe")

	// Convert to gob
	gobData, err := user.ToGob()

	// Verify no error is returned
	if err != nil {
		t.Error("Error must be nil, but found:", err.Error())
		return
	}

	// Verify gob data is not empty
	if len(gobData) == 0 {
		t.Error("Expected gob data to not be empty")
		return
	}

	// Decode the gob data to verify it contains the correct data
	var decodedData map[string]string
	decoder := gob.NewDecoder(bytes.NewReader(gobData))
	err = decoder.Decode(&decodedData)

	if err != nil {
		t.Error("Failed to decode gob data:", err.Error())
		return
	}

	// Verify the decoded data contains the expected values
	if decodedData["first_name"] != "Jon" {
		t.Errorf("Expected decodedData[\"first_name\"] to be \"Jon\", but found %s", decodedData["first_name"])
	}

	if decodedData["last_name"] != "Doe" {
		t.Errorf("Expected decodedData[\"last_name\"] to be \"Doe\", but found %s", decodedData["last_name"])
	}

	// Verify the ID was preserved
	if decodedData["id"] == "" {
		t.Error("Expected ID to be preserved in gob data, but it was empty")
	}
}

func Test_DataObject_ToGobAndBack(t *testing.T) {
	// Create a data object with some test data
	original := NewDataObject()
	original.Set("first_name", "Jane")
	original.Set("last_name", "Smith")
	original.Set("age", "30")

	// Convert to gob
	gobData, err := original.ToGob()
	if err != nil {
		t.Error("Failed to convert to gob:", err.Error())
		return
	}

	// Create a new data object from the gob data
	restored, err := NewFromGob(gobData)
	if err != nil {
		t.Error("Failed to create from gob:", err.Error())
		return
	}

	// Verify all data was preserved correctly
	if restored.ID() != original.ID() {
		t.Errorf("Expected ID to be preserved, original: %s, restored: %s", original.ID(), restored.ID())
	}

	if restored.Get("first_name") != "Jane" {
		t.Errorf("Expected first_name to be \"Jane\", but found %s", restored.Get("first_name"))
	}

	if restored.Get("last_name") != "Smith" {
		t.Errorf("Expected last_name to be \"Smith\", but found %s", restored.Get("last_name"))
	}

	if restored.Get("age") != "30" {
		t.Errorf("Expected age to be \"30\", but found %s", restored.Get("age"))
	}

	// Verify the restored object is not marked as dirty (as it's from existing data)
	if restored.IsDirty() {
		t.Error("Expected restored object to not be dirty, but it was")
	}
}

func Test_DataObject_Init(t *testing.T) {
	// Test initialization of a new data object
	do := &DataObject{}

	// Verify data and dataChanged are nil before initialization
	if do.data != nil {
		t.Error("Expected data to be nil before Init, but it was not nil")
	}

	if do.dataChanged != nil {
		t.Error("Expected dataChanged to be nil before Init, but it was not nil")
	}

	// Call Init
	do.Init()

	// Verify data and dataChanged are initialized
	if do.data == nil {
		t.Error("Expected data to be initialized after Init, but it was nil")
	}

	if do.dataChanged == nil {
		t.Error("Expected dataChanged to be initialized after Init, but it was nil")
	}

	// Call Init again to ensure it doesn't overwrite existing maps
	do.data["test"] = "value"
	do.dataChanged["test"] = "value"

	do.Init()

	// Verify data and dataChanged were not reset
	if do.data["test"] != "value" {
		t.Error("Expected data to retain values after second Init, but it did not")
	}

	if do.dataChanged["test"] != "value" {
		t.Error("Expected dataChanged to retain values after second Init, but it did not")
	}
}

func Test_DataObject_Data(t *testing.T) {
	do := NewDataObject()
	do.Set("key1", "value1")
	do.Set("key2", "value2")

	data := do.Data()

	if len(data) != 3 { // 2 keys + id
		t.Errorf("Expected data to have 3 entries, but found %d", len(data))
	}

	if data["key1"] != "value1" {
		t.Errorf("Expected data[\"key1\"] to be \"value1\", but found %s", data["key1"])
	}

	if data["key2"] != "value2" {
		t.Errorf("Expected data[\"key2\"] to be \"value2\", but found %s", data["key2"])
	}
}

func Test_DataObject_DataChanged(t *testing.T) {
	do := NewDataObject()
	do.Set("key1", "value1")
	do.Set("key2", "value2")

	dataChanged := do.DataChanged()

	if len(dataChanged) != 3 { // 2 keys + id
		t.Errorf("Expected dataChanged to have 3 entries, but found %d", len(dataChanged))
	}

	if dataChanged["key1"] != "value1" {
		t.Errorf("Expected dataChanged[\"key1\"] to be \"value1\", but found %s", dataChanged["key1"])
	}

	if dataChanged["key2"] != "value2" {
		t.Errorf("Expected dataChanged[\"key2\"] to be \"value2\", but found %s", dataChanged["key2"])
	}

	// Test after MarkAsNotDirty
	do.MarkAsNotDirty()
	dataChanged = do.DataChanged()

	if len(dataChanged) != 0 {
		t.Errorf("Expected dataChanged to be empty after MarkAsNotDirty, but found %d entries", len(dataChanged))
	}
}

func Test_DataObject_IsDirty(t *testing.T) {
	do := NewDataObject()

	// New object with ID should be dirty
	if !do.IsDirty() {
		t.Error("Expected new object to be dirty, but it was not")
	}

	// After marking as not dirty, should not be dirty
	do.MarkAsNotDirty()
	if do.IsDirty() {
		t.Error("Expected object to not be dirty after MarkAsNotDirty, but it was")
	}

	// After setting a value, should be dirty again
	do.Set("key", "value")
	if !do.IsDirty() {
		t.Error("Expected object to be dirty after Set, but it was not")
	}
}

func Test_DataObject_SetData(t *testing.T) {
	do := NewDataObject()
	do.MarkAsNotDirty() // Clear initial ID setting

	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	do.SetData(data)

	// Check data was set correctly
	if do.Get("key1") != "value1" {
		t.Errorf("Expected do.Get(\"key1\") to be \"value1\", but found %s", do.Get("key1"))
	}

	if do.Get("key2") != "value2" {
		t.Errorf("Expected do.Get(\"key2\") to be \"value2\", but found %s", do.Get("key2"))
	}

	// Check object is marked as dirty
	if !do.IsDirty() {
		t.Error("Expected object to be dirty after SetData, but it was not")
	}

	// Check dataChanged contains the new values
	dataChanged := do.DataChanged()
	if dataChanged["key1"] != "value1" || dataChanged["key2"] != "value2" {
		t.Error("Expected dataChanged to contain the new values, but it did not")
	}
}

func Test_DataObject_Hydrate(t *testing.T) {
	do := NewDataObject()
	do.MarkAsNotDirty() // Clear initial ID setting

	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	do.Hydrate(data)

	// Check data was set correctly
	if do.Get("key1") != "value1" {
		t.Errorf("Expected do.Get(\"key1\") to be \"value1\", but found %s", do.Get("key1"))
	}

	if do.Get("key2") != "value2" {
		t.Errorf("Expected do.Get(\"key2\") to be \"value2\", but found %s", do.Get("key2"))
	}

	// Check object is NOT marked as dirty after Hydrate
	if do.IsDirty() {
		t.Error("Expected object to NOT be dirty after Hydrate, but it was")
	}
}

func Test_DataObject_IDMethods(t *testing.T) {
	do := NewDataObject()
	originalID := do.ID()

	// Verify ID is not empty
	if originalID == "" {
		t.Error("Expected ID to not be empty, but it was")
	}

	// Set a new ID
	newID := "custom-id-123"
	do.SetID(newID)

	// Verify ID was updated
	if do.ID() != newID {
		t.Errorf("Expected ID to be %s after SetID, but found %s", newID, do.ID())
	}
}
