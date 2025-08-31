# DataObject vs PHP DTOs: A Comparison

## Introduction

This document compares the Go-based DataObject implementation with PHP Data Transfer Objects (DTOs). While both serve as structured data containers, they have different approaches and features based on their language capabilities and typical use cases.

## Core Concepts

### PHP Data Transfer Objects (DTOs)

PHP DTOs are simple classes that:
- Transport data between system components
- Have no or minimal behavior
- Typically have public or protected properties with getters/setters
- May be immutable (especially in modern PHP)
- Usually don't have built-in change tracking
- Often use PHP attributes/annotations for validation and serialization metadata
- May implement interfaces for specific behaviors

### Go DataObject

Our Go DataObject implementation follows these principles:
- Has a default constructor
- Has a unique identifier (ID)
- All fields are private
- Fields are accessed via public getter and setter methods
- Tracks changes to data internally
- Can return all data as a map
- Designed for efficient serialization to data stores

## Detailed Comparison

| Feature | PHP DTOs | Go DataObject |
|---------|----------|---------------|
| **Language** | PHP | Go |
| **Data Storage** | Class properties | Internal map[string]string |
| **Type Safety** | Varies (weak in PHP 7, stronger with typed properties in PHP 8+) | Strong typing via getters/setters |
| **Constructor** | Typically accepts all properties | Multiple constructors (`New()`, `NewFromData()`, `NewFromJSON()`) |
| **Property Access** | Direct or via getters/setters | Via explicit getter/setter methods |
| **Immutability** | Optional, common in modern PHP DTOs | Mutable by design |
| **Change Tracking** | Not typically built-in | Built-in via `IsDirty()` and `DataChanged()` |
| **Unique Identifier** | Not required (but common) | Required (ID field) |
| **Serialization** | Via PHP's built-in functions or libraries like Symfony Serializer | Custom JSON serialization |
| **Validation** | Often uses attributes/annotations or separate validators | Not built-in, typically handled separately |
| **Method Chaining** | Sometimes supported | Supported by design |

## Implementation Differences

### PHP DTO Example (Traditional Style)

```php
<?php

namespace App\DTOs;

class UserDTO
{
    private string $id;
    private string $firstName;
    private string $lastName;
    private string $email;
    private string $status;
    
    public function __construct(
        string $id,
        string $firstName,
        string $lastName,
        string $email,
        string $status = 'active'
    ) {
        $this->id = $id;
        $this->firstName = $firstName;
        $this->lastName = $lastName;
        $this->email = $email;
        $this->status = $status;
    }
    
    public function getId(): string
    {
        return $this->id;
    }
    
    public function getFirstName(): string
    {
        return $this->firstName;
    }
    
    public function getLastName(): string
    {
        return $this->lastName;
    }
    
    public function getEmail(): string
    {
        return $this->email;
    }
    
    public function getStatus(): string
    {
        return $this->status;
    }
    
    public function isActive(): bool
    {
        return $this->status === 'active';
    }
    
    public function toArray(): array
    {
        return [
            'id' => $this->id,
            'first_name' => $this->firstName,
            'last_name' => $this->lastName,
            'email' => $this->email,
            'status' => $this->status,
        ];
    }
    
    public static function fromArray(array $data): self
    {
        return new self(
            $data['id'] ?? '',
            $data['first_name'] ?? '',
            $data['last_name'] ?? '',
            $data['email'] ?? '',
            $data['status'] ?? 'active'
        );
    }
}
```

### PHP DTO Example (Modern Immutable Style with PHP 8)

```php
<?php

namespace App\DTOs;

use JsonSerializable;

final class UserDTO implements JsonSerializable
{
    public function __construct(
        public readonly string $id,
        public readonly string $firstName,
        public readonly string $lastName,
        public readonly string $email,
        public readonly string $status = 'active'
    ) {
    }
    
    public function isActive(): bool
    {
        return $this->status === 'active';
    }
    
    public function withFirstName(string $firstName): self
    {
        return new self(
            $this->id,
            $firstName,
            $this->lastName,
            $this->email,
            $this->status
        );
    }
    
    public function withLastName(string $lastName): self
    {
        return new self(
            $this->id,
            $this->firstName,
            $lastName,
            $this->email,
            $this->status
        );
    }
    
    public function jsonSerialize(): array
    {
        return [
            'id' => $this->id,
            'first_name' => $this->firstName,
            'last_name' => $this->lastName,
            'email' => $this->email,
            'status' => $this->status,
        ];
    }
    
    public static function fromArray(array $data): self
    {
        return new self(
            $data['id'] ?? '',
            $data['first_name'] ?? '',
            $data['last_name'] ?? '',
            $data['email'] ?? '',
            $data['status'] ?? 'active'
        );
    }
}
```

### Go DataObject Example

