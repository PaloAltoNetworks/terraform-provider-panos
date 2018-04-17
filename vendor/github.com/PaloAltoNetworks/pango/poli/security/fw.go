package security

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// FwSecurity is the client.Policies.Security namespace.
type FwSecurity struct {
    con util.XapiClient
}

// Initialize is invoed by client.Initialize().
func (c *FwSecurity) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of security policies.
func (c *FwSecurity) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of security policies")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of security policies.
func (c *FwSecurity) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of security policies")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given security policy.
func (c *FwSecurity) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) security policy %q", name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given security policy.
func (c *FwSecurity) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) security policy %q", name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more security policies.
func (c *FwSecurity) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given security policy configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "rules"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) security policies: %v", names)

    // Set xpath.
    path := c.xpath(vsys, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the security policies.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// VerifiableSet behaves like Set(), except policies with LogEnd as true
// will first be created with LogEnd as false, and then a second Set() is
// performed which will do LogEnd as true.  This is due to the unique
// combination of being a boolean value that is true by default, the XML
// returned from querying the rule details will omit the LogEnd setting,
// which will be interpreted as false, when in fact it is true.  We can
// get around this by setting the value to a non-standard value, then back
// again, in which case it will properly show up in the returned XML.
func (c *FwSecurity) VerifiableSet(vsys string, e ...Entry) error {
    c.con.LogAction("(set) performing verifiable set")
    again := make([]Entry, 0, len(e))

    for i := range e {
        if e[i].LogEnd {
            again = append(again, e[i])
            e[i].LogEnd = false
        }
    }

    if err := c.Set(vsys, e...); err != nil {
        return err
    }

    if len(again) == 0 {
        return nil
    }

    return c.Set(vsys, again...)
}

// Edit performs EDIT to create / update a security policy.
func (c *FwSecurity) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) security policy %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Edit the security policy.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given security policies.
//
// Security policies can be either a string or an Entry object.
func (c *FwSecurity) Delete(vsys string, e ...interface{}) error {
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
    c.con.LogAction("(delete) security policies: %v", names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

// DeleteAll removes all security policies from the specified vsys / rulebase.
func (c *FwSecurity) DeleteAll(vsys string) error {
    c.con.LogAction("(delete) all security policies")
    list, err := c.GetList(vsys)
    if err != nil || len(list) == 0 {
        return err
    }
    li := make([]interface{}, len(list))
    for i := range list {
        li[i] = list[i]
    }
    return c.Delete(vsys, li...)
}

/** Internal functions for the FwSecurity struct **/

func (c *FwSecurity) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *FwSecurity) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwSecurity) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    if vsys == "shared" {
        return []string{
            "config",
            "shared",
            "rulebase",
            "security",
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
        "security",
        "rules",
        util.AsEntryXpath(vals),
    }
}
