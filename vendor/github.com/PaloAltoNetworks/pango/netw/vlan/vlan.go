// Package vlan is the client.Network.Vlan namespace.
//
// Normalized object:  Entry
package vlan

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a VLAN.
//
// Static MAC addresses are given as a map[string] string, where the key is
// the MAC address and the value is the interface it should be associated with.
type Entry struct {
    Name string
    VlanInterface string
    Interfaces []string
    StaticMacs map[string] string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.VlanInterface = s.VlanInterface
    o.Interfaces = s.Interfaces
    o.StaticMacs = s.StaticMacs
}

// Vlan is the client.Network.Vlan namespace.
type Vlan struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *Vlan) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of VLANs.
func (c *Vlan) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of VLANs")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of VLANs.
func (c *Vlan) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of VLANs")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given VLAN.
func (c *Vlan) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) VLAN %q", name)
    return c.details(c.con.Get, name)
}

// Show performs SHOW to retrieve information for the given VLAN.
func (c *Vlan) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) VLAN %q", name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more VLANs.
//
// Specify a non-empty vsys to import the VLAN(s) into the given vsys
// after creating, allowing the vsys to use them.
func (c *Vlan) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given VLAN configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "vlan"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) VLANs: %v", names)

    // Set xpath.
    path := c.xpath(names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the VLANs.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    if err != nil {
        return err
    }

    // Perform vsys import next.
    if vsys == "" {
        return nil
    }
    return c.con.ImportVlans(vsys, names)
}

// Edit performs EDIT to create / update a VLAN.
//
// Specify a non-empty vsys to import the VLAN into the given vsys
// after creating, allowing the vsys to use it.
func (c *Vlan) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) VLAN %q", e.Name)

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit the VLAN.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    if err != nil {
        return err
    }

    // Perform vsys import next.
    if vsys == "" {
        return nil
    }
    return c.con.ImportVlans(vsys, []string{e.Name})
}

// Delete removes the given VLAN(s) from the firewall.
//
// Specify a non-empty vsys to have this function remove the VLAN(s) from
// the vsys prior to deleting them.
//
// VLANs can be a string or an Entry object.
func (c *Vlan) Delete(vsys string, e ...interface{}) error {
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
            return fmt.Errorf("Unknown type sent to delete: %s", v)
        }
    }
    c.con.LogAction("(delete) VLANs: %v", names)

    // Unimport VLANs from the given vsys.
    err = c.con.UnimportVlans(vsys, names)
    if err != nil {
        return err
    }

    // Remove VLANs next.
    path := c.xpath(names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the Vlan struct **/

func (c *Vlan) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *Vlan) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Vlan) xpath(vals []string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "vlan",
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
        VlanInterface: o.Answer.VlanInterface,
        Interfaces: util.MemToStr(o.Answer.Interfaces),
    }
    if len(o.Answer.Mac.Entry) > 0 {
        ans.StaticMacs = make(map[string] string, len(o.Answer.Mac.Entry))
        for i := range o.Answer.Mac.Entry {
            ans.StaticMacs[o.Answer.Mac.Entry[i].Mac] = o.Answer.Mac.Entry[i].Interface
        }
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    VlanInterface string `xml:"virtual-interface>interface"`
    Interfaces *util.Member `xml:"interface"`
    Mac mac `xml:"mac"`
}

type mac struct {
    Entry []macList `xml:"entry"`
}

type macList struct {
    XMLName xml.Name `xml:"entry"`
    Mac string `xml:"name,attr"`
    Interface string `xml:"interface"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        VlanInterface: e.VlanInterface,
        Interfaces: util.StrToMem(e.Interfaces),
    }

    i := 0
    ans.Mac.Entry = make([]macList, len(e.StaticMacs))
    for key := range e.StaticMacs {
        ans.Mac.Entry[i] = macList{Mac: key, Interface: e.StaticMacs[key]}
        i = i + 1
    }

    return ans
}
