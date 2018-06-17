package nat

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)

// FwNat is the client.Policies.Nat namespace.
type FwNat struct {
    con util.XapiClient
}

// Initialize is invoed by client.Initialize().
func (c *FwNat) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of NAT policies.
func (c *FwNat) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of NAT policies")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of NAT policies.
func (c *FwNat) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of NAT policies")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given NAT policy.
func (c *FwNat) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) NAT policy %q", name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given NAT policy.
func (c *FwNat) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) NAT policy %q", name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more NAT policies.
func (c *FwNat) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "rules"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) NAT policies: %v", names)

    // Set xpath.
    path := c.xpath(vsys, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the NAT policies.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update a NAT policy.
func (c *FwNat) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) NAT policy %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Edit the NAT policy.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given NAT policies.
//
// NAT policies can be either a string or an Entry object.
func (c *FwNat) Delete(vsys string, e ...interface{}) error {
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
    c.con.LogAction("(delete) NAT policies: %v", names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the Zone struct **/

func (c *FwNat) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{8, 1, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *FwNat) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwNat) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    if vsys == "shared" {
        return []string{
            "config",
            "shared",
            "rulebase",
            "nat",
            "rules",
            util.AsEntryXpath(vals),
        }
    }

    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "rulebase",
        "nat",
        "rules",
        util.AsEntryXpath(vals),
    }
}
