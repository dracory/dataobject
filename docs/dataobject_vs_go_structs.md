# DataObject vs Standard Go Structs: A Comparison

## Introduction

This document compares the DataObject implementation with standard Go structs. While both are used to organize and manage data in Go applications, they have different approaches and features. This comparison will help developers understand when to use each approach in their Go applications.

## Core Concepts

### Standard Go Structs

Go structs are:
- Composite data types that group together variables under a single name
- Statically typed with explicitly defined fields
- Can have methods attached to them
- Support embedding for composition
- Do not have built-in change tracking
- Typically serialized/deserialized using the `encoding/json` package
- Can use struct tags for metadata

### DataObject

Our DataObject implementation follows these principles:
- Has a default constructor
- Has a unique identifier (ID)
- All fields are private (stored in an internal map)
- Fields are accessed via public getter and setter methods
- Tracks changes to data internally
- Can return all data as a map
- Designed for efficient serialization to data stores

## Detailed Comparison

| Feature | Standard Go Structs | DataObject |
|---------|-------------------|---------------|
| **Data Storage** | Explicitly defined fields with specific types | Internal map[string]string |
| **Type Safety** | Strong static typing at compile time | Strong typing via getters/setters at runtime |
| **Constructor** | No built-in constructors (use factory functions) | Multiple constructors (`New()`, `NewFromData()`, `NewFromJSON()`) |
| **Property Access** | Direct field access or via getter/setter methods | Via explicit getter/setter methods only |
| **Change Tracking** | Not built-in (requires custom implementation) | Built-in via `IsDirty()` and `DataChanged()` |
| **Unique Identifier** | Not required (add manually if needed) | Required (ID field) |
| **Serialization** | Via `encoding/json` and struct tags | Custom JSON serialization |
| **Memory Efficiency** | More efficient for known fields | Less efficient due to map storage |
| **Dynamic Fields** | Not supported (fields must be defined at compile time) | Supported (can add any key/value pair) |
| **Composition** | Via struct embedding | Via struct embedding |
| **Method Chaining** | Possible but not idiomatic | Supported by design |

## Implementation Differences

### Standard Go Struct Example

```go
package models

import (
    "encoding/json"
    "time"
    
    "github.com/gouniverse/uid"
)

// User represents a user in the system
type User struct {
    ID        string    `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user
func NewUser() *User {
    now := time.Now().UTC()
    return &User{
        ID:        uid.HumanUid(),
        Status:    "active",
        CreatedAt: now,
        UpdatedAt: now,
    }
}

// NewUserFromMap creates a user from a map
func NewUserFromMap(data map[string]string) *User {
    user := &User{
        ID:        data["id"],
        FirstName: data["first_name"],
        LastName:  data["last_name"],
        Email:     data["email"],
        Status:    data["status"],
    }
    
    if data["created_at"] != "" {
        createdAt, _ := time.Parse(time.RFC3339, data["created_at"])
        user.CreatedAt = createdAt
    }
    
    if data["updated_at"] != "" {
        updatedAt, _ := time.Parse(time.RFC3339, data["updated_at"])
        user.UpdatedAt = updatedAt
    }
    
    return user
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
    return u.Status == "active"
}

// FullName returns the user's full name
func (u *User) FullName() string {
    return u.FirstName + " " + u.LastName
}

// SetFirstName sets the first name and updates the UpdatedAt timestamp
func (u *User) SetFirstName(firstName string) *User {
    u.FirstName = firstName
    u.UpdatedAt = time.Now().UTC()
    return u
}

// SetLastName sets the last name and updates the UpdatedAt timestamp
func (u *User) SetLastName(lastName string) *User {
    u.LastName = lastName
    u.UpdatedAt = time.Now().UTC()
    return u
}

// SetEmail sets the email and updates the UpdatedAt timestamp
func (u *User) SetEmail(email string) *User {
    u.Email = email
    u.UpdatedAt = time.Now().UTC()
    return u
}

