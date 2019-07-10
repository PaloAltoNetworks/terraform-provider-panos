package vlan

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// FwVlan is the client.Network.Vlan namespace.
type FwVlan struct {
    con util.XapiClient
}

/*
SetInterface performs a SET to add an interface to a VLAN.

The VLAN can be either a string or an Entry object.
*/
func (c *FwVlan) SetInterface(vlan interface{}, iface string) error {
    var name string

    switch v := vlan.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to %s set interface: %s", singular, v)
    }

    c.con.LogAction("(set) interface for %s %q: %s", singular, name, iface)

    path := c.xpath([]string{name})
    path = append(path, "interface")

    _, err := c.con.Set(path, util.Member{Value: iface}, nil, nil)
    return err
}

/*
DeleteInterface performs a DELETE to remove an interface from a VLAN.

The VLAN can be either a string or an Entry object.
*/
func (c *FwVlan) DeleteInterface(vlan interface{}, iface string) error {
    var name string

    switch v := vlan.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to %s delete interface: %s", singular, v)
    }

    c.con.LogAction("(delete) interface for %s %q: %s", singular, name, iface)

    path := c.xpath([]string{name})
    path = append(path, "interface", util.AsMemberXpath([]string{iface}))

    _, err := c.con.Delete(path, nil, nil)
    return err
}

// Initialize is invoked by client.Initialize().
func (c *FwVlan) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of VLANs.
func (c *FwVlan) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of VLANs.
func (c *FwVlan) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given VLAN.
func (c *FwVlan) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, name)
}

// Show performs SHOW to retrieve information for the given VLAN.
func (c *FwVlan) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more VLANs.
//
// Specify a non-empty vsys to import the VLAN(s) into the given vsys
// after creating, allowing the vsys to use them.
func (c *FwVlan) Set(vsys string, e ...Entry) error {
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
    c.con.LogAction("(set) %s: %v", plural, names)

    // Set xpath.
    path := c.xpath(names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the VLANs.
    if _, err = c.con.Set(path, d.Config(), nil, nil); err != nil {
        return err
    }

    // Remove the VLANs from any vsys they're currently in.
    if err = c.con.VsysUnimport(util.VlanImport, "", "", names); err != nil {
        return err
    }

    // Perform vsys import next.
    return c.con.VsysImport(util.VlanImport, "", "", vsys, names)
}

// Edit performs EDIT to create / update a VLAN.
//
// Specify a non-empty vsys to import the VLAN into the given vsys
// after creating, allowing the vsys to use it.
func (c *FwVlan) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit the VLAN.
    if _, err = c.con.Edit(path, fn(e), nil, nil); err != nil {
        return err
    }

    // Remove the VLANs from any vsys they're currently in.
    if err = c.con.VsysUnimport(util.VlanImport, "", "", []string{e.Name}); err != nil {
        return err
    }

    // Perform vsys import next.
    return c.con.VsysImport(util.VlanImport, "", "", vsys, []string{e.Name})
}

// Delete removes the given VLAN(s) from the firewall.
//
// VLANs can be a string or an Entry object.
func (c *FwVlan) Delete(e ...interface{}) error {
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
    c.con.LogAction("(delete) %s: %v", plural, names)

    // Unimport VLANs.
    if err = c.con.VsysUnimport(util.VlanImport, "", "", names); err != nil {
        return err
    }

    // Remove VLANs next.
    path := c.xpath(names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *FwVlan) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *FwVlan) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwVlan) xpath(vals []string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "vlan",
        util.AsEntryXpath(vals),
    }
}
