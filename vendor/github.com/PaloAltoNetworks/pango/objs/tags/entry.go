package tags

import (
	"encoding/xml"
	"fmt"
)

// These are the color constants you can use in Entry.SetColor().  Note that
// each version of PANOS has added colors, so if you are looking for maximum
// compatibility, only use the first 16 colors (17 including None).
const (
	None = iota
	Red
	Green
	Blue
	Yellow
	Copper
	Orange
	Purple
	Gray
	LightGreen
	Cyan
	LightGray
	BlueGray
	Lime
	Black
	Gold
	Brown
	Olive
	_
	Maroon
	RedOrange
	YellowOrange
	ForestGreen
	TurquoiseBlue
	AzureBlue
	CeruleanBlue
	MidnightBlue
	MediumBlue
	CobaltBlue
	VioletBlue
	BlueViolet
	MediumViolet
	MediumRose
	Lavender
	Orchid
	Thistle
	Peach
	Salmon
	Magenta
	RedViolet
	Mahogany
	BurntSienna
	Chestnut
)

// Entry is a normalized, version independent representation of an
// administrative tag.  Note that colors should be set to a string
// such as `color5` or `color13`.  If you want to set a color using the
// color name (e.g. - "red"), use the SetColor function.
type Entry struct {
	Name    string
	Color   string
	Comment string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Color = s.Color
	o.Comment = s.Comment
}

// SetColor takes a color constant (e.g. - Olive) and converts it to a color
// enum (e.g. - "color17").
//
// Note that color availability varies according to version:
//
// * 6.1 - 7.0:  None - Brown
// * 7.1 - 8.0:  None - Olive
// * 8.1:  None - Chestnut
func (o *Entry) SetColor(v int) {
	if v == 0 {
		o.Color = ""
	} else {
		o.Color = fmt.Sprintf("color%d", v)
	}
}

/** Structs / functions for normalization. **/

type normalizer interface {
	Normalize() Entry
}

type container_v1 struct {
	Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
	ans := Entry{
		Name:    o.Answer.Name,
		Color:   o.Answer.Color,
		Comment: o.Answer.Comment,
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	Color   string   `xml:"color,omitempty"`
	Comment string   `xml:"comments,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:    e.Name,
		Color:   e.Color,
		Comment: e.Comment,
	}

	return ans
}
