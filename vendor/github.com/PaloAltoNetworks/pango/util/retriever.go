package util

// Retriever is a type that is intended to act as a stand-in for using
// either the Get or Show pango Client functions.
type Retriever func(interface{}, interface{}, interface{}) ([]byte, error)
