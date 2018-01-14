// Package srvcgrp is the client.Objects.ServiceGroup namespace.
//
// Normalized object:  Entry
package srvcgrp

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a service
// group.
type Entry struct {
    Name string
    Services []string
    Tags []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Services = s.Services
    o.Tags = s.Tags
}

// SrvcGrp is a namespace struct, included as part of pango.Client.
type SrvcGrp struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *SrvcGrp) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of service groups.
func (c *SrvcGrp) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of service groups")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of service groups.
func (c *SrvcGrp) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of service groups")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given service group.
func (c *SrvcGrp) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) service group %q", name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given service group.
func (c *SrvcGrp) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) service group %q", name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more service groups.
func (c *SrvcGrp) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "service-group"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) service groups: %v", names)

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

// Edit performs EDIT to create / update a service group.
func (c *SrvcGrp) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) service group %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Create the objects.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given service groups from the firewall.
//
// Service groups can be either a string or an Entry object.
func (c *SrvcGrp) Delete(vsys string, e ...interface{}) error {
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
    c.con.LogAction("(delete) service groups: %v", names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the SrvcGrp struct **/

func (c *SrvcGrp) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *SrvcGrp) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *SrvcGrp) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "service-group",
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
        Services: util.MemToStr(o.Answer.Services),
        Tags: util.MemToStr(o.Answer.Tags),
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Services *util.Member `xml:"members"`
    Tags *util.Member `xml:"tag"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Services: util.StrToMem(e.Services),
        Tags: util.StrToMem(e.Tags),
    }

    return ans
}