// SetStatus sets the status and updates the UpdatedAt timestamp
func (u *User) SetStatus(status string) *User {
    u.Status = status
    u.UpdatedAt = time.Now().UTC()
    return u
}

// ToJSON converts the user to JSON
func (u *User) ToJSON() (string, error) {
    bytes, err := json.Marshal(u)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

// FromJSON populates the user from JSON
func (u *User) FromJSON(jsonStr string) error {
    return json.Unmarshal([]byte(jsonStr), u)
}

// ToMap converts the user to a map
func (u *User) ToMap() map[string]string {
    return map[string]string{
        "id":         u.ID,
        "first_name": u.FirstName,
        "last_name":  u.LastName,
        "email":      u.Email,
        "status":     u.Status,
        "created_at": u.CreatedAt.Format(time.RFC3339),
        "updated_at": u.UpdatedAt.Format(time.RFC3339),
    }
}
```

### Standard Go Struct with Change Tracking

```go
package models

import (
    "encoding/json"
    "time"
    
    "github.com/gouniverse/uid"
)

// User represents a user in the system
type User struct {
    ID             string            `json:"id"`
    FirstName      string            `json:"first_name"`
    LastName       string            `json:"last_name"`
    Email          string            `json:"email"`
    Status         string            `json:"status"`
    CreatedAt      time.Time         `json:"created_at"`
    UpdatedAt      time.Time         `json:"updated_at"`
    changedFields  map[string]string `json:"-"`
    isDirty        bool              `json:"-"`
}

// NewUser creates a new user
func NewUser() *User {
    now := time.Now().UTC()
    return &User{
        ID:            uid.HumanUid(),
        Status:        "active",
        CreatedAt:     now,
        UpdatedAt:     now,
        changedFields: make(map[string]string),
        isDirty:       false,
    }
}

// NewUserFromMap creates a user from a map
func NewUserFromMap(data map[string]string) *User {
    user := &User{
        ID:            data["id"],
        FirstName:     data["first_name"],
        LastName:      data["last_name"],
        Email:         data["email"],
        Status:        data["status"],
        changedFields: make(map[string]string),
        isDirty:       false,
    }
    
    if data["created_at"] != "" {
        createdAt, _ := time.Parse(time.RFC3339, data["created_at"])
        user.CreatedAt = createdAt
    }
    
    if data["updated_at"] != "" {
        updatedAt, _ := time.Parse(time.RFC3339, data["updated_at"])
        user.UpdatedAt = updatedAt
    }
    
    return user
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
    return u.Status == "active"
}

// FullName returns the user's full name
func (u *User) FullName() string {
    return u.FirstName + " " + u.LastName
}

// SetFirstName sets the first name and tracks the change
func (u *User) SetFirstName(firstName string) *User {
    if u.FirstName != firstName {
        u.isDirty = true
        u.changedFields["first_name"] = firstName
        u.FirstName = firstName
        u.UpdatedAt = time.Now().UTC()
        u.changedFields["updated_at"] = u.UpdatedAt.Format(time.RFC3339)
    }
    return u
}

// SetLastName sets the last name and tracks the change
func (u *User) SetLastName(lastName string) *User {
    if u.LastName != lastName {
        u.isDirty = true
        u.changedFields["last_name"] = lastName
        u.LastName = lastName
        u.UpdatedAt = time.Now().UTC()
        u.changedFields["updated_at"] = u.UpdatedAt.Format(time.RFC3339)
    }
    return u
}

// SetEmail sets the email and tracks the change
func (u *User) SetEmail(email string) *User {
    if u.Email != email {
        u.isDirty = true
        u.changedFields["email"] = email
        u.Email = email
        u.UpdatedAt = time.Now().UTC()
        u.changedFields["updated_at"] = u.UpdatedAt.Format(time.RFC3339)
    }
    return u
}

// SetStatus sets the status and tracks the change
func (u *User) SetStatus(status string) *User {
    if u.Status != status {
        u.isDirty = true
        u.changedFields["status"] = status
        u.Status = status
        u.UpdatedAt = time.Now().UTC()
        u.changedFields["updated_at"] = u.UpdatedAt.Format(time.RFC3339)
    }
    return u
}

// IsDirty returns whether the user has been modified
func (u *User) IsDirty() bool {
    return u.isDirty
}

// ChangedFields returns the changed fields
func (u *User) ChangedFields() map[string]string {
    return u.changedFields
}

// ResetChanges resets the change tracking
func (u *User) ResetChanges() {
    u.isDirty = false
    u.changedFields = make(map[string]string)
}

// ToJSON converts the user to JSON
func (u *User) ToJSON() (string, error) {
    bytes, err := json.Marshal(u)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

// FromJSON populates the user from JSON
func (u *User) FromJSON(jsonStr string) error {
    return json.Unmarshal([]byte(jsonStr), u)
}

// ToMap converts the user to a map
func (u *User) ToMap() map[string]string {
    return map[string]string{
        "id":         u.ID,
        "first_name": u.FirstName,
        "last_name":  u.LastName,
        "email":      u.Email,
        "status":     u.Status,
        "created_at": u.CreatedAt.Format(time.RFC3339),
        "updated_at": u.UpdatedAt.Format(time.RFC3339),
    }
}
```

### DataObject Example

```go
package models

import (
    "github.com/dracory/dataobject"
    "github.com/dracory/uid"
    "time"
)

// User is a data object
type User struct {
    dataobject.DataObject
}

// NewUser instantiates a new user
func NewUser() *User {
    o := &User{}
    o.SetID(uid.HumanUid())
    o.SetStatus("active")
    o.SetCreatedAt(time.Now().UTC().Format(time.RFC3339))
    o.SetUpdatedAt(time.Now().UTC().Format(time.RFC3339))
    return o
}

// NewUserFromData helper method to hydrate an existing user data object
func NewUserFromData(data map[string]string) *User {
    o := &User{}
    o.Hydrate(data)
    return o
}

// NewUserFromJSON creates a user from JSON string
func NewUserFromJSON(jsonString string) (*User, error) {
    do, err := dataobject.NewFromJSON(jsonString)
    if err != nil {
        return nil, err
    }
    
    return &User{DataObject: *do}, nil
}

// Getters and Setters
func (o *User) FirstName() string {
    return o.Get("first_name")
}

func (o *User) SetFirstName(firstName string) *User {
    o.Set("first_name", firstName)
    o.SetUpdatedAt(time.Now().UTC().Format(time.RFC3339))
    return o
}

func (o *User) LastName() string {
    return o.Get("last_name")
}

func (o *User) SetLastName(lastName string) *User {
    o.Set("last_name", lastName)
    o.SetUpdatedAt(time.Now().UTC().Format(time.RFC3339))
    return o
}

func (o *User) Email() string {
    return o.Get("email")
}

func (o *User) SetEmail(email string) *User {
    o.Set("email", email)
    o.SetUpdatedAt(time.Now().UTC().Format(time.RFC3339))
    return o
}

func (o *User) Status() string {
    return o.Get("status")
}

func (o *User) SetStatus(status string) *User {
    o.Set("status", status)
    o.SetUpdatedAt(time.Now().UTC().Format(time.RFC3339))
    return o
}

func (o *User) CreatedAt() string {
    return o.Get("created_at")
}

func (o *User) SetCreatedAt(createdAt string) *User {
    o.Set("created_at", createdAt)
    return o
}

func (o *User) UpdatedAt() string {
    return o.Get("updated_at")
}

func (o *User) SetUpdatedAt(updatedAt string) *User {
    o.Set("updated_at", updatedAt)
    return o
}

// Helper methods
func (o *User) IsActive() bool {
    return o.Status() == "active"
}

func (o *User) FullName() string {
    return o.FirstName() + " " + o.LastName()
}
```

## Key Differences

1. **Data Storage and Type Safety**:
   - Standard Go Structs: Store data in typed fields defined at compile time
   - DataObject: Stores data in a generic map[string]string, with type conversion in getters/setters

2. **Change Tracking**:
   - Standard Go Structs: No built-in change tracking (requires custom implementation)
   - DataObject: Built-in change tracking with `IsDirty()` and `DataChanged()`

3. **Dynamic Fields**:
   - Standard Go Structs: Fields must be defined at compile time
   - DataObject: Can add any key/value pair at runtime

4. **Memory and Performance**:
   - Standard Go Structs: More memory-efficient and faster for known fields
   - DataObject: Less memory-efficient due to map storage and string conversions

5. **Type Conversion**:
   - Standard Go Structs: No type conversion needed (fields have specific types)
   - DataObject: Requires type conversion in getters/setters (all values stored as strings)

## Use Cases

### Standard Go Structs

- When the data structure is well-defined and unlikely to change
- When performance and memory efficiency are critical
- When strong compile-time type checking is important
- When working with the standard library's JSON encoding/decoding
- For simple data structures without complex behavior

### DataObject

- When change tracking is important (e.g., for partial database updates)
- When the data structure might evolve or have dynamic fields
- When you need a consistent interface for different types of objects
- When working with string-based data sources (e.g., form submissions, URL parameters)
- For domain objects in MVC architecture with complex behavior

## Practical Examples

### Standard Go Struct in a Web Application

```go
// HTTP handler using standard struct
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var requestData struct {
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
        Email     string `json:"email"`
        Status    string `json:"status"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Create new user
    user := NewUser()
    user.SetFirstName(requestData.FirstName)
    user.SetLastName(requestData.LastName)
    user.SetEmail(requestData.Email)
    if requestData.Status != "" {
        user.SetStatus(requestData.Status)
    }
    
    // Save user to database
    if err := db.SaveUser(user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return created user
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

### Standard Go Struct with Change Tracking in a Database Context

```go
// Update user with change tracking
func UpdateUser(id string, updates map[string]string) error {
    // Get existing user
    user, err := GetUserByID(id)
    if err != nil {
        return err
    }
    
    // Apply updates
    if firstName, ok := updates["first_name"]; ok {
        user.SetFirstName(firstName)
    }
    if lastName, ok := updates["last_name"]; ok {
        user.SetLastName(lastName)
    }
    if email, ok := updates["email"]; ok {
        user.SetEmail(email)
    }
    if status, ok := updates["status"]; ok {
        user.SetStatus(status)
    }
    
    // Save changes if any
    if user.IsDirty() {
        // Build SQL for partial update
        query := "UPDATE users SET "
        args := []interface{}{}
        i := 0
        
        changedFields := user.ChangedFields()
        for field, value := range changedFields {
            if i > 0 {
                query += ", "
            }
            
            // Convert field name to database column name
            column := fieldToColumn(field)
            query += column + " = ?"
            args = append(args, value)
            i++
        }
        
        query += " WHERE id = ?"
        args = append(args, id)
        
        // Execute query
        _, err := db.Exec(query, args...)
        if err != nil {
            return err
        }
        
        // Reset change tracking
        user.ResetChanges()
    }
    
    return nil
}
```

### DataObject in a Web Application

```go
// HTTP handler using DataObject
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var requestData map[string]string
    json.NewDecoder(r.Body).Decode(&requestData)
    
    // Create new user
    user := NewUser()
    if firstName, ok := requestData["first_name"]; ok {
        user.SetFirstName(firstName)
    }
    if lastName, ok := requestData["last_name"]; ok {
        user.SetLastName(lastName)
    }
    if email, ok := requestData["email"]; ok {
        user.SetEmail(email)
    }
    if status, ok := requestData["status"]; ok {
        user.SetStatus(status)
    }
    
    // Save user to database
    if err := db.SaveUser(user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return created user
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user.Data())
}
```

### DataObject with Change Tracking in a Database Context

```go
// Update user with DataObject's built-in change tracking
func UpdateUser(id string, updates map[string]string) error {
    // Get existing user
    user, err := GetUserByID(id)
    if err != nil {
        return err
    }
    
    // Apply updates
    if firstName, ok := updates["first_name"]; ok {
        user.SetFirstName(firstName)
    }
    if lastName, ok := updates["last_name"]; ok {
        user.SetLastName(lastName)
    }
    if email, ok := updates["email"]; ok {
        user.SetEmail(email)
    }
    if status, ok := updates["status"]; ok {
        user.SetStatus(status)
    }
    
    // Save changes if any
    if user.IsDirty() {
        // Build SQL for partial update
        query := "UPDATE users SET "
        args := []interface{}{}
        i := 0
        
        changedData := user.DataChanged()
        for field, value := range changedData {
            if i > 0 {
                query += ", "
            }
            
            // Convert field name to database column name
            column := fieldToColumn(field)
            query += column + " = ?"
            args = append(args, value)
            i++
        }
        
        query += " WHERE id = ?"
        args = append(args, id)
        
        // Execute query
        _, err := db.Exec(query, args...)
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

## Integration with MVC Architecture

### Standard Go Structs in MVC

In an MVC architecture using standard Go structs:

1. **Model**: The struct itself represents the data model
   ```go
   type User struct {
       ID        string    `json:"id"`
       FirstName string    `json:"first_name"`
       LastName  string    `json:"last_name"`
       // Other fields...
   }
   ```

2. **Controller**: Handles requests and updates the model
   ```go
   func UserController(w http.ResponseWriter, r *http.Request) {
       // Parse request
       // Update model
       // Render view
   }
   ```

3. **View**: Renders the model data
   ```go
   func RenderUserView(user *User) string {
       // Generate HTML or JSON representation
   }
   ```

### DataObject in MVC

In an MVC architecture using DataObject:

1. **Model**: The DataObject represents the data model
   ```go
   type User struct {
       dataobject.DataObject
   }
   
   func (u *User) FirstName() string {
       return u.Get("first_name")
   }
   
   func (u *User) SetFirstName(firstName string) *User {
       u.Set("first_name", firstName)
       return u
   }
   // Other getters and setters...
   ```

2. **Controller**: Handles requests and updates the model
   ```go
   func UserController(w http.ResponseWriter, r *http.Request) {
       // Parse request
       // Update model using getters/setters
       // Track changes with IsDirty() and DataChanged()
       // Render view
   }
   ```

3. **View**: Renders the model data
   ```go
   func RenderUserView(user *User) string {
       // Generate HTML or JSON representation using getters
   }
   ```

## Conclusion

Standard Go structs and DataObject represent different approaches to data management in Go applications:

- **Standard Go Structs** provide strong compile-time type safety, better performance, and memory efficiency. They are ideal for well-defined data structures where the schema is unlikely to change and performance is critical.

- **DataObject** provides built-in change tracking, dynamic fields, and a consistent interface. It is ideal for applications where tracking changes for partial updates is important, or where the data structure might evolve over time.

In practice, the choice between these approaches depends on your specific requirements:

1. **Use Standard Go Structs when**:
   - Performance and memory efficiency are critical
   - The data structure is well-defined and stable
   - You need strong compile-time type checking
   - You're working with simple data structures

2. **Use DataObject when**:
   - Change tracking for partial updates is important
   - You need a consistent interface across different types of objects
   - The data structure might evolve or have dynamic fields
   - You're working with string-based data sources

Both approaches can be used effectively in an MVC architecture, with DataObject providing additional features that are particularly useful in database-backed applications where efficient updates are important.

For our MVC pattern with independent data stores and controllers, DataObject provides a convenient abstraction that aligns well with the architecture, particularly for database operations where partial updates can significantly improve performance.
