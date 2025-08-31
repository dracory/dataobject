# DataObject vs Elixir Structs: A Comparison

## Introduction

This document compares the Go-based DataObject implementation with Elixir Structs. While both serve as data containers in their respective languages, they have fundamentally different approaches based on their underlying language paradigms: Go's imperative, mutable approach versus Elixir's functional, immutable approach.

## Core Concepts

### Elixir Structs

Elixir Structs are:
- Named map-like data structures with a defined set of fields
- Immutable by design (following Elixir's functional paradigm)
- Compile-time checked with defined field names
- Updated via transformation functions that return new instances
- Often used with pattern matching
- Typically extended with behavior via protocols
- Can include default values for fields
- Defined within modules

### Go DataObject

Our Go DataObject implementation follows these principles:
- Has a default constructor
- Has a unique identifier (ID)
- All fields are private (stored in an internal map)
- Fields are accessed via public getter and setter methods
- Tracks changes to data internally
- Can return all data as a map
- Designed for efficient serialization to data stores
- Mutable by design

## Detailed Comparison

| Feature | Elixir Structs | Go DataObject |
|---------|----------------|---------------|
| **Language Paradigm** | Functional, immutable | Imperative, mutable |
| **Data Storage** | Named fields with compile-time checking | Internal map[string]string |
| **Type Safety** | Optional static typing with typespecs | Strong typing via getters/setters |
| **Constructor** | `%ModuleName{}` syntax | Multiple constructors (`New()`, `NewFromData()`, `NewFromJSON()`) |
| **Property Access** | Direct field access with dot notation | Via explicit getter/setter methods |
| **Immutability** | Immutable by design | Mutable by design |
| **Change Tracking** | Not needed (immutable) | Built-in via `IsDirty()` and `DataChanged()` |
| **Unique Identifier** | Not required | Required (ID field) |
| **Serialization** | Via libraries like Jason or Poison | Custom JSON serialization |
| **Method Generation** | No automatic method generation | Manual implementation |
| **Default Values** | Supported in struct definition | Set in constructors |
| **Inheritance** | No inheritance (composition via protocols) | Composition-based |
| **Method Chaining** | Not applicable (immutable) | Supported by design |
| **Pattern Matching** | First-class support | Not supported |

## Implementation Differences

### Elixir Struct Example (Basic)

```elixir
defmodule User do
  @derive {Jason.Encoder, only: [:id, :first_name, :last_name, :email, :status, :created_at, :updated_at]}
  
  defstruct id: nil,
            first_name: "",
            last_name: "",
            email: "",
            status: "active",
            created_at: nil,
            updated_at: nil

  def new do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %__MODULE__{
      id: UUID.uuid4(),
      created_at: now,
      updated_at: now
    }
  end
  
  def from_map(data) when is_map(data) do
    # Convert string keys to atoms safely
    data_with_atoms = for {key, val} <- data, into: %{} do
      atom_key = if is_atom(key), do: key, else: String.to_existing_atom(key)
      {atom_key, val}
    end
    
    struct(__MODULE__, data_with_atoms)
  end
  
  def is_active(%__MODULE__{status: "active"}), do: true
  def is_active(_), do: false
  
  def full_name(%__MODULE__{first_name: first_name, last_name: last_name}) do
    "#{first_name} #{last_name}"
  end
  
  def set_first_name(user, first_name) do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %{user | first_name: first_name, updated_at: now}
  end
  
  def set_last_name(user, last_name) do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %{user | last_name: last_name, updated_at: now}
  end
  
  def set_email(user, email) do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %{user | email: email, updated_at: now}
  end
  
  def set_status(user, status) do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %{user | status: status, updated_at: now}
  end
  
  def to_map(user) do
    Map.from_struct(user)
  end
  
  def to_json(user) do
    Jason.encode!(user)
  end
  
  def from_json(json) when is_binary(json) do
    json
    |> Jason.decode!()
    |> from_map()
  end
end
```

### Elixir Struct Example (With Change Tracking)

```elixir
defmodule User do
  @derive {Jason.Encoder, only: [:id, :first_name, :last_name, :email, :status, :created_at, :updated_at]}
  
  defstruct id: nil,
            first_name: "",
            last_name: "",
            email: "",
            status: "active",
            created_at: nil,
            updated_at: nil,
            changed_fields: %{}

  def new do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %__MODULE__{
      id: UUID.uuid4(),
      created_at: now,
      updated_at: now,
      changed_fields: %{}
    }
  end
  
  def from_map(data) when is_map(data) do
    # Convert string keys to atoms safely
    data_with_atoms = for {key, val} <- data, into: %{} do
      atom_key = if is_atom(key), do: key, else: String.to_existing_atom(key)
      {atom_key, val}
    end
    
    # Ensure changed_fields is initialized
    data_with_atoms = Map.put_new(data_with_atoms, :changed_fields, %{})
    
    struct(__MODULE__, data_with_atoms)
  end
  
  def is_active(%__MODULE__{status: "active"}), do: true
  def is_active(_), do: false
  
  def full_name(%__MODULE__{first_name: first_name, last_name: last_name}) do
    "#{first_name} #{last_name}"
  end
  
  def set_first_name(%__MODULE__{first_name: current} = user, current), do: user
  def set_first_name(user, first_name) do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %{user | 
      first_name: first_name, 
      updated_at: now,
      changed_fields: Map.put(user.changed_fields, :first_name, first_name)
    }
  end
  
  def set_last_name(%__MODULE__{last_name: current} = user, current), do: user
  def set_last_name(user, last_name) do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %{user | 
      last_name: last_name, 
      updated_at: now,
      changed_fields: Map.put(user.changed_fields, :last_name, last_name)
    }
  end
  
  def set_email(%__MODULE__{email: current} = user, current), do: user
  def set_email(user, email) do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %{user | 
      email: email, 
      updated_at: now,
      changed_fields: Map.put(user.changed_fields, :email, email)
    }
  end
  
  def set_status(%__MODULE__{status: current} = user, current), do: user
  def set_status(user, status) do
    now = DateTime.utc_now() |> DateTime.to_iso8601()
    %{user | 
      status: status, 
      updated_at: now,
      changed_fields: Map.put(user.changed_fields, :status, status)
    }
  end
  
  def is_dirty(%__MODULE__{changed_fields: changed_fields}) do
    map_size(changed_fields) > 0
  end
  
  def changed_fields(%__MODULE__{changed_fields: changed_fields}) do
    # Convert atom keys to strings for compatibility with DataObject
    for {key, val} <- changed_fields, into: %{} do
      {Atom.to_string(key), val}
    end
  end
  
  def reset_changes(user) do
    %{user | changed_fields: %{}}
  end
  
  def to_map(user) do
    user
    |> Map.from_struct()
    |> Map.drop([:changed_fields])
  end
  
  def to_string_map(user) do
    # Convert atom keys to strings
    for {key, val} <- to_map(user), into: %{} do
      {Atom.to_string(key), to_string(val)}
    end
  end
  
  def to_json(user) do
    user
    |> Map.from_struct()
    |> Map.drop([:changed_fields])
    |> Jason.encode!()
  end
  
  def from_json(json) when is_binary(json) do
    json
    |> Jason.decode!()
    |> from_map()
  end
end
```

### Go DataObject Example

```go
package models

import (
    "github.com/dracory/dataobject"
    "github.com/gouniverse/uid"
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

1. **Mutability vs. Immutability**:
   - Elixir Structs: Immutable by design; updates create new instances
   - DataObject: Mutable by design; updates modify the existing instance

2. **Functional vs. Imperative Paradigm**:
   - Elixir Structs: Follow functional programming principles with transformation functions
   - DataObject: Follow imperative programming principles with state modification

3. **Change Tracking**:
   - Elixir Structs: Not typically needed due to immutability (but can be implemented)
   - DataObject: Built-in via `IsDirty()` and `DataChanged()`

4. **Pattern Matching**:
   - Elixir Structs: First-class support for pattern matching
   - DataObject: No pattern matching support

5. **Method Chaining**:
   - Elixir Structs: Not idiomatic due to immutability
   - DataObject: Supported by design for fluent interfaces

## Use Cases

### Elixir Structs

- Domain modeling in functional applications
- Data transformation pipelines
- Concurrent and distributed systems
- Event-driven architectures
- When immutability is a requirement
- Phoenix web applications

### Go DataObject

- Database entity representation
- Data Transfer Objects in web services
- Domain objects in MVC architecture
- Efficient data storage with change tracking for partial updates
- When mutable state with change tracking is needed

## Practical Examples

### Elixir Struct (Basic Usage)

```elixir
# Creating a new user
user = User.new()
|> User.set_first_name("John")
|> User.set_last_name("Doe")
|> User.set_email("john@example.com")

# Accessing attributes
IO.puts(user.first_name)  # "John"
IO.puts(User.full_name(user))  # "John Doe"

# Modifying attributes (creates a new struct)
updated_user = User.set_first_name(user, "Jane")
IO.puts(User.full_name(updated_user))  # "Jane Doe"
IO.puts(User.full_name(user))  # Still "John Doe" (original unchanged)

# Converting to map
user_map = User.to_map(user)
```

### Elixir Struct (With Change Tracking)

```elixir
# Creating a new user
user = User.new()
|> User.set_first_name("John")
|> User.set_last_name("Doe")
|> User.set_email("john@example.com")

# Checking for changes
if User.is_dirty(user) do
  changes = User.changed_fields(user)
  IO.inspect(changes)  # %{"first_name" => "John", "last_name" => "Doe", "email" => "john@example.com"}
  
  # In a database context
  Repo.update_user(user.id, changes)
  
  # Reset changes after saving
  user = User.reset_changes(user)
end

# Converting to JSON
json_str = User.to_json(user)
IO.puts(json_str)

# Creating from JSON
new_user = User.from_json(~s({"id": "123", "first_name": "Alice", "last_name": "Brown", "email": "alice@example.com"}))
```

### Go DataObject

```go
// Create a new user
user := NewUser()
user.SetFirstName("John")
user.SetLastName("Doe")
user.SetEmail("john@example.com")

// Accessing data via getters
fmt.Println(user.FirstName())  // "John"
fmt.Println(user.FullName())   // "John Doe"

// Method chaining for setters
user.SetFirstName("Jane").SetLastName("Smith")

// Checking for changes
if user.IsDirty() {
    changes := user.DataChanged()
    
    // In a database context
    db.UpdateUserPartial(user.ID(), changes)
}

// Converting to JSON
jsonStr := user.ToJSON()

// Creating from JSON
newUser, err := NewUserFromJSON(`{"id":"123","first_name":"Alice","last_name":"Brown","email":"alice@example.com"}`)
if err != nil {
    log.Fatal(err)
}
```

## Framework Integration

### Elixir Structs with Phoenix Framework

```elixir
# Phoenix controller using Elixir structs
defmodule MyApp.UserController do
  use MyApp.Web, :controller
  alias MyApp.User
  
  def create(conn, %{"user" => user_params}) do
    user = User.new()
    |> User.set_first_name(user_params["first_name"])
    |> User.set_last_name(user_params["last_name"])
    |> User.set_email(user_params["email"])
    
    case Repo.insert(user) do
      {:ok, user} ->
        conn
        |> put_status(:created)
        |> render("show.json", user: user)
      {:error, changeset} ->
        conn
        |> put_status(:unprocessable_entity)
        |> render("error.json", changeset: changeset)
    end
  end
  
  def update(conn, %{"id" => id, "user" => user_params}) do
    with {:ok, user} <- Repo.get_user(id) do
      # Apply updates (each function returns a new struct)
      updated_user = user
      |> maybe_update_field(user_params, "first_name", &User.set_first_name/2)
      |> maybe_update_field(user_params, "last_name", &User.set_last_name/2)
      |> maybe_update_field(user_params, "email", &User.set_email/2)
      |> maybe_update_field(user_params, "status", &User.set_status/2)
      
      # Save changes if any
      if User.is_dirty(updated_user) do
        changes = User.changed_fields(updated_user)
        case Repo.update_user(id, changes) do
          {:ok, _} ->
            conn
            |> render("show.json", user: updated_user)
          {:error, reason} ->
            conn
            |> put_status(:unprocessable_entity)
            |> render("error.json", error: reason)
        end
      else
        conn
        |> render("show.json", user: user)
      end
    else
      nil ->
        conn
        |> put_status(:not_found)
        |> render("error.json", error: "User not found")
    end
  end
  
  # Helper function to conditionally update a field
  defp maybe_update_field(user, params, field, update_fn) do
    case Map.get(params, field) do
      nil -> user
      value -> update_fn.(user, value)
    end
  end
end
```

### Go DataObject with Web Framework

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
    
    // Save user
    err := SaveUser(user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return created user
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user.Data())
}

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
    if email, ok := requestData["email"]; ok {
        user.SetEmail(email)
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

## Philosophical Differences

The contrast between Elixir Structs and Go DataObject reflects deeper philosophical differences in programming paradigms:

1. **Functional vs. Imperative**:
   - Elixir embraces functional programming with immutability and transformation
   - Go embraces imperative programming with mutable state and in-place modification

2. **Concurrency Models**:
   - Elixir's immutable structs work well with its actor-based concurrency (no shared state)
   - Go's mutable DataObject requires careful handling in concurrent contexts (locks, channels)

3. **Error Handling**:
   - Elixir uses pattern matching and the "let it crash" philosophy
   - Go uses explicit error returns and checks

4. **Type Systems**:
   - Elixir has dynamic typing with optional typespecs
   - Go has static typing with compile-time checks

## Conclusion

Elixir Structs and Go DataObject represent fundamentally different approaches to data containers, reflecting the core philosophies of their respective languages:

- **Elixir Structs** embrace immutability and functional programming. They create new instances for every change, making them ideal for concurrent systems where avoiding shared mutable state is important. They excel in transformation pipelines and pattern matching scenarios.

- **Go DataObject** embraces mutability with controlled change tracking. It modifies data in place while keeping track of what changed, making it ideal for database operations where partial updates are important. It excels in imperative programming contexts where performance and memory efficiency are priorities.

The choice between these approaches depends on your language ecosystem, application requirements, and architectural preferences:

1. **Choose Elixir Structs when**:
   - You're working in a functional programming paradigm
   - Immutability is a requirement
   - Concurrency is a major concern
   - You need pattern matching capabilities
   - You're building transformation pipelines

2. **Choose Go DataObject when**:
   - You're working in an imperative programming paradigm
   - You need mutable objects with change tracking
   - You're performing partial database updates
   - You're building an MVC architecture with data stores
   - You need a consistent interface for different types of objects

Both approaches can be effective in their respective ecosystems, and the best choice depends on the specific requirements and constraints of your application.
