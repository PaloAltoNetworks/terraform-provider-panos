package util

import (
	"encoding/xml"
)

// BulkElement is a generic bulk container for bulk operations.
type BulkElement struct {
	XMLName xml.Name
	Data    []interface{}
}

// Config returns an interface to be Marshaled.
func (o BulkElement) Config() interface{} {
	if len(o.Data) == 1 {
		return o.Data[0]
	}
	return o
}
