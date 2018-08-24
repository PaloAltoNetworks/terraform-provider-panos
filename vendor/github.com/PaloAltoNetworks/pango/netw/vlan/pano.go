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

// ShowList performs SHOW to retrieve a list of VLANs.
func (c *PanoVlan) ShowList(tmpl, ts string) ([]string, error) {
    c.con.LogQuery("(show) list of VLANs")
    path := c.xpath(tmpl, ts, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of VLANs.
func (c *PanoVlan) GetList(tmpl, ts string) ([]string, error) {
    c.con.LogQuery("(get) list of VLANs")
    path := c.xpath(tmpl, ts, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given VLAN.
func (c *PanoVlan) Get(tmpl, ts, name string) (Entry, error) {
    c.con.LogQuery("(get) VLAN %q", name)
    return c.details(c.con.Get, tmpl, ts, name)
}

// Show performs SHOW to retrieve information for the given VLAN.
func (c *PanoVlan) Show(tmpl, ts, name string) (Entry, error) {
    c.con.LogQuery("(show) VLAN %q", name)
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
    c.con.LogAction("(set) VLANs: %v", names)

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

    c.con.LogAction("(edit) VLAN %q", e.Name)

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
    c.con.LogAction("(delete) VLANs: %v", names)

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
