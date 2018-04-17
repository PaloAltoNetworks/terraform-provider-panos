package addr

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// PanoAddr is a namespace struct, included as part of pango.Firewall.
type PanoAddr struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *PanoAddr) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of address objects.
func (c *PanoAddr) GetList(dg string) ([]string, error) {
    c.con.LogQuery("(get) list of address objects")
    path := c.xpath(dg, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of address objects.
func (c *PanoAddr) ShowList(dg string) ([]string, error) {
    c.con.LogQuery("(show) list of address objects")
    path := c.xpath(dg, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given address object.
func (c *PanoAddr) Get(dg, name string) (Entry, error) {
    c.con.LogQuery("(get) address object %q", name)
    return c.details(c.con.Get, dg, name)
}

// Get performs SHOW to retrieve information for the given address object.
func (c *PanoAddr) Show(dg, name string) (Entry, error) {
    c.con.LogQuery("(show) address object %q", name)
    return c.details(c.con.Show, dg, name)
}

// Set performs SET to create / update one or more address objects.
func (c *PanoAddr) Set(dg string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "address"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) address objects: %v", names)

    // Set xpath.
    path := c.xpath(dg, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the objects.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update an address object.
func (c *PanoAddr) Edit(dg string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) address object %q", e.Name)

    // Set xpath.
    path := c.xpath(dg, []string{e.Name})

    // Create the objects.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given address objects from the firewall.
//
// Address objects can be either a string or an Entry object.
func (c *PanoAddr) Delete(dg string, e ...interface{}) error {
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
            return fmt.Errorf("Unsupported type to delete: %s", v)
        }
    }
    c.con.LogAction("(delete) address objects: %v", names)

    path := c.xpath(dg, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the PanoAddr struct **/

func (c *PanoAddr) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *PanoAddr) details(fn util.Retriever, dg, name string) (Entry, error) {
    path := c.xpath(dg, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoAddr) xpath(dg string, vals []string) []string {
    if dg == "" {
        dg = "shared"
    }

    if dg == "shared" {
        return []string {
            "config",
            "shared",
            "address",
            util.AsEntryXpath(vals),
        }
    }

    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "device-group",
        util.AsEntryXpath([]string{dg}),
        "address",
        util.AsEntryXpath(vals),
    }
}
