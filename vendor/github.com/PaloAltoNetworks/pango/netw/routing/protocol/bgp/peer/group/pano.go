package group

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// PanoGroup is the client.Network.BgpPeerGroup namespace.
type PanoGroup struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *PanoGroup) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *PanoGroup) ShowList(tmpl, ts, vr string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(tmpl, ts, vr, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of values.
func (c *PanoGroup) GetList(tmpl, ts, vr string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(tmpl, ts, vr, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given uid.
func (c *PanoGroup) Get(tmpl, ts, vr, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, tmpl, ts, vr, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *PanoGroup) Show(tmpl, ts, vr, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, tmpl, ts, vr, name)
}

// Set performs SET to create / update one or more objects.
func (c *PanoGroup) Set(tmpl, ts, vr string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct.
    d := util.BulkElement{XMLName: xml.Name{Local: "peer-group"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) %s: %v", plural, names)

    // Set xpath.
    path := c.xpath(tmpl, ts, vr, names)
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
func (c *PanoGroup) Edit(tmpl, ts, vr string, e Entry) error {
    var err error

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    }

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath(tmpl, ts, vr, []string{e.Name})

    // Edit the object.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *PanoGroup) Delete(tmpl, ts, vr string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
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
    path := c.xpath(tmpl, ts, vr, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *PanoGroup) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *PanoGroup) details(fn util.Retriever, tmpl, ts, vr, name string) (Entry, error) {
    path := c.xpath(tmpl, ts, vr, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoGroup) xpath(tmpl, ts, vr string, vals []string) []string {
    ans := make([]string, 0, 15)
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
        "peer-group",
        util.AsEntryXpath(vals),
    )

    return ans
}
