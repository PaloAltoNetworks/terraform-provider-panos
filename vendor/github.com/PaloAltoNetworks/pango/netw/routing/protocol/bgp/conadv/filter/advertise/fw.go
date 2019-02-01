package advertise

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// FwAdvertise is the client.Network.BgpConAdvAdvertiseFilter namespace.
type FwAdvertise struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwAdvertise) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *FwAdvertise) ShowList(vr, ca string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(vr, ca, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of values.
func (c *FwAdvertise) GetList(vr, ca string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(vr, ca, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given uid.
func (c *FwAdvertise) Get(vr, ca, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, vr, ca, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *FwAdvertise) Show(vr, ca, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, vr, ca, name)
}

// Set performs SET to create / update one or more objects.
func (c *FwAdvertise) Set(vr, ca string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    } else if ca == "" {
        return fmt.Errorf("ca must be specified")
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
    path := c.xpath(vr, ca, names)
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
func (c *FwAdvertise) Edit(vr, ca string, e Entry) error {
    var err error

    if vr == "" {
        return fmt.Errorf("vr must be specified")
    } else if ca == "" {
        return fmt.Errorf("ca must be specified")
    }

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath(vr, ca, []string{e.Name})

    // Edit the object.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *FwAdvertise) Delete(vr, ca string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    } else if ca == "" {
        return fmt.Errorf("ca must be specified")
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
    path := c.xpath(vr, ca, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *FwAdvertise) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{8, 0, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *FwAdvertise) details(fn util.Retriever, vr, ca, name string) (Entry, error) {
    path := c.xpath(vr, ca, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwAdvertise) xpath(vr, ca string, vals []string) []string {
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
        "conditional-advertisement",
        "policy",
        util.AsEntryXpath([]string{ca}),
        "advertise-filters",
        util.AsEntryXpath(vals),
    }
}
