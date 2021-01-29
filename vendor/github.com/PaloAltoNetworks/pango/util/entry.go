package util

import (
	"encoding/xml"
)

// EntryType defines an entry config node used for sending and receiving XML
// from PAN-OS.
type EntryType struct {
	Entries []Entry `xml:"entry"`
}

// Entry is a standalone entry struct.
type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Value   string   `xml:"name,attr"`
}

// EntToStr normalizes an EntryType pointer into a list of strings.
func EntToStr(e *EntryType) []string {
	if e == nil {
		return nil
	}

	ans := make([]string, len(e.Entries))
	for i := range e.Entries {
		ans[i] = e.Entries[i].Value
	}

	return ans
}

// StrToEnt converts a list of strings into an EntryType pointer.
func StrToEnt(e []string) *EntryType {
	if e == nil {
		return nil
	}

	ans := make([]Entry, len(e))
	for i := range e {
		ans[i] = Entry{Value: e[i]}
	}

	return &EntryType{ans}
}

// EntToOneStr normalizes an EntryType pointer for a max_items=1 XML node
// into a string.
func EntToOneStr(e *EntryType) string {
	if e == nil || len(e.Entries) == 0 {
		return ""
	}

	return e.Entries[0].Value
}

// OneStrToEnt converts a string into an EntryType pointer for a max_items=1
// XML node.
func OneStrToEnt(e string) *EntryType {
	if e == "" {
		return nil
	}

	return &EntryType{[]Entry{
		{Value: e},
	}}
}
