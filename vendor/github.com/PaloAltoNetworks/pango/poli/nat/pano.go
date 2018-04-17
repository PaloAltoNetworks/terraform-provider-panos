package nat

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// PanoNat is the client.Policies.Nat namespace.
type PanoNat struct {
    con util.XapiClient
}

// Initialize is invoed by client.Initialize().
func (c *PanoNat) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of NAT policies.
func (c *PanoNat) GetList(dg, base string) ([]string, error) {
    c.con.LogQuery("(get) list of NAT policies")
    path := c.xpath(dg, base, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of NAT policies.
func (c *PanoNat) ShowList(dg, base string) ([]string, error) {
    c.con.LogQuery("(show) list of NAT policies")
    path := c.xpath(dg, base, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given NAT policy.
func (c *PanoNat) Get(dg, base, name string) (Entry, error) {
    c.con.LogQuery("(get) NAT policy %q", name)
    return c.details(c.con.Get, dg, base, name)
}

// Get performs SHOW to retrieve information for the given NAT policy.
func (c *PanoNat) Show(dg, base, name string) (Entry, error) {
    c.con.LogQuery("(show) NAT policy %q", name)
    return c.details(c.con.Show, dg, base, name)
}

// Set performs SET to create / update one or more NAT policies.
func (c *PanoNat) Set(dg, base string, e ...Entry) error {
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
    path := c.xpath(dg, base, names)
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
func (c *PanoNat) Edit(dg, base string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) NAT policy %q", e.Name)

    // Set xpath.
    path := c.xpath(dg, base, []string{e.Name})

    // Edit the NAT policy.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given NAT policies.
//
// NAT policies can be either a string or an Entry object.
func (c *PanoNat) Delete(dg, base string, e ...interface{}) error {
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

    path := c.xpath(dg, base, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the Zone struct **/

func (c *PanoNat) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *PanoNat) details(fn util.Retriever, dg, base, name string) (Entry, error) {
    path := c.xpath(dg, base, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoNat) xpath(dg, base string, vals []string) []string {
    if dg == "" {
        dg = "shared"
    }
    if base == "" {
        base = util.PreRulebase
    }

    if dg == "shared" {
        return []string{
            "config",
            "shared",
            base,
            "nat",
            "rules",
            util.AsEntryXpath(vals),
        }
    }

    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "device-group",
        util.AsEntryXpath([]string{dg}),
        base,
        "nat",
        "rules",
        util.AsEntryXpath(vals),
    }
}