```go
package models

import (
    "github.com/dracory/dataobject"
    "github.com/gouniverse/uid"
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
    return o
}

// NewUserFromData helper method to hydrate an existing user data object
func NewUserFromData(data map[string]string) *User {
    o := &User{}
    o.Hydrate(data)
    return o
}

// Getters and Setters
func (o *User) FirstName() string {
    return o.Get("first_name")
}

func (o *User) SetFirstName(firstName string) *User {
    o.Set("first_name", firstName)
    return o
}

func (o *User) LastName() string {
    return o.Get("last_name")
}

func (o *User) SetLastName(lastName string) *User {
    o.Set("last_name", lastName)
    return o
}

func (o *User) Email() string {
    return o.Get("email")
}

func (o *User) SetEmail(email string) *User {
    o.Set("email", email)
    return o
}

func (o *User) Status() string {
    return o.Get("status")
}

func (o *User) SetStatus(status string) *User {
    o.Set("status", status)
    return o
}

// Helper methods
func (o *User) IsActive() bool {
    return o.Status() == "active"
}
```

## Key Differences

1. **Immutability vs. Mutability**:
   - Modern PHP DTOs: Often designed to be immutable with readonly properties
   - DataObject: Designed to be mutable with change tracking

2. **Data Storage Approach**:
   - PHP DTOs: Each property is a separate typed field
   - DataObject: All data is stored in a generic map[string]string

3. **Change Tracking**:
   - PHP DTOs: No built-in change tracking
   - DataObject: Built-in change tracking with `IsDirty()` and `DataChanged()`

4. **Property Updates**:
   - Modern PHP DTOs: Often use "with" methods that return new instances
   - DataObject: Uses setter methods that modify the existing instance

5. **Type Safety**:
   - PHP DTOs: Type safety varies (stronger in PHP 8+ with typed properties)
   - DataObject: Type conversion is handled by getters/setters

## Use Cases

### PHP DTOs

- API request/response objects
- Form data validation
- Cross-boundary data transfer
- Command/Query objects in CQRS
- Value objects in Domain-Driven Design
- Immutable data structures

### Go DataObject

- Database entity representation
- Data Transfer Objects in web services
- Domain objects in MVC architecture
- Efficient data storage with change tracking for partial updates

## Practical Examples

### PHP DTO in Laravel API

```php
// Controller
public function updateUser(UpdateUserRequest $request, string $id)
{
    // Request data is validated by Laravel
    $userDTO = new UserDTO(
        $id,
        $request->input('first_name'),
        $request->input('last_name'),
        $request->input('email'),
        $request->input('status', 'active')
    );
    
    // Pass DTO to service layer
    $user = $this->userService->updateUser($userDTO);
    
    return response()->json($user);
}
```

### PHP Immutable DTO with Symfony

```php
// Using an immutable DTO with "with" methods
$userDTO = new UserDTO('123', 'John', 'Doe', 'john@example.com');

// Create a new instance with modified data
$updatedDTO = $userDTO
    ->withFirstName('Jane')
    ->withLastName('Smith');

// Original DTO remains unchanged
echo $userDTO->firstName; // "John"
echo $updatedDTO->firstName; // "Jane"
```

### Go DataObject with Change Tracking

```go
// Create a new user
user := NewUser()
user.SetFirstName("John")
user.SetLastName("Smith")

// Later, update some fields
user.SetFirstName("Jane")

// Check if the object has been modified
if user.IsDirty() {
    // Get only the changed fields
    changedData := user.DataChanged()
    
    // Use the changed data for a partial update
    db.UpdateUserPartial(user.ID(), changedData)
}
```

## Framework Integration

### PHP DTOs in Modern Frameworks

PHP DTOs are often integrated with frameworks in various ways:

```php
// Symfony with attributes for validation and serialization
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Serializer\Annotation\Groups;

class UserDTO
{
    #[Assert\NotBlank]
    #[Assert\Uuid]
    #[Groups(['user:read'])]
    public readonly string $id;
    
    #[Assert\NotBlank]
    #[Assert\Length(min: 2, max: 50)]
    #[Groups(['user:read', 'user:write'])]
    public readonly string $firstName;
    
    // Other properties...
}
```

### Go DataObject in Web Applications

```go
// HTTP handler using DataObject
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var requestData map[string]string
    json.NewDecoder(r.Body).Decode(&requestData)
    
    // Get user ID from URL
    id := chi.URLParam(r, "id")
    
    // Load existing user
    user := GetUserByID(id)
    if user == nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    // Update user with request data
    if firstName, ok := requestData["first_name"]; ok {
        user.SetFirstName(firstName)
    }
    if lastName, ok := requestData["last_name"]; ok {
        user.SetLastName(lastName)
    }
    
    // Save only changed fields
    if user.IsDirty() {
        err := SaveUserChanges(user)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    
    // Return updated user
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user.Data())
}
```

## Conclusion

PHP DTOs and Go DataObject represent different approaches to structured data containers:

- **PHP DTOs** have evolved from simple data carriers to more sophisticated immutable objects with rich metadata through attributes/annotations. Modern PHP DTOs often emphasize immutability and type safety, especially with PHP 8's features.

- **Go DataObject** takes a more dynamic approach with a generic data store and built-in change tracking. It emphasizes mutability with tracking, making it particularly well-suited for database operations where partial updates are important.

The choice between these approaches depends on your specific requirements:

- If you need immutable data structures with strong validation through framework integration, modern PHP DTOs are well-suited.

- If you need mutable objects with built-in change tracking for efficient database operations, Go DataObject provides these features out of the box.

Both approaches have their place in modern application development, and the best choice depends on your language ecosystem, application requirements, and architectural preferences.
