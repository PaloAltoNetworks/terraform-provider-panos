package eth

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// FwEth is the client.Network.EthernetInterface namespace.
type FwEth struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwEth) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of ethernet interfaces.
func (c *FwEth) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of ethernet interfaces")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of ethernet interfaces.
func (c *FwEth) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of ethernet interfaces")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given ethernet interface.
func (c *FwEth) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) ethernet interface %q", name)
    return c.details(c.con.Get, name)
}

// Show performs SHOW to retrieve information for the given ethernet interface.
func (c *FwEth) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) ethernet interface %q", name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more ethernet interfaces.
//
// Specifying a non-empty vsys will import the interfaces into that vsys,
// allowing the vsys to use them, as long as the interface does not have a
// mode of "ha" or "aggregate-group".  Interfaces that have either of those
// modes are omitted from this function's followup vsys import.
func (c *FwEth) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    n1 := make([]string, len(e))
    n2 := make([]string, 0, len(e))

    // Build up the struct with the given interface configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "ethernet"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        n1[i] = e[i].Name
        if e[i].Mode != "ha" && e[i].Mode != "aggregate-group" {
            n2 = append(n2, e[i].Name)
        }
    }
    c.con.LogAction("(set) ethernet interfaces: %v", n1)

    // Set xpath.
    path := c.xpath(n1)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the interfaces.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    if err != nil {
        return err
    }

    // Remove the interfaces from any vsys they're currently in.
    if err = c.con.VsysUnimport(util.InterfaceImport, "", "", n2); err != nil {
        return err
    }

    // Perform vsys import next.
    return c.con.VsysImport(util.InterfaceImport, "", "", vsys, n2)
}

// Edit performs EDIT to create / update the specified ethernet interface.
//
// Specifying a non-empty vsys will import the interface into that vsys,
// allowing the vsys to use it, as long as the interface does not have a
// mode of "ha" or "aggregate-group".  Interfaces that have either of those
// modes are omitted from this function's followup vsys import.
func (c *FwEth) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) ethernet interface %q", e.Name)

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit the interface.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    if err != nil {
        return err
    }

    // Check if we should skip the import step.
    if e.Mode == "ha" || e.Mode == "aggregate-group" {
        return nil
    }

    // Remove the interface from any vsys it's currently in.
    if err = c.con.VsysUnimport(util.InterfaceImport, "", "", []string{e.Name}); err != nil {
        return err
    }

    // Import the interface.
    return c.con.VsysImport(util.InterfaceImport, "", "", vsys, []string{e.Name})
}

// Delete removes the given interface(s) from the firewall.
//
// Interfaces can be a string or an Entry object.
func (c *FwEth) Delete(e ...interface{}) error {
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
    c.con.LogAction("(delete) ethernet interface(s): %v", names)

    // Unimport interfaces.
    if err = c.con.VsysUnimport(util.InterfaceImport, "", "", names); err != nil {
        return err
    }

    // Remove interfaces next.
    path := c.xpath(names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *FwEth) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{7, 1, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *FwEth) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwEth) xpath(vals []string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "interface",
        "ethernet",
        util.AsEntryXpath(vals),
    }
}
