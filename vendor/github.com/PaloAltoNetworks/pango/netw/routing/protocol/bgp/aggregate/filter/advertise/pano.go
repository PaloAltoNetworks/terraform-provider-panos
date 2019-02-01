package advertise

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// PanoAdvertise is the client.Network.BgpAggAdvertiseFilter namespace.
type PanoAdvertise struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *PanoAdvertise) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *PanoAdvertise) ShowList(tmpl, ts, vr, ag string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(tmpl, ts, vr, ag, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of values.
func (c *PanoAdvertise) GetList(tmpl, ts, vr, ag string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(tmpl, ts, vr, ag, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given uid.
func (c *PanoAdvertise) Get(tmpl, ts, vr, ag, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, tmpl, ts, vr, ag, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *PanoAdvertise) Show(tmpl, ts, vr, ag, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, tmpl, ts, vr, ag, name)
}

// Set performs SET to create / update one or more objects.
func (c *PanoAdvertise) Set(tmpl, ts, vr, ag string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    } else if ag == "" {
        return fmt.Errorf("ag must be specified")
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct.
    d := util.BulkElement{XMLName: xml.Name{Local: "advertise-filters"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) %s: %v", plural, names)

    // Set xpath.
    path := c.xpath(tmpl, ts, vr, ag, names)
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
func (c *PanoAdvertise) Edit(tmpl, ts, vr, ag string, e Entry) error {
    var err error

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    } else if ag == "" {
        return fmt.Errorf("ag must be specified")
    }

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath(tmpl, ts, vr, ag, []string{e.Name})

    // Edit the object.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *PanoAdvertise) Delete(tmpl, ts, vr, ag string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
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
    path := c.xpath(tmpl, ts, vr, ag, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *PanoAdvertise) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{8, 0, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *PanoAdvertise) details(fn util.Retriever, tmpl, ts, vr, ag, name string) (Entry, error) {
    path := c.xpath(tmpl, ts, vr, ag, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoAdvertise) xpath(tmpl, ts, vr, ag string, vals []string) []string {
    ans := make([]string, 0, 19)
    ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
    ans = append(ans,
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
        "advertise-filters",
        util.AsEntryXpath(vals),
    )

    return ans
}
