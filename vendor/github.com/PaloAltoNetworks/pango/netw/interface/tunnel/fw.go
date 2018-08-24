package tunnel

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// FwTunnel is the client.Network.TunnelInterface namespace.
type FwTunnel struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwTunnel) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of tunnel interfaces.
func (c *FwTunnel) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of tunnel interfaces")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of tunnel interfaces.
func (c *FwTunnel) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of tunnel interfaces")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given tunnel interface.
func (c *FwTunnel) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) tunnel interface %q", name)
    return c.details(c.con.Get, name)
}

// Show performs SHOW to retrieve information for the given tunnel interface.
func (c *FwTunnel) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) tunnel interface %q", name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more tunnel interfaces.
//
// Specifying a non-empty vsys will import the interfaces into that vsys,
// allowing the vsys to use them.
func (c *FwTunnel) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given interface configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "units"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) tunnel interfaces: %v", names)

    // Set xpath.
    path := c.xpath(names)
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
    if err = c.con.VsysUnimport(util.InterfaceImport, "", "", names); err != nil {
        return err
    }

    // Perform vsys import next.
    return c.con.VsysImport(util.InterfaceImport, "", "", vsys, names)
}

// Edit performs EDIT to create / update the specified tunnel interface.
//
// Specifying a non-empty vsys will import the interface into that vsys,
// allowing the vsys to use it.
func (c *FwTunnel) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) tunnel interface %q", e.Name)

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit the interface.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    if err != nil {
        return err
    }

    // Remove the interface from any vsys it's currently in.
    if err = c.con.VsysUnimport(util.InterfaceImport, "", "", []string{e.Name}); err != nil {
        return err
    }

    // Import the interface.
    return c.con.VsysImport(util.InterfaceImport, "", "", vsys, []string{e.Name})
}

// Delete removes the given tunnel interface(s) from the firewall.
//
// Interfaces can be a string or an Entry object.
func (c *FwTunnel) Delete(e ...interface{}) error {
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
    c.con.LogAction("(delete) tunnel interfaces: %v", names)

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

func (c *FwTunnel) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *FwTunnel) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwTunnel) xpath(vals []string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "interface",
        "tunnel",
        "units",
        util.AsEntryXpath(vals),
    }
}
