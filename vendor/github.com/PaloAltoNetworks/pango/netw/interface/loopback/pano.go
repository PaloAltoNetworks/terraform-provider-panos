package loopback

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// PanoLoopback is the client.Network.LoopbackInterface namespace.
type PanoLoopback struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *PanoLoopback) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of loopback interfaces.
func (c *PanoLoopback) ShowList(tmpl, ts string) ([]string, error) {
    c.con.LogQuery("(show) list of loopback interfaces")
    path := c.xpath(tmpl, ts, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of loopback interfaces.
func (c *PanoLoopback) GetList(tmpl, ts string) ([]string, error) {
    c.con.LogQuery("(get) list of loopback interfaces")
    path := c.xpath(tmpl, ts, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given loopback interface.
func (c *PanoLoopback) Get(tmpl, ts, name string) (Entry, error) {
    c.con.LogQuery("(get) loopback interface %q", name)
    return c.details(c.con.Get, tmpl, ts, name)
}

// Show performs SHOW to retrieve information for the given loopback interface.
func (c *PanoLoopback) Show(tmpl, ts, name string) (Entry, error) {
    c.con.LogQuery("(show) loopback interface %q", name)
    return c.details(c.con.Show, tmpl, ts, name)
}

// Set performs SET to create / update one or more loopback interfaces.
//
// Specifying a non-empty vsys will import the interfaces into that vsys,
// allowing the vsys to use them.
func (c *PanoLoopback) Set(tmpl, ts, vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vsys == "" {
        return fmt.Errorf("vsys must be specified")
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given interface configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "units"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) loopback interfaces: %v", names)

    // Set xpath.
    path := c.xpath(tmpl, ts, names)
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
    if err = c.con.VsysUnimport(util.InterfaceImport, tmpl, ts, names); err != nil {
        return err
    }

    // Perform vsys import next.
    return c.con.VsysImport(util.InterfaceImport, tmpl, ts, vsys, names)
}

// Edit performs EDIT to create / update the specified loopback interface.
//
// Specifying a non-empty vsys will import the interface into that vsys,
// allowing the vsys to use it.
func (c *PanoLoopback) Edit(tmpl, ts, vsys string, e Entry) error {
    var err error

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vsys == "" {
        return fmt.Errorf("vsys must be specified")
    }

    _, fn := c.versioning()

    c.con.LogAction("(edit) loopback interface %q", e.Name)

    // Set xpath.
    path := c.xpath(tmpl, ts, []string{e.Name})

    // Edit the interface.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    if err != nil {
        return err
    }

    // Remove the interface from any vsys it's currently in.
    if err = c.con.VsysUnimport(util.InterfaceImport, tmpl, ts, []string{e.Name}); err != nil {
        return err
    }

    // Import the interface.
    return c.con.VsysImport(util.InterfaceImport, tmpl, ts, vsys, []string{e.Name})
}

// Delete removes the given loopback interface(s) from the firewall.
//
// Interfaces can be a string or an Entry object.
func (c *PanoLoopback) Delete(tmpl, ts string, e ...interface{}) error {
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
    c.con.LogAction("(delete) loopback interfaces: %v", names)

    // Unimport interfaces.
    if err = c.con.VsysUnimport(util.InterfaceImport, tmpl, ts, names); err != nil {
        return err
    }

    // Remove interfaces next.
    path := c.xpath(tmpl, ts, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *PanoLoopback) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{7, 1, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *PanoLoopback) details(fn util.Retriever, tmpl, ts, name string) (Entry, error) {
    path := c.xpath(tmpl, ts, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoLoopback) xpath(tmpl, ts string, vals []string) []string {
    ans := make([]string, 0, 13)
    ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
    ans = append(ans,
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "interface",
        "loopback",
        "units",
        util.AsEntryXpath(vals),
    )

    return ans
}
