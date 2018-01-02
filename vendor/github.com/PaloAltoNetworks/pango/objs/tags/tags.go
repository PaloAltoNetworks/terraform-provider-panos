// Package tags is the client.Objects.Tags namespace.
//
// Normalized object:  Entry
package tags

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
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
    Name string
    Color string
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

// Tags is a namespace struct, included as part of pango.Client.
type Tags struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *Tags) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of administrative tags.
func (c *Tags) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of administrative tags")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of administrative tags.
func (c *Tags) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of administrative tags")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given administrative tag.
func (c *Tags) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) administrative tag %q", name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given administrative tag.
func (c *Tags) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) administrative tag %q", name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more administrative tags.
func (c *Tags) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "tag"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) administrative tags: %v", names)

    // Set xpath.
    path := c.xpath(vsys, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the objects.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update an administrative tag.
func (c *Tags) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) administrative tag %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Create the objects.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given administrative tags from the firewall.
//
// Administrative tags can be either a string or an Entry object.
func (c *Tags) Delete(vsys string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    names := make([]string, len(e))
    for i := range e {
        switch v := e[i].(type) {
        case string:
            names[i] = v
        case Entry:
            names[i] = v.Name
        default:
            return fmt.Errorf("Unsupported type to delete: %s", v)
        }
    }
    c.con.LogAction("(delete) administrative tags: %v", names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the Tags struct **/

func (c *Tags) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *Tags) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Tags) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "tag",
        util.AsEntryXpath(vals),
    }
}

/** Structs / functions for this namespace. **/

type normalizer interface {
    Normalize() Entry
}

type container_v1 struct {
    Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Color: o.Answer.Color,
        Comment: o.Answer.Comment,
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Color string `xml:"color,omitempty"`
    Comment string `xml:"comments,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Color: e.Color,
        Comment: e.Comment,
    }

    return ans
}
