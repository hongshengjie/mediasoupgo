package mediasoupgo

// AppData represents a map of string keys to arbitrary values.
type AppData map[string]interface{}

// Either represents a type that can hold either one of two types of data, but not both.
// In Go, this cannot be enforced at compile time as in TypeScript.
// Instead, this is a placeholder type (interface{}), and the mutual exclusivity
// of the fields must be enforced at runtime by the application logic.
// For example, if you intend to use Either with two specific structs T and U,
// you can define them explicitly and enforce the constraint in your code.
type Either interface{}
