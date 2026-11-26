package dataobject

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

// New creates a new data object with a unique ID
//
// Business logic:
// - instantiates a new data object
// - generates an ID, using uuid.NewV7()
//
// Note! The object is marked as dirty, as ID is set
//
// Returns:
// - a new data object
func New() *DataObject {
	o := &DataObject{}
	o.SetID(generateID())
	return o
}

// NewFromJSON creates a new data object
// and hydrates it with the passed JSON string.
//
// # The JSON string is expected to be a valid DataObject JSON object
//
// Note: the object is marked as not dirty, as it is existing data
//
// Business logic:
// - instantiates a new data object
// - hydrates it with the passed JSON string
//
// Returns:
// - a new data object
// - an error if any
func NewFromJSON(jsonString string) (*DataObject, error) {
	if !isValidDataObjectJSON(jsonString) {
		return nil, errors.New("invalid json: must be a valid dataobject json object")
	}

	var e any

	jsonError := json.Unmarshal([]byte(jsonString), &e)

	if jsonError != nil {
		return nil, jsonError
	}

	data := mapStringAnyToMapStringString(e.(map[string]any))

	if data == nil {
		return nil, errors.New("invalid data from json")
	}

	if data[propertyId] == "" {
		return nil, errors.New("invalid json: missing id")
	}

	do := NewFromData(data)

	return do, nil
}

// NewFromData creates a new data object
// and hydrates it with the passed data
//
// Note: the object is marked as not dirty, as it is existing data,
// unless ID is missing, in which case it is added (and the object is marked as dirty)
//
// Business logic:
// - instantiates a new data object
// - hydrates it with the passed data
// - generates an ID, using uuid.NewV7() if missing (but is marked as dirty)
//
// Returns:
// - a new data object
func NewFromData(data map[string]string) *DataObject {
	if data == nil {
		return nil
	}

	o := &DataObject{}
	o.Hydrate(data)

	if o.ID() == "" {
		o.SetID(generateID())
	}

	return o
}

// NewFromGob creates a new data object
// and hydrates it with the passed gob-encoded byte array.
//
// # The gob data is expected to be an encoded map[string]string
//
// Note: the object is marked as not dirty, as it is existing data
//
// Business logic:
// - instantiates a new data object
// - decodes the gob data
// - hydrates it with the decoded data
//
// Returns:
// - a new data object
// - an error if any
func NewFromGob(gobData []byte) (*DataObject, error) {
	var data map[string]string

	decoder := gob.NewDecoder(bytes.NewReader(gobData))
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	if data[propertyId] == "" {
		return nil, errors.New("invalid gob data: missing id")
	}

	do := NewFromData(data)

	return do, nil
}

// Deprecated: NewDataObject is deprecated, use New() instead.
// Creates a new data object with a unique ID
func NewDataObject() *DataObject {
	return New()
}

// Deprecated: NewDataObjectFromExistingData is deprecated, use NewFromData() instead.
// Creates a new data object and hydrates it with the passed data
func NewDataObjectFromExistingData(data map[string]string) *DataObject {
	return NewFromData(data)
}

// Deprecated: NewDataObjectFromJSON is deprecated, use NewFromJSON() instead.
// Creates a new data object and hydrates it with the passed JSON string
func NewDataObjectFromJSON(jsonString string) (do *DataObject, err error) {
	return NewFromJSON(jsonString)
}
