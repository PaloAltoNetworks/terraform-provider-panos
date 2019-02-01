package suppress

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// FwSuppress is the client.Network.BgpAggSuppressFilter namespace.
type FwSuppress struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwSuppress) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *FwSuppress) ShowList(vr, ag string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(vr, ag, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of values.
func (c *FwSuppress) GetList(vr, ag string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(vr, ag, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given uid.
func (c *FwSuppress) Get(vr, ag, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, vr, ag, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *FwSuppress) Show(vr, ag, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, vr, ag, name)
}

// Set performs SET to create / update one or more objects.
func (c *FwSuppress) Set(vr, ag string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    } else if ag == "" {
        return fmt.Errorf("ag must be specified")
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct.
    d := util.BulkElement{XMLName: xml.Name{Local: "suppress-filters"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) %s: %v", plural, names)

    // Set xpath.
    path := c.xpath(vr, ag, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the objects.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update one object.
func (c *FwSuppress) Edit(vr, ag string, e Entry) error {
    var err error

    if vr == "" {
        return fmt.Errorf("vr must be specified")
    } else if ag == "" {
        return fmt.Errorf("ag must be specified")
    }

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath(vr, ag, []string{e.Name})

    // Edit the object.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *FwSuppress) Delete(vr, ag string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    } else if ag == "" {
        return fmt.Errorf("ag must be specified")
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

    // Remove the objects.
    path := c.xpath(vr, ag, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *FwSuppress) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{8, 0, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *FwSuppress) details(fn util.Retriever, vr, ag, name string) (Entry, error) {
    path := c.xpath(vr, ag, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwSuppress) xpath(vr, ag string, vals []string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "virtual-router",
        util.AsEntryXpath([]string{vr}),
        "protocol",
        "bgp",
        "policy",
        "aggregation",
        "address",
        util.AsEntryXpath([]string{ag}),
        "suppress-filters",
        util.AsEntryXpath(vals),
    }
}
