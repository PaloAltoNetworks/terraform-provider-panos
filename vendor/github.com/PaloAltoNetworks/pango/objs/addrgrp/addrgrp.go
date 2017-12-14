// Package addrgrp is the client.Objects.AddressGroup namespace.
//
// Normalized object:  Entry
package addrgrp

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of an address
// group.  The value set in Dynamic should be something like the following:
//
//  * 'tag1'
//  * 'tag1' or 'tag2' and 'tag3'
//
// The tags param is for administrative tags for this address object
// group itself.
type Entry struct {
    Name string
    Description string
    Static []string
    Dynamic string
    Tag []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.Static = s.Static
    o.Dynamic = s.Dynamic
    o.Tag = s.Tag
}

// AddrGrp is a namespace struct, included as part of pango.Client.
type AddrGrp struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *AddrGrp) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of address groups.
func (c *AddrGrp) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of address groups")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of address groups.
func (c *AddrGrp) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of address groups")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given address group.
func (c *AddrGrp) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) address group %q", name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given address group.
func (c *AddrGrp) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) address group %q", name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more address groups.
func (c *AddrGrp) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "address-group"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) address groups: %v", names)

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

// Edit performs EDIT to create / update an address group.
func (c *AddrGrp) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) address group %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Create the objects.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given address groups from the firewall.
//
// Address groups can be either a string or an Entry object.
func (c *AddrGrp) Delete(vsys string, e ...interface{}) error {
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
    c.con.LogAction("(delete) address groups: %v", names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the AddrGrp struct **/

func (c *AddrGrp) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *AddrGrp) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *AddrGrp) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "address-group",
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
        Static: util.MemToStr(o.Answer.Static),
        Tag: util.MemToStr(o.Answer.Tag),
    }
    if o.Answer.Dynamic != nil {
        ans.Dynamic = *o.Answer.Dynamic
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Description string `xml:"description"`
    Static *util.Member `xml:"static"`
    Dynamic *string `xml:"dynamic>filter"`
    Tag *util.Member `xml:"tag"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        Static: util.StrToMem(e.Static),
        Tag: util.StrToMem(e.Tag),
    }
    if e.Dynamic != "" {
        ans.Dynamic = &e.Dynamic
    }

    return ans
}
