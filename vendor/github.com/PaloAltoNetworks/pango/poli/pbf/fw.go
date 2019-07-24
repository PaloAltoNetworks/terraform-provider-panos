package pbf

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// FwPbf is the client.Policies.PolicyBasedForwarding namespace.
type FwPbf struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwPbf) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *FwPbf) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of values.
func (c *FwPbf) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given uid.
func (c *FwPbf) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, vsys, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *FwPbf) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more objects.
func (c *FwPbf) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct.
    d := util.BulkElement{XMLName: xml.Name{Local: "temp"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) %s: %v", plural, names)

    // Set xpath.
    path := c.xpath(vsys, names)
    d.XMLName = xml.Name{Local: path[len(path) - 2]}
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
func (c *FwPbf) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Edit the object.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *FwPbf) Delete(vsys string, e ...interface{}) error {
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
            return fmt.Errorf("Unknown type sent to delete: %s", v)
        }
    }
    c.con.LogAction("(delete) %s: %v", plural, names)

    // Remove the objects.
    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

// MoveGroup moves a logical group of policy based forwarding rules somewhere
// in relation to another rule.
func (c *FwPbf) MoveGroup(vsys string, mvt int, rule string, e ...Entry) error {
    var err error

    c.con.LogAction("(move) %s group", singular)

    if len(e) < 1 {
        return fmt.Errorf("Requires at least one rule")
    }

    path := c.xpath(vsys, []string{e[0].Name})
    list, err := c.GetList(vsys)
    if err != nil {
        return err
    }

    // Set the first entity's position.
    if err = c.con.PositionFirstEntity(mvt, rule, e[0].Name, path, list); err != nil {
        return err
    }

    // Move all the rest under it.
    li := len(path) - 1
    for i := 1; i < len(e); i++ {
        path[li] = util.AsEntryXpath([]string{e[i].Name})
        if _, err = c.con.Move(path, "after", e[i - 1].Name, nil, nil); err != nil {
            return err
        }
    }

    return nil
}


/** Internal functions for this namespace struct **/

func (c *FwPbf) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{9, 0, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *FwPbf) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwPbf) xpath(vsys string, vals []string) []string {
    ans := make([]string, 0, 9)
    ans = append(ans, util.VsysXpathPrefix(vsys)...)
    ans = append(ans,
        "rulebase",
        "pbf",
        "rules",
        util.AsEntryXpath(vals),
    )

    return ans
}
