package ikegw

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// FwIkeGw is a namespace struct, included as part of pango.Client.
type FwIkeGw struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *FwIkeGw) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of IKE gateways.
func (c *FwIkeGw) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of ike gateways")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of IKE gateways.
func (c *FwIkeGw) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of ike gateways")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given IKE gateway.
func (c *FwIkeGw) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) ike gateway %q", name)
    return c.details(c.con.Get, name)
}

// Get performs SHOW to retrieve information for the given IKE gateway.
func (c *FwIkeGw) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) ike gateway %q", name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more IKE gateways.
func (c *FwIkeGw) Set(e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "gateway"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) ike gateway: %v", names)

    // Set xpath.
    path := c.xpath(names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update an IKE gateway.
func (c *FwIkeGw) Edit(e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) ike gateway %q", e.Name)

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given IKE gateways from the firewall.
//
// IKE gateways can be either a string or an Entry object.
func (c *FwIkeGw) Delete(e ...interface{}) error {
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
    c.con.LogAction("(delete) ike gateways: %v", names)

    path := c.xpath(names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *FwIkeGw) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{8, 1, 0, ""}) {
        return &container_v4{}, specify_v4
    } else if v.Gte(version.Number{7, 1, 0, ""}) {
        return &container_v3{}, specify_v3
    } else if v.Gte(version.Number{7, 0, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *FwIkeGw) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwIkeGw) xpath(vals []string) []string {
    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "ike",
        "gateway",
        util.AsEntryXpath(vals),
    }
}
