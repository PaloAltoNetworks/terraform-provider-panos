// Package addr is the client.Objects.Address namespace.
//
// Normalized object:  Entry
package addr

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// Constants for Entry.Type field.
const (
    IpNetmask string = "ip-netmask"
    IpRange string = "ip-range"
    Fqdn string = "fqdn"
)

// Entry is a normalized, version independent representation of an address
// object.
type Entry struct {
    Name string
    Value string
    Type string
    Description string
    Tag []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Value = s.Value
    o.Type = s.Type
    o.Description = s.Description
    o.Tag = s.Tag
}

// Addr is a namespace struct, included as part of pango.Client.
type Addr struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *Addr) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of address objects.
func (c *Addr) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of address objects")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of address objects.
func (c *Addr) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of address objects")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given address object.
func (c *Addr) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) address object %q", name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given address object.
func (c *Addr) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) address object %q", name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more address objects.
func (c *Addr) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "address"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) address objects: %v", names)

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

// Edit performs EDIT to create / update an address object.
func (c *Addr) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) address object %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Create the objects.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given address objects from the firewall.
//
// Address objects can be either a string or an Entry object.
func (c *Addr) Delete(vsys string, e ...interface{}) error {
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
    c.con.LogAction("(delete) address objects: %v", names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the Addr struct **/

func (c *Addr) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *Addr) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Addr) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "address",
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
        Description: o.Answer.Description,
        Tag: util.MemToStr(o.Answer.Tag),
    }
    switch {
    case o.Answer.IpNetmask != nil:
        ans.Type = IpNetmask
        ans.Value = o.Answer.IpNetmask.Value
    case o.Answer.IpRange != nil:
        ans.Type = IpRange
        ans.Value = o.Answer.IpRange.Value
    case o.Answer.Fqdn != nil:
        ans.Type = Fqdn
        ans.Value = o.Answer.Fqdn.Value
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    IpNetmask *valType `xml:"ip-netmask"`
    IpRange *valType `xml:"ip-range"`
    Fqdn *valType `xml:"fqdn"`
    Description string `xml:"description"`
    Tag *util.Member `xml:"tag"`
}

type valType struct {
    Value string `xml:",chardata"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        Tag: util.StrToMem(e.Tag),
    }
    vt := &valType{e.Value}
    switch e.Type {
    case IpNetmask:
        ans.IpNetmask = vt
    case IpRange:
        ans.IpRange = vt
    case Fqdn:
        ans.Fqdn = vt
    }

    return ans
}
