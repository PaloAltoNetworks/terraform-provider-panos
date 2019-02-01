package exp

import (
    "fmt"
    "encoding/xml"
    "strings"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// PanoExp is the client.Network.BgpExport namespace.
type PanoExp struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *PanoExp) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *PanoExp) ShowList(tmpl, ts, vr string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(tmpl, ts, vr, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of values.
func (c *PanoExp) GetList(tmpl, ts, vr string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(tmpl, ts, vr, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given uid.
func (c *PanoExp) Get(tmpl, ts, vr, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, tmpl, ts, vr, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *PanoExp) Show(tmpl, ts, vr, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, tmpl, ts, vr, name)
}

// Set performs SET to create / update one or more objects.
func (c *PanoExp) Set(tmpl, ts, vr string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    } else {
        // Make sure rule names are unique.
        m := make(map[string] int)
        for i := range e {
            m[e[i].Name] = m[e[i].Name] + 1
            if m[e[i].Name] > 1 {
                return fmt.Errorf("%s is defined multiple times: %s", singular, e[i].Name)
            }
        }
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct.
    d := util.BulkElement{XMLName: xml.Name{Local: "rules"}}
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

    // On error: find the rule that's causing the error if multiple rules
    // were given.
    if err != nil && strings.Contains(err.Error(), "rules is invalid") {
        for i := 0; i < len(e); i++ {
            if e2 := c.Set(tmpl, ts, vr, e[i]); e2 != nil {
                return fmt.Errorf("Error with rule %d: %s", i + 1, e2)
            } else {
                _ = c.Delete(tmpl, ts, vr, e[i])
            }
        }

        // Couldn't find it, just return the original error.
        return err
    }

    return err
}

// Edit performs EDIT to create / update one object.
func (c *PanoExp) Edit(tmpl, ts, vr string, e Entry) error {
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
func (c *PanoExp) Delete(tmpl, ts, vr string, e ...interface{}) error {
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

// MoveGroup moves a logical group of BGP export rules somewhere in relation
// to another rule.
func (c *PanoExp) MoveGroup(tmpl, ts, vr string, mvt int, rule string, e ...Entry) error {
    var err error

    c.con.LogAction("(move) %s group", singular)

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if len(e) < 1 {
        return fmt.Errorf("Requires at least one rule")
    }

    path := c.xpath(tmpl, ts, vr, []string{e[0].Name})
    list, err := c.GetList(tmpl, ts, vr)
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

func (c *PanoExp) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{8, 0, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *PanoExp) details(fn util.Retriever, tmpl, ts, vr, name string) (Entry, error) {
    path := c.xpath(tmpl, ts, vr, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoExp) xpath(tmpl, ts, vr string, vals []string) []string {
    ans := make([]string, 0, 17)
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
        "export",
        "rules",
        util.AsEntryXpath(vals),
    )

    return ans
}
