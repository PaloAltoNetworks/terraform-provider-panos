package vlan

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// PanoVlan is the client.Network.Vlan namespace.
type PanoVlan struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *PanoVlan) Initialize(con util.XapiClient) {
    c.con = con
}

/*
SetInterface performs a SET to add an interface to a VLAN.

The VLAN can be either a string or an Entry object.
The iface variable is the interface.
The rmMacs and addMacs params are MAC addresses to remove/add that
will reference the iface interface.
*/
func (c *PanoVlan) SetInterface(tmpl, ts string, vlan interface{}, iface string, rmMacs, addMacs []string) error {
    var (
        name string
        err error
    )

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    }

    switch v := vlan.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to %s set interface: %s", singular, v)
    }

    c.con.LogAction("(set) interface for %s %q: %s", singular, name, iface)

    basePath := c.xpath(tmpl, ts, []string{name})
    iPath := append(basePath, "interface")

    if _, err = c.con.Set(iPath, util.Member{Value: iface}, nil, nil); err != nil {
        return err
    }

    if len(rmMacs) > 0 {
        c.con.LogAction("(delete) removing %q mac addresses: %#v", name, rmMacs)
        rPath := append(basePath, "mac", util.AsEntryXpath(rmMacs))
        if _, err = c.con.Delete(rPath, nil, nil); err != nil {
            return err
        }
    }

    if len(addMacs) > 0 {
        c.con.LogAction("(set) adding %q mac addresses: %#v", name, addMacs)
        d := util.BulkElement{XMLName: xml.Name{Local: "mac"}}
        for i := range addMacs {
            d.Data = append(d.Data, macList{Mac: addMacs[i], Interface: iface})
        }
        aPath := make([]string, 0, len(basePath) + 1)
        aPath = append(aPath, basePath...)
        if len(addMacs) == 1 {
            aPath = append(aPath, "mac")
        }
        if _, err = c.con.Set(aPath, d.Config(), nil, nil); err != nil {
            return err
        }
    }

    return nil
}

/*
DeleteInterface performs a DELETE to remove an interface from a VLAN.

The VLAN can be either a string or an Entry object.

All MAC addresses referencing this interface are deleted.
*/
func (c *PanoVlan) DeleteInterface(tmpl, ts string, vlan interface{}, iface string) error {
    var (
        name string
        err error
    )

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    }

    switch v := vlan.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to %s delete interface: %s", singular, v)
    }

    o, err := c.Get(tmpl, ts, name)
    if err != nil {
        return err
    }
    rmMacs := make([]string, 0, len(o.StaticMacs))
    for k, v := range o.StaticMacs {
        if v == iface {
            rmMacs = append(rmMacs, k)
        }
    }

    c.con.LogAction("(delete) interface for %s %q: %s", singular, name, iface)

    basePath := c.xpath(tmpl, ts, []string{name})
    mPath := append(basePath, "mac", util.AsEntryXpath(rmMacs))
    iPath := append(basePath, "interface", util.AsMemberXpath([]string{iface}))

    if len(rmMacs) > 0 {
        c.con.LogAction("(delete) removing %q mac addresses: %#v", iface, rmMacs)
        if _, err = c.con.Delete(mPath, nil, nil); err != nil {
            return err
        }
    }

    _, err = c.con.Delete(iPath, nil, nil)
    return err
}

// ShowList performs SHOW to retrieve a list of VLANs.
func (c *PanoVlan) ShowList(tmpl, ts string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(tmpl, ts, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of VLANs.
func (c *PanoVlan) GetList(tmpl, ts string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(tmpl, ts, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given VLAN.
func (c *PanoVlan) Get(tmpl, ts, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, tmpl, ts, name)
}

// Show performs SHOW to retrieve information for the given VLAN.
func (c *PanoVlan) Show(tmpl, ts, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, tmpl, ts, name)
}

// Set performs SET to create / update one or more VLANs.
//
// Specify a non-empty vsys to import the VLAN(s) into the given vsys
// after creating, allowing the vsys to use them.
func (c *PanoVlan) Set(tmpl, ts, vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
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
    path := c.xpath(tmpl, ts, names)
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
    if err = c.con.VsysUnimport(util.VlanImport, tmpl, ts, names); err != nil {
        return err
    }

    // Perform vsys import next.
    return c.con.VsysImport(util.VlanImport, tmpl, ts, vsys, names)
}

// Edit performs EDIT to create / update a VLAN.
//
// Specify a non-empty vsys to import the VLAN into the given vsys
// after creating, allowing the vsys to use it.
func (c *PanoVlan) Edit(tmpl, ts, vsys string, e Entry) error {
    var err error

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    }

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath(tmpl, ts, []string{e.Name})

    // Edit the VLAN.
    if _, err = c.con.Edit(path, fn(e), nil, nil); err != nil {
        return err
    }

    // Remove the VLANs from any vsys they're currently in.
    if err = c.con.VsysUnimport(util.VlanImport, tmpl, ts, []string{e.Name}); err != nil {
        return err
    }

    // Perform vsys import next.
    return c.con.VsysImport(util.VlanImport, tmpl, ts, vsys, []string{e.Name})
}

// Delete removes the given VLAN(s) from the firewall.
//
// VLANs can be a string or an Entry object.
func (c *PanoVlan) Delete(tmpl, ts string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
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
    if err = c.con.VsysUnimport(util.VlanImport, tmpl, ts, names); err != nil {
        return err
    }

    // Remove VLANs next.
    path := c.xpath(tmpl, ts, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *PanoVlan) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *PanoVlan) details(fn util.Retriever, tmpl, ts, name string) (Entry, error) {
    path := c.xpath(tmpl, ts, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoVlan) xpath(tmpl, ts string, vals []string) []string {
    ans := make([]string, 0, 11)
    ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
    ans = append(ans,
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "vlan",
        util.AsEntryXpath(vals),
    )

    return ans
}
