package dataobject

// DataObjectInterface is an interface for a data object
type DataObjectInterface interface {

	// ID returns the ID of the object
	ID() string

	// SetID sets the ID of the object
	SetID(id string)

	// GetData returns the data for the object
	Data() map[string]string

	// SetData sets the data for the object
	SetData(data map[string]string)

	// GetChangedData returns the data that has been changed since the last hydration
	DataChanged() map[string]string

	// MarkAsNotDirty marks the object as not dirty
	// If no columns specified, marks all fields as not dirty
	// If columns specified, marks only those columns as not dirty
	MarkAsNotDirty(columns ...string)

	// MarkAsDirty marks the object as dirty
	// If no columns specified, marks all fields as dirty
	// If columns specified, marks only those columns as dirty
	MarkAsDirty(columns ...string)

	// IsDirty returns if data has been modified
	IsDirty() bool

	// Hydrates the data object with data
	Hydrate(map[string]string)

	// Get returns the value for a key
	Get(key string) string

	// Set sets a key-value pair in the data object
	Set(key string, value string)

	// ToJSON converts the DataObject to a JSON string
	ToJSON() (string, error)

	// ToGob converts the DataObject to a gob-encoded byte array
	ToGob() ([]byte, error)
}
